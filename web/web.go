package web

import (
	"github.com/DryginAlexander/notifier"
	"github.com/DryginAlexander/notifier/settings"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine
var operator Operator
var storage notifier.Storage

func Init(stor notifier.Storage) *gin.Engine {
	storage = stor
	operator = NewOperator(stor)

	Router = gin.Default()
	Router.LoadHTMLGlob(settings.StaticPath)
	initializeRoutes()

	return Router
}
