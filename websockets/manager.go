package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,                                       // Size of the read buffer (in bytes)
	WriteBufferSize: 1024,                                       // Size of the write buffer (in bytes)
	CheckOrigin:     func(r *http.Request) bool { return true }, // Function to validate the origin of the request
}

type Manager struct {
	mu             sync.Mutex
	connections    map[string][]*websocket.Conn // Map of user IDs to their WebSocket connections
	allConnections []*websocket.Conn
}

func NewManager() *Manager {
	return &Manager{
		connections:    make(map[string][]*websocket.Conn),
		allConnections: []*websocket.Conn{},
	}
}

// NewManager creates a new WebSocket manager
func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Extract user ID from query parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("Missing user_id parameter")
		conn.Close()
		return
	}

	// Add the connection to the manager
	m.mu.Lock()
	if m.connections[userID] == nil {
		m.connections[userID] = []*websocket.Conn{}
	}
	m.connections[userID] = append(m.connections[userID], conn)
	m.allConnections = append(m.allConnections, conn) // Add to allConnections list
	m.mu.Unlock()

	log.Printf("User %s connected via WebSocket", userID)

	// Handle disconnection
	go func() {
		defer func() {
			m.mu.Lock()
			conns := m.connections[userID]
			for i, c := range conns {
				if c == conn {
					m.connections[userID] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			if len(m.connections[userID]) == 0 {
				delete(m.connections, userID)
			}

			// Remove the connection from allConnections
			for i, c := range m.allConnections {
				if c == conn {
					m.allConnections = append(m.allConnections[:i], m.allConnections[i+1:]...)
					break
				}
			}

			m.mu.Unlock()
			log.Printf("User %s disconnected from WebSocket", userID)
			conn.Close()
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				break
			}
		}
	}()
}

func (m *Manager) BroadcastToAll(message interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	for _, conn := range m.allConnections {
		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Printf("Failed to send message to WebSocket client: %v", err)
			// Remove faulty connections
			for i, c := range m.allConnections {
				if c == conn {
					m.allConnections = append(m.allConnections[:i], m.allConnections[i+1:]...)
					break
				}
			}
		}
	}
}
