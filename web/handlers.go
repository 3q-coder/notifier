package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatal(err)
	}
	go operator.SubscribeToNotifications(username, ws)
}

func sendNote(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)

	note := notifier.Notification{}
	_ = json.Unmarshal(data, &note)
	note.Sent = false

	operator.SendNotification(&note)
}

func sendNoteAll(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)

	note := &notifier.Notification{}
	_ = json.Unmarshal(data, note)

	operator.SendNotificationAll(note.Message)
}

func showNewsPage(c *gin.Context) {
	render(c, gin.H{"title": "News Page"}, "news.html")
}

func metrics(c *gin.Context) {
	allClients, _ := storage.UsersNumber()
	onlineClients := operator.OnlineClientsNumber()
	allNotes, _ := storage.NotesNumber()
	sentNotes, _ := storage.SentNotesNumber()

	output := fmt.Sprintf("all_clients %d\nonline_clients %d\nall_note %d\nsent_note %d\n",
		allClients, onlineClients, allNotes, sentNotes)

	c.String(http.StatusOK, output)
}

// authentification -------------------------------------------------
// ------------------------------------------------------------------

func showIndexPage(c *gin.Context) {
	render(c, gin.H{"title": "Home Page"}, "index.html")
}

func showRegistrationPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"register.html",
		gin.H{
			"title": "Register",
		},
	)
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
	storage.CreateUser(&user)

	// If the user is created, set the token in a cookie and log the user in
	token := username
	c.SetCookie("token", token, 3600, "", "", false, false)
	c.Set("is_logged_in", true)

	render(c, gin.H{"title": "Successful registration & Login"},
		"login-successful.html")
}

func showLoginPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func performLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	status, err := storage.IsUserValid(username, password)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if status {
		// TODO use secure token
		token := username
		c.SetCookie("token", token, 3600, "", "", false, false)
		c.Set("is_logged_in", true)

		render(c, gin.H{
			"title": "Successful Login"}, "login-successful.html")
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
