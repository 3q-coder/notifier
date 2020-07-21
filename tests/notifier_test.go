package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DryginAlexander/notifier"
	"github.com/DryginAlexander/notifier/models"
	"github.com/DryginAlexander/notifier/operator"
	"github.com/DryginAlexander/notifier/settings"
	"github.com/DryginAlexander/notifier/web"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var client = http.Client{Timeout: time.Duration(5 * time.Second)}

func initSettingsAndDB(t *testing.T) (models.Storage, operator.Operator, *gin.Engine) {
	err := settings.Init("../env/test.env")
	assert.Nil(t, err, "failed to initialize settings")
	stor := models.NewStorage()
	err = stor.MigrateDB()
	assert.Nil(t, err, "failed to migrate db", err)
	// clear test db
	if err = stor.DB.Delete(&models.User{}).Error; err != nil {
		t.Errorf("failed to delete all records in table Users %v", err)
	}
	if err = stor.DB.Delete(&models.Notification{}).Error; err != nil {
		t.Errorf("failed to delete all records in table Notifications %v", err)
	}
	oper := operator.NewOperator(&stor)
	return stor, oper, web.Init(&stor, &oper)
}

func TestSubscribeToNotifications(t *testing.T) {
	name := "test"
	msg := "test message"

	// init server
	stor, oper, router := initSettingsAndDB(t)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// create notification
	note := notifier.Notification{
		Username: name,
		Message:  msg,
		Sent:     false,
	}
	oper.SendNotification(&note)

	// connect to the server
	url := "ws" + strings.TrimPrefix(ts.URL, "http")
	url += "/notifications?token=" + name
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.Nil(t, err, "failed to connect to the server ", err)
	defer ws.Close()

	// check that notification has been sent
	_, data, err := ws.ReadMessage()
	assert.Nil(t, err, "failed to read message ", err)
	var wsMsg notifier.Notification
	json.Unmarshal(data, &wsMsg)
	assert.Equal(t, wsMsg.Message, msg)

	// check that note has been marked as sent in db
	time.Sleep(100 * time.Millisecond)
	var noteDB models.Notification
	err = stor.DB.Where("username = ? AND message = ?", name, msg).First(&noteDB).Error
	assert.Nil(t, err, "failed to load note from db", err)
	assert.Equal(t, true, noteDB.Sent)
}

