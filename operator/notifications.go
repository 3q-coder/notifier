package operator

import (
	"encoding/json"

	"github.com/DryginAlexander/notifier"
	"github.com/gorilla/websocket"
)

func (o *Operator) SubscribeToNotifications(username string, client *websocket.Conn) {
	// init channel
	o.noteChanals[username] = make(chan string, 10)

	// TODO load messages from db

	inChan := o.noteChanals[username]
	var msg string
	for {
		msg = <-inChan
		data, _ := json.Marshal(map[string]string{"message": msg})

		w, err := client.NextWriter(websocket.TextMessage)
		if err != nil {
			break
		}
		w.Write(data)
		w.Close()
	}
}

func (o *Operator) SendNotification(note notifier.Notification) {
	// TODO handle user is not exist
	outChan := o.noteChanals[note.Username]
	outChan <- note.Message
}

func (o *Operator) SendNotificationAll(message string) {
	for _, outChan := range o.noteChanals {
		outChan <- message
	}
}
