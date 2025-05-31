package handler

import (
	"log"
	"net/http"
	"quiz-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type WebSocketHandler struct {
	wsService service.WebSocketService
}

func NewWebSocketHandler(wsService service.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{wsService: wsService}
}

func (h *WebSocketHandler) HandleLeaderboardWebSocket(c *gin.Context) {
	quizID := c.Param("quiz_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	h.wsService.RegisterLeaderboardViewer(quizID, conn)
}
