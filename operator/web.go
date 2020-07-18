package operator

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/DryginAlexander/notifier"
	"github.com/gorilla/websocket"
	"github.com/icrowley/fake"
)

func (o *Operator) SendNotification(client *websocket.Conn) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		w, err := client.NextWriter(websocket.TextMessage)
		if err != nil {
			ticker.Stop()
			break
		}

		msg := newMessage()
		w.Write(msg)
		w.Close()

		<-ticker.C
	}
}

func newMessage() []byte {
	data, _ := json.Marshal(map[string]string{
		"email":   fake.EmailAddress(),
		"name":    fake.FirstName() + " " + fake.LastName(),
		"subject": fake.Product() + " " + fake.Model(),
	})
	return data
}

func (o *Operator) RegisterUser(username, password string) (*notifier.User, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !o.storage.IsUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	user := notifier.User{
		Username: username,
		Password: password,
	}
	o.storage.CreateUser(&user)
	return &user, nil
}
