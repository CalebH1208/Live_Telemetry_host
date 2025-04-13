package ws

import (
	"car"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"serial"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Manager struct {
	clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
	writeMu   sync.Mutex
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

	go func() {
		defer func() {
			conn.Close()
			m.clientsMu.Lock()
			delete(m.clients, conn)
			m.clientsMu.Unlock()
			log.Println("WebSocket client disconnected")
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading WebSocket message:", err)
				return
			}

			// Process the message
			var msgData map[string]interface{}
			if err := json.Unmarshal(message, &msgData); err != nil {
				log.Println("Error parsing WebSocket message:", err)
				continue
			}

			// Process the parsed message
			m.handleWSMessage(conn, msgData)
		}
	}()
}

func (m *Manager) BroadcastTelemetry(cars []*car.Car) {
	message := map[string]any{
		"type": "telemetry",
		"cars": cars,
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("error serializing data:", err)
		return
	}
	if len(cars) == 0 {
		return
	}

	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()
	//log.Print("length", len(m.clients), " | ", len(data))
	for conn := range m.clients {
		m.writeMu.Lock()
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Error writing to client:", err)
			conn.Close()
			delete(m.clients, conn)
		} else {
			//log.Print("sending data")
		}
		m.writeMu.Unlock()
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

func (m *Manager) Send_available_ports() {
	for {
		ports, err := serial.Get_port_list()
		if err != nil {
			log.Fatal("serialport failed", err)
			return
		}

		message := map[string]any{
			"type":  "port_list",
			"ports": ports,
		}

		data, _ := json.Marshal(message)

		for conn := range m.clients {
			m.writeMu.Lock()
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Error writing to client:", err)
				conn.Close()
				delete(m.clients, conn)
			} else {
				//log.Print("sending data")
			}
			m.writeMu.Unlock()
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (m *Manager) Send_lap_times(lapTimesaddr *[1024]time.Duration) {
	for {
		lapTimes := *lapTimesaddr
		// Find highest index with data
		highestUsed := -1
		for i, lap := range lapTimes {
			if lap != 0 {
				highestUsed = i
			}
		}

		laps := make([]string, highestUsed+1)
		for i := 0; i <= highestUsed; i++ {
			if lapTimes[i] == 0 {
				laps[i] = "00:00.00"
			} else {
				minutes := int(lapTimes[i].Minutes())
				seconds := int((lapTimes[i] % time.Minute).Seconds())
				hundredths := int((lapTimes[i] % time.Second).Milliseconds() / 10)
				laps[i] = fmt.Sprintf("%02d:%02d.%02d", minutes, seconds, hundredths)
			}
		}

		message := map[string]any{
			"type":      "lap_times",
			"lap_times": laps,
		}

		data, _ := json.Marshal(message)

		for conn := range m.clients {
			m.writeMu.Lock()
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Error writing to client:", err)
				conn.Close()
				delete(m.clients, conn)
			}
			m.writeMu.Unlock()
		}

		time.Sleep(400 * time.Millisecond)
	}
}

func (m *Manager) handleWSMessage(conn *websocket.Conn, message map[string]interface{}) {
	msgType, ok := message["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "select_port":
		portName, _ := message["port"].(string)
		if portName != "" && portName != "no active ports" {
			log.Print("connected to : ", portName)
			go serial.Read_serial_message(portName)
		}
	case "reset":
		serial.Restart_serial()
	default:
		log.Print("handle WS message fail")

	}

}
