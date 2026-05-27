package internal

import (
	"errors"
	"net/http"

	pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/agamrai0123/wanderplan/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handlers wires HTTP requests to NotificationService methods.
type Handlers struct {
	svc *NotificationService
	hub *Hub
}

// NewHandlers constructs the handler set.
func NewHandlers(svc *NotificationService, hub *Hub) *Handlers { return &Handlers{svc: svc, hub: hub} }

// ListNotifications handles GET /notifications
func (h *Handlers) ListNotifications(c *gin.Context) {
	userID := middleware.UserID(c)
	ns, err := h.svc.GetByUser(c.Request.Context(), userID)
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, ns)
}

// MarkRead handles PATCH /notifications/:id/read
func (h *Handlers) MarkRead(c *gin.Context) {
	userID := middleware.UserID(c)
	err := h.svc.MarkRead(c.Request.Context(), c.Param("id"), userID)
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("notification not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.NoContent(c)
}

// MarkAllRead handles PATCH /notifications/read-all
func (h *Handlers) MarkAllRead(c *gin.Context) {
	userID := middleware.UserID(c)
	if err := h.svc.MarkAllRead(c.Request.Context(), userID); err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.NoContent(c)
}

// WebSocket handles GET /ws/notifications — upgrades to WS and registers connection.
func (h *Handlers) WebSocket(c *gin.Context) {
	userID := middleware.UserID(c)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.Err(c, pkgerr.Internal("websocket upgrade failed"))
		return
	}
	h.hub.Register(userID, conn)
	defer func() {
		h.hub.Unregister(userID, conn)
		conn.Close()
	}()
	// Keep alive: read until client closes.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
