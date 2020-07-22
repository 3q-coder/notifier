package main

import (
	"fmt"

	"github.com/DryginAlexander/notifier/models"
	"github.com/DryginAlexander/notifier/settings"
	"github.com/DryginAlexander/notifier/web"
)

func main() {

	fmt.Println("init settings")
	err := settings.Init("./settings/dev.env")
	if err != nil {
		fmt.Println("failed to init settings")
		return
	}

	fmt.Println("init db")
	stor, err := models.NewStorage()
	if err != nil {
		fmt.Println("failed to init db")
		return
	}
	defer stor.CloseDB()

	fmt.Println("applying migration if needed")
	err = stor.MigrateDB()
	if err != nil {
		fmt.Println("failed to migrate db")
		return
	}

	router := web.Init(&stor)
	router.Run()
}
