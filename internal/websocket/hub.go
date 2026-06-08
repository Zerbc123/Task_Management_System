package websocket

import (
	"encoding/json"
	"log"
	"sync"

	gorilla "github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*gorilla.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*gorilla.Conn]bool),
	}
}

func (h *Hub) AddClient(conn *gorilla.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = true
	log.Println("websocket client connected")
}

func (h *Hub) RemoveClient(conn *gorilla.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, conn)
	conn.Close()

	log.Println("websocket client disconnected")
}

func (h *Hub) Broadcast(message string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteMessage(gorilla.TextMessage, []byte(message)); err != nil {
			log.Println("failed to send websocket message", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}

func (h *Hub) BroadcastJSON(event interface{}) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Println("failed to marshal websocket event:", err)
		return
	}

	h.Broadcast(string(data))
}
