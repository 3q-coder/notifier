package notifier

type User struct {
	Username string
	Password string
}

type Notification struct {
	Username string `json:",omitempty"`
	Message  string
	Sent     bool `json:"-"`
}

type Storage interface {
	// user
	IsUsernameAvailable(username string) (bool, error)
	IsUserValid(username, password string) (bool, error)
	CreateUser(user *User) error
	// notification
	CreateNotification(note *Notification) (uint, error)
	CreateNotificationAll(message string) error
	SetSentNoteStatus(id uint) error
	NotificationsByUsername(name string) ([]Notification, []uint, error)
	// metrics
	UsersNumber() (int, error)
	NotesNumber() (int, error)
	SentNotesNumber() (int, error)
}
