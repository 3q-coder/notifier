package web

import (
	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {

	router.Use(setUserStatus())
	router.GET("/", showIndexPage)
	router.GET("/news", ensureLoggedIn(), showNewsPage)
	router.GET("/notifications", ensureLoggedIn(), initNoteSocket)

	router.POST("/send-note", ensureApiKey(), sendNote)
	router.POST("/send-note-all", ensureApiKey(), sendNoteAll)

	router.GET("/metrics", ensureApiKey(), metrics)

	userRoutes := router.Group("/u")
	{
		userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)
		userRoutes.POST("/register", ensureNotLoggedIn(), register)
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)
		userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)
		userRoutes.GET("/logout", ensureLoggedIn(), logout)
	}
}
