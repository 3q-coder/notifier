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

	// TODO rewrite - for not delete own client
	time.Sleep(3 * time.Second)

	// init channel
	ctx, finish := context.WithCancel(context.Background())
	inChan := make(chan ChanMessage, 10)
	o.clients[username] = Client{
		username: username,
		channel:  inChan,
		finish:   finish,
	}

	// load unsent message to channel
	notes, ids, _ := o.storage.NotificationsByUsername(username)
	for i, note := range notes {
		inChan <- ChanMessage{
			text: note.Message,
			id:   ids[i],
		}
	}

	go reader(finish, ws)
	go writer(ctx, inChan, ws, o.storage)

	// remove connection
	<-ctx.Done()
	delete(o.clients, username)
}

func reader(finish func(), ws *websocket.Conn) {
	defer func() {
		finish()
		ws.Close()
	}()
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

func writer(ctx context.Context, inChan <-chan ChanMessage, ws *websocket.Conn, storage notifier.Storage) {
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

			data, _ := json.Marshal(map[string]string{"message": msg.text})

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

			storage.SetSentNoteStatus(msg.id)
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(settings.WriteWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (o *Operator) SendNotification(note *notifier.Notification) {
	id, _ := o.storage.CreateNotification(note)

	// send message if client online
	if client, ok := o.clients[note.Username]; ok == true {
		client.channel <- ChanMessage{
			text: note.Message,
			id:   id,
		}

	}
}

func (o *Operator) SendNotificationAll(message string) {
	o.storage.CreateNotificationAll(message)

	// send message to all online clients
	for _, client := range o.clients {
		notes, ids, _ := o.storage.NotificationsByUsername(client.username)
		for i, note := range notes {
			client.channel <- ChanMessage{
				text: note.Message,
				id:   ids[i],
			}
		}
	}
}

func (o *Operator) OnlineClientsNumber() int {
	return len(o.clients)
}
