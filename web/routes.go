package web

func initializeRoutes() {

	Router.Use(setUserStatus())

	Router.GET("/", showIndexPage)
	Router.GET("/news", ensureLoggedIn(), showNewsPage)
	Router.GET("/notifications", ensureLoggedIn(), initNoteSocket)

	Router.POST("/send-note", ensureApiKey(), sendNote)
	Router.POST("/send-note-all", ensureApiKey(), sendNoteAll)

	userRoutes := Router.Group("/u")
	{
		userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)
		userRoutes.POST("/register", ensureNotLoggedIn(), register)
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)
		userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)
		userRoutes.GET("/logout", ensureLoggedIn(), logout)
	}
}
