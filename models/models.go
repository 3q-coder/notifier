package models

import (
	"fmt"

	"github.com/DryginAlexander/notifier"
	"github.com/DryginAlexander/notifier/settings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	gormigrate "gopkg.in/gormigrate.v1"
)

type Storage struct {
	DB *gorm.DB
}

func NewStorage() Storage {
	var DBConnStr string
	var dialect string
	switch settings.DBDialect {
	case "sqlite":
		DBConnStr = settings.DBName
		dialect = "sqlite3"
	case "postgresql":
		DBConnStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			settings.DBHost,
			settings.DBPort,
			settings.DBUser,
			settings.DBPw,
			settings.DBName,
		)
		dialect = "postgres"
		// default:
		// 	return errors.New(fmt.Sprintf("unknown DBDialect %s", settings.DBDialect))
	}

	db, _ := gorm.Open(dialect, DBConnStr)
	// DB.SetLogger(logger)
	return Storage{
		DB: db,
	}
}

func (s *Storage) CloseDB() {
	s.DB.Close()
}

func (s *Storage) MigrateDB() error {
	m := gormigrate.New(s.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// inital migration
		{
			ID: "202007182355",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					gorm.Model
					notifier.User
				}
				type Notification struct {
					gorm.Model
					notifier.Notification
				}
				return tx.AutoMigrate(&User{}, &Notification{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("users", "notifications").Error
			},
		},
		// future migrations ...
	})
	return m.Migrate()
}
