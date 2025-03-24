package ws

import (
	"car"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Manager struct {
	clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
	upgrader  websocket.Upgrader
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			// For testing; restrict origins in production.
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (m *Manager) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	m.clientsMu.Lock()
	m.clients[conn] = true
	m.clientsMu.Unlock()
	log.Println("new websocket client connected")

	for {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			m.clientsMu.Lock()
			delete(m.clients, conn)
			m.clientsMu.Unlock()
			log.Println("WebSocket client disconnected")
			break
		}
	}
}

func (m *Manager) BroadcastTelemetry(cars []*car.Car) {
	data, err := json.Marshal(cars)
	if err != nil {
		log.Println("error serializing data:", err)
		return
	}

	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()
	for conn := range m.clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Error writing to client:", err)
			conn.Close()
			delete(m.clients, conn)
		}
	}
}

func (m *Manager) StartBroadcast(getCars func() []*car.Car) {
	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		cars := getCars()
		m.BroadcastTelemetry(cars)
	}
}
