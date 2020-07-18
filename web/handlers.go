package web

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

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
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	go operator.SendNotification(ws)
}

func generateSessionToken() string {
	// TODO use secure way
	return strconv.FormatInt(rand.Int63(), 16)
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

	if _, err := operator.RegisterUser(username, password); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

		render(c, gin.H{"title": "Successful registration & Login"},
			"login-successful.html")

	} else {
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
		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
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