func TestSendNotification(t *testing.T) {
	names := [2]string{"test", "test2"}
	msgs := [2]string{"test message", "test message 2"}

	// init server
	stor, _, router := initSettingsAndDB(t)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// connect to the server
	var url string
	var wss []*websocket.Conn
	for _, name := range names {
		url = "ws" + strings.TrimPrefix(ts.URL, "http")
		url += "/notifications?token=" + name
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		assert.Nil(t, err, "failed to connect to the server ", err)
		defer ws.Close()
		wss = append(wss, ws)
	}

	// send notification
	url = ts.URL + "/send-note"
	for i, msg := range msgs {
		data, _ := json.Marshal(notifier.Notification{
			Username: names[i],
			Message:  msg,
		})
		body := bytes.NewBufferString(string(data))

		req, _ := http.NewRequest(http.MethodPost, url, body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(data)))
		req.Header.Add("API-KEY", settings.APIKey)

		resp, err := client.Do(req)
		assert.Nil(t, err, "failed to send message", err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// check that notification has been sent
	for i, ws := range wss {
		_, data, err := ws.ReadMessage()
		assert.Nil(t, err, "failed to read message ", err)
		var wsMsg notifier.Notification
		json.Unmarshal(data, &wsMsg)
		assert.Equal(t, wsMsg.Message, msgs[i])
	}

	// check that note has been marked as sent in db
	time.Sleep(100 * time.Millisecond)
	var err error
	var noteDB models.Notification
	for i, msg := range msgs {
		noteDB = models.Notification{}
		err = stor.DB.Where("username = ? AND message = ?", names[i], msg).First(&noteDB).Error
		assert.Nil(t, err, "failed to load note from db", err)
		assert.Equal(t, true, noteDB.Sent)
	}
}

func TestSendNotificationAll(t *testing.T) {
	var names []string
	msgs := [1]string{"testMessage"}
	for i := 0; i < 20; i++ {
		names = append(names, "test"+strconv.Itoa(i))
	}

	// init server
	stor, _, router := initSettingsAndDB(t)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// create users
	for _, name := range names {
		stor.CreateUser(&notifier.User{
			Username: name,
			Password: name,
		})
	}

	// connect to the server
	var url string
	var wss []*websocket.Conn
	for _, name := range names {
		url = "ws" + strings.TrimPrefix(ts.URL, "http")
		url += "/notifications?token=" + name
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		assert.Nil(t, err, "failed to connect to the server ", err)
		defer ws.Close()
		wss = append(wss, ws)
	}

	// send notification
	url = ts.URL + "/send-note-all"
	for _, msg := range msgs {
		data, _ := json.Marshal(notifier.Notification{
			Message: msg,
		})
		body := bytes.NewBufferString(string(data))

		req, _ := http.NewRequest(http.MethodPost, url, body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(data)))
		req.Header.Add("API-KEY", settings.APIKey)

		resp, err := client.Do(req)
		assert.Nil(t, err, "failed to send message", err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// check that notification has been sent
	for _, msg := range msgs {
		for _, ws := range wss {
			_, data, err := ws.ReadMessage()
			assert.Nil(t, err, "failed to read message ", err)
			var wsMsg notifier.Notification
			json.Unmarshal(data, &wsMsg)
			assert.Equal(t, wsMsg.Message, msg)
		}
	}

	// check that note has been marked as sent in db
	time.Sleep(100 * time.Millisecond)
	var err error
	var noteDB models.Notification
	for _, msg := range msgs {
		for _, name := range names {
			noteDB = models.Notification{}
			err = stor.DB.Where("username = ? AND message = ?", name, msg).First(&noteDB).Error
			assert.Nil(t, err, "failed to load note from db", err)
			assert.Equal(t, true, noteDB.Sent)
		}
	}
}

func TestMetrics(t *testing.T) {
	var names []string
	msgs := [1]string{"testMessage"}
	for i := 0; i < 100; i++ {
		names = append(names, "test"+strconv.Itoa(i))
	}

	// init server
	stor, _, router := initSettingsAndDB(t)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// create users
	for _, name := range names {
		stor.CreateUser(&notifier.User{
			Username: name,
			Password: name,
		})
	}

	// connect to the server
	var url string
	var wss []*websocket.Conn
	for i, name := range names {
		if i == 50 {
			break
		}
		url = "ws" + strings.TrimPrefix(ts.URL, "http")
		url += "/notifications?token=" + name
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		assert.Nil(t, err, "failed to connect to the server ", err)
		defer ws.Close()
		wss = append(wss, ws)
	}

	// send notification
	url = ts.URL + "/send-note-all"
	for _, msg := range msgs {
		data, _ := json.Marshal(notifier.Notification{
			Message: msg,
		})
		body := bytes.NewBufferString(string(data))

		req, _ := http.NewRequest(http.MethodPost, url, body)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(data)))
		req.Header.Add("API-KEY", settings.APIKey)

		resp, err := client.Do(req)
		assert.Nil(t, err, "failed to send message", err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// check that notification has been sent
	for _, msg := range msgs {
		for _, ws := range wss {
			_, data, err := ws.ReadMessage()
			assert.Nil(t, err, "failed to read message ", err)
			var wsMsg notifier.Notification
			json.Unmarshal(data, &wsMsg)
			assert.Equal(t, wsMsg.Message, msg)
		}
	}

	// check metrics
	time.Sleep(150 * time.Millisecond)
	url = ts.URL + "/metrics"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("API-KEY", settings.APIKey)
	resp, err := client.Do(req)
	assert.Nil(t, err, "failed to request metrics", err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(
		t,
		"all_clients 100\nonline_clients 50\nall_note 100\nsent_note 50\n",
		string(body),
	)
}
