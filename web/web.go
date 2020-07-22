package web

import (
	"github.com/DryginAlexander/notifier"
	"github.com/DryginAlexander/notifier/settings"
	"github.com/gin-gonic/gin"
)

var operator Operator
var storage notifier.Storage

func Init(stor notifier.Storage) *gin.Engine {
	storage = stor
	operator = NewOperator(stor)

	router := gin.Default()
	router.LoadHTMLGlob(settings.StaticPath)
	initializeRoutes(router)

	return router
}
