package internal

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// Hub manages active WebSocket connections keyed by userID.
type Hub struct {
	mu      sync.RWMutex
	clients map[string][]*websocket.Conn
}

// NewHub constructs an empty hub.
func NewHub() *Hub {
	return &Hub{clients: make(map[string][]*websocket.Conn)}
}

// Register adds a connection for a given user.
func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = append(h.clients[userID], conn)
}

// Unregister removes a specific connection for a user.
func (h *Hub) Unregister(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.clients[userID]
	remaining := conns[:0]
	for _, c := range conns {
		if c != conn {
			remaining = append(remaining, c)
		}
	}
	h.clients[userID] = remaining
}

// Broadcast sends a JSON payload to all connections belonging to userID.
func (h *Hub) Broadcast(userID string, payload any) {
	h.mu.RLock()
	conns := make([]*websocket.Conn, len(h.clients[userID]))
	copy(conns, h.clients[userID])
	h.mu.RUnlock()

	for _, c := range conns {
		if err := c.WriteJSON(payload); err != nil {
			log.Warn().Err(err).Str("user_id", userID).Msg("ws write failed")
		}
	}
}
