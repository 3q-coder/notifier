package web

import (
	"github.com/DryginAlexander/notifier"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine
var operator notifier.Operator
var storage notifier.Storage

func Init(stor notifier.Storage, oper notifier.Operator) {
	operator = oper
	storage = stor

	Router = gin.Default()
	Router.LoadHTMLGlob("./web/templates/*")
	initializeRoutes()

	Router.Run()
}
