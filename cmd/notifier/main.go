package main

import (
	"fmt"

	"github.com/DryginAlexander/notifier/models"
	"github.com/DryginAlexander/notifier/settings"
	"github.com/DryginAlexander/notifier/web"
)

func main() {

	fmt.Println("init settings")
	_ = settings.Init("./settings/dev.env")

	fmt.Println("connecting to db")
	stor := models.NewStorage()
	defer stor.CloseDB()

	fmt.Println("applying migration if needed")
	_ = stor.MigrateDB()

	router := web.Init(&stor)
	router.Run()
}
