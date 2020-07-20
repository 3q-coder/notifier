package models

import (
	"github.com/DryginAlexander/notifier"
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	notifier.Notification
}

func (s *Storage) CreateNotification(_note *notifier.Notification) (uint, error) {
	note := Notification{Notification: *_note}
	err := s.db.Create(&note).Error
	return note.ID, err
}

func (s *Storage) CreateNotificationAll(message string) error {
	var users []*User
	err := s.db.Set("gorm:auto_preload", true).Find(&users).Error
	if err != nil {
		return err
	}

	for _, user := range users {
		note := notifier.Notification{
			Username: user.User.Username,
			Message:  message,
			Sent:     false,
		}
		_, err = s.CreateNotification(&note)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) SetSentNoteStatus(id uint) error {
	var note Notification
	err := s.db.Set("gorm:auto_preload", true).
		Where("id = ?", id).First(&note).Error
	if err != nil {
		return err
	}
	note.Notification.Sent = true
	return s.db.Save(note).Error
}

func (s *Storage) NotificationsByUsername(name string) ([]notifier.Notification, []uint, error) {
	notes := []*Notification{}
	err := s.db.Set("gorm:auto_preload", true).
		Where("username = ?", name).Where("sent = ?", false).Find(&notes).Error
	var notifier_notes []notifier.Notification
	var ids []uint
	for _, note := range notes {
		notifier_notes = append(notifier_notes, note.Notification)
		ids = append(ids, note.ID)
	}
	return notifier_notes, ids, err
}

func (s *Storage) NotesNumber() (int, error) {
	var count int
	err := s.db.Table("notifications").Count(&count).Error
	return count, err
}

func (s *Storage) SentNotesNumber() (int, error) {
	var count int
	err := s.db.Model(&Notification{}).Where("sent = ?", true).
		Count(&count).Error
	return count, err
}
