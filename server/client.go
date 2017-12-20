package server;

import (
	"log"
	"time"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

var count = 0

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type Client struct {
	id int
	name string
	send chan []byte
	router *Router
	connection *websocket.Conn
}

func (client *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.connection.Close()
	}()
	for {
		select {
		case message, ok := <- client.send:
			if !ok {
				client.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := client.connection.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Println("[Faild]", "Message cloudn't be written")
				return
			}

			writer.Write(message)

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) read() {
	defer func () {
		client.router.unregister <- client
		client.connection.Close()
	}()
	client.connection.SetReadLimit(maxMessageSize)
	client.connection.SetReadDeadline(time.Now().Add(pongWait))
	client.connection.SetPongHandler(func (string) error {
		client.connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := client.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("[Error] %v", err)
			}
			break
		}

		client.router.broadcast <- message
	}
}

func handleUpgradeRequest(
	router *Router,
	w http.ResponseWriter,
	r *http.Request,
) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("[Upgrade faild]", err)
		return
	}

	count++
	client := &Client{
		id: count,
		name: "connection: " + strconv.Itoa(count),
		router: router,
		connection: connection,
		send: make(chan []byte, 256),
	}

	client.router.register <- client

	go client.read()
	go client.write()
}
