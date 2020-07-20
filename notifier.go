package notifier

// TODO move dependency from root module
import "github.com/gorilla/websocket"

type User struct {
	Username string
	Password string
}

type Notification struct {
	Username string
	Message  string
	Sent     bool
}

type Storage interface {
	// user
	IsUsernameAvailable(username string) bool
	IsUserValid(username, password string) bool
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

type Operator interface {
	// notification
	SubscribeToNotifications(username string, client *websocket.Conn)
	SendNotification(note *Notification)
	SendNotificationAll(message string)
	// metrics
	OnlineClientsNumber() int
}
