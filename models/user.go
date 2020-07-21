package models

import (
	"github.com/DryginAlexander/notifier"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	notifier.User
}

func (s *Storage) IsUsernameAvailable(username string) bool {
	var user User
	if err := s.DB.First(&user, "username = ?", username).Error; err != nil {
		return true
	}
	return false
}

func (s *Storage) IsUserValid(username, password string) bool {
	var user User
	if err := s.DB.First(&user, "username = ?", username).Error; err != nil {
		return false
	}
	if user.Password != password {
		return false
	}
	return true
}

func (s *Storage) CreateUser(_user *notifier.User) error {
	user := User{User: *_user}
	err := s.DB.Create(&user).Error
	return err
}

func (s *Storage) UsersNumber() (int, error) {
	var count int
	err := s.DB.Table("users").Count(&count).Error
	return count, err
}
