package operator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DryginAlexander/notifier"
	"github.com/DryginAlexander/notifier/settings"
	"github.com/gorilla/websocket"
)

func (o *Operator) SubscribeToNotifications(username string, ws *websocket.Conn) {
	// close existing connection
	if client, ok := o.clients[username]; ok == true {
		client.finish()
	}

	// init channel
	ctx, finish := context.WithCancel(context.Background())
	inChan := make(chan string, 10)
	o.clients[username] = Client{
		channel: inChan,
		finish:  finish,
	}

	// TODO load messages from db

	go reader(ws)
	go writer(ctx, inChan, ws)
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(settings.PongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(settings.PongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ctx context.Context, inChan <-chan string, ws *websocket.Conn) {
	pingTicker := time.NewTicker(settings.PingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-inChan:
			// TODO add msg respond

			data, _ := json.Marshal(map[string]string{"message": msg})

			ws.SetWriteDeadline(time.Now().Add(settings.WriteWait))
			w, err := ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(data)
			if err != nil {
				return
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(settings.WriteWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (o *Operator) SendNotification(note notifier.Notification) {
	// TODO handle user is not exist
	client := o.clients[note.Username]
	client.channel <- note.Message
}

func (o *Operator) SendNotificationAll(message string) {
	for _, client := range o.clients {
		client.channel <- message
	}
}
