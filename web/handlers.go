package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/DryginAlexander/notifier"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// news -------------------------------------------------------------
// ------------------------------------------------------------------

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func initNoteSocket(c *gin.Context) {
	nameInterface, _ := c.Get("username")
	username := nameInterface.(string)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	go operator.SubscribeToNotifications(username, ws)
}

func sendNote(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	note := notifier.Notification{}
	err = json.Unmarshal(data, &note)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	note.Sent = false

	err = operator.SendNotification(&note)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func sendNoteAll(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	note := &notifier.Notification{}
	err = json.Unmarshal(data, note)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	err = operator.SendNotificationAll(note.Message)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func showNewsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "news.html", gin.H{"title": "News Page"})
}

func metrics(c *gin.Context) {
	allClients, err := storage.UsersNumber()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	onlineClients := operator.OnlineClientsNumber()
	allNotes, err := storage.NotesNumber()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	sentNotes, err := storage.SentNotesNumber()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	output := fmt.Sprintf("all_clients %d\nonline_clients %d\nall_note %d\nsent_note %d\n",
		allClients, onlineClients, allNotes, sentNotes)

	c.String(http.StatusOK, output)
}

// authentification -------------------------------------------------
// ------------------------------------------------------------------

func showIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "Home Page"})
}

func showRegistrationPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"title": "Register"})
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	status, err := storage.IsUsernameAvailable(username)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if strings.TrimSpace(password) == "" {
		err = errors.New("The password can't be empty")
	} else if !status {
		err = errors.New("The username isn't available")
	}

	if err != nil {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(
			http.StatusBadRequest,
			"register.html",
			gin.H{
				"ErrorTitle":   "Registration Failed",
				"ErrorMessage": err.Error(),
			},
		)
	}

	user := notifier.User{
		Username: username,
		Password: password,
	}
	err = storage.CreateUser(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	// If the user is created, set the token in a cookie and log the user in
	token := username
	c.SetCookie("token", token, 3600, "", "", false, true)
	c.Set("is_logged_in", true)

	c.HTML(http.StatusOK, "login-successful.html", gin.H{"title": "Successful registration"})
}

func showLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
}

func performLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	status, err := storage.IsUserValid(username, password)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if status {
		token := username
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		c.HTML(http.StatusOK, "login-successful.html", gin.H{
			"title": "Successful Login",
		})
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided",
		})
	}
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
