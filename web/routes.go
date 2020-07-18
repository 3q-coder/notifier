package web

func initializeRoutes() {

	Router.Use(setUserStatus())

	Router.GET("/", showIndexPage)
	Router.GET("/news", ensureLoggedIn(), showNewsPage)
	Router.GET("/notifications", initNoteSocket)

	userRoutes := Router.Group("/u")
	{
		userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)
		userRoutes.POST("/register", ensureNotLoggedIn(), register)
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)
		userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)
		userRoutes.GET("/logout", ensureLoggedIn(), logout)
	}
}
