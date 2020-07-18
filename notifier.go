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
}

type Storage interface {
	// user
	IsUsernameAvailable(username string) bool
	IsUserValid(username, password string) bool
	CreateUser(user *User) error
	// notification
	CreateNotification(note *Notification) error
	NotificationsByUsername(name string) ([]Notification, error)
	DeleteNotification(note *Notification) error
}

type Operator interface {
	// web
	SendNotification(client *websocket.Conn)
	RegisterUser(username, password string) (*User, error)
}
