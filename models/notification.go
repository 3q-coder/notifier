package models

import (
	"github.com/DryginAlexander/notifier"
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	notifier.Notification
}

func (s *Storage) CreateNotification(_note *notifier.Notification) error {
	note := Notification{Notification: *_note}

	if err := s.db.Create(&note).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) NotificationsByUsername(name string) ([]notifier.Notification, error) {
	notes := []*Notification{}
	err := s.db.Set("gorm:auto_preload", true).
		Where("username = ?", name).Find(&notes).Error
	var notifier_notes []notifier.Notification
	for _, note := range notes {
		notifier_notes = append(notifier_notes, note.Notification)
	}
	return notifier_notes, err
}

func (s *Storage) DeleteNotification(_note *notifier.Notification) error {
	note := Notification{Notification: *_note}

	if err := s.db.Delete(&note).Error; err != nil {
		return err
	}
	return nil
}
