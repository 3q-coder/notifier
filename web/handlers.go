package web

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/DryginAlexander/notifier"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func showIndexPage(c *gin.Context) {
	render(c, gin.H{"title": "Home Page"}, "index.html")
}

func showNewsPage(c *gin.Context) {
	render(c, gin.H{"title": "News Page"}, "news.html")
}

func initNoteSocket(c *gin.Context) {
	username := c.Query("token")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	go operator.SubscribeToNotifications(username, ws)
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

	var err error
	if strings.TrimSpace(password) == "" {
		err = errors.New("The password can't be empty")
	} else if !storage.IsUsernameAvailable(username) {
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

	operator.InitNewsChanel(username)
}

func showLoginPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func performLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if storage.IsUserValid(username, password) {
		// TODO use secure token
		token := username
		c.SetCookie("token", token, 3600, "", "", false, false)
		c.Set("is_logged_in", true)

		render(c, gin.H{
			"title": "Successful Login"}, "login-successful.html")

		operator.InitNewsChanel(username)
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

func sendNote(c *gin.Context) {
	username := c.PostForm("username")
	msg := c.PostForm("message")
	note := notifier.Notification{
		Username: username,
		Message:  msg,
	}
	operator.SendNotification(note)
}

func sendNoteAll(c *gin.Context) {
	msg := c.PostForm("message")
	operator.SendNotificationAll(msg)
}
