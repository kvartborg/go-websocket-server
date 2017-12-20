package server

import (
	"log"
)

type Router struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

func newRouter() *Router {
	return &Router{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (router *Router) run() {
	for {
		select {
		case client := <- router.register:
			log.Println("[Register]", client.id, client.name)
			router.clients[client] = true
		case client := <- router.unregister:
			log.Println("[Unregister]", client.id, client.name)
			if _, ok := router.clients[client]; ok {
				delete(router.clients, client)
				close(client.send)
			}
		case message := <- router.broadcast:
			log.Println("[Received]", string(message))
			for client := range router.clients {
				client.send <- message
			}
		}
	}
}
