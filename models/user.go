package models

import (
	"github.com/DryginAlexander/notifier"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	notifier.User
}

func (s *Storage) IsUsernameAvailable(username string) (bool, error) {
	var user User
	if err := s.DB.First(&user, "username = ?", username).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (s *Storage) IsUserValid(username, password string) (bool, error) {
	var user User
	if err := s.DB.First(&user, "username = ?", username).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	if user.Password != password {
		return false, nil
	}
	return true, nil
}

func (s *Storage) CreateUser(_user *notifier.User) error {
	user := User{User: *_user}
	err := s.DB.Create(&user).Error
	return err
}

func (s *Storage) UsersNumber() (int, error) {
	var count int
	err := s.DB.Model(&User{}).Count(&count).Error
	return count, err
}
