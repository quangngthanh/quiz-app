package service

import (
	"encoding/json"
	"log"
	"quiz-app/internal/model"
	"quiz-app/internal/repository"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService interface {
	RegisterLeaderboardViewer(quizID string, conn *websocket.Conn)
	UnregisterLeaderboardViewer(quizID string, conn *websocket.Conn)
	HasLeaderboardViewers(quizID string) bool
	BroadcastLeaderboardUpdate(quizID string, leaderboard []model.LeaderboardEntry)
}

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	quizID string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	quizID     string
}

type webSocketService struct {
	hubs      map[string]*Hub
	hubsMutex sync.RWMutex
	redisRepo repository.RedisRepository
}

func NewWebSocketService(redisRepo repository.RedisRepository) WebSocketService {
	service := &webSocketService{
		hubs:      make(map[string]*Hub),
		redisRepo: redisRepo,
	}

	return service
}

func (s *webSocketService) getOrCreateHub(quizID string) *Hub {
	s.hubsMutex.Lock()
	defer s.hubsMutex.Unlock()

	if hub, exists := s.hubs[quizID]; exists {
		return hub
	}

	hub := &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		quizID:     quizID,
	}

	s.hubs[quizID] = hub
	go hub.run()

	return hub
}

func (s *webSocketService) RegisterLeaderboardViewer(quizID string, conn *websocket.Conn) {
	hub := s.getOrCreateHub(quizID)
	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		quizID: quizID,
	}

	hub.register <- client

	go client.writePump()
	go client.readPump(hub)
}

func (s *webSocketService) UnregisterLeaderboardViewer(quizID string, conn *websocket.Conn) {
	s.hubsMutex.RLock()
	hub, exists := s.hubs[quizID]
	s.hubsMutex.RUnlock()

	if !exists {
		return
	}

	// Find and unregister client
	for client := range hub.clients {
		if client.conn == conn {
			hub.unregister <- client
			break
		}
	}
}

func (s *webSocketService) HasLeaderboardViewers(quizID string) bool {
	s.hubsMutex.RLock()
	defer s.hubsMutex.RUnlock()

	hub, exists := s.hubs[quizID]
	if !exists {
		return false
	}

	return len(hub.clients) > 0
}

func (s *webSocketService) BroadcastLeaderboardUpdate(quizID string, leaderboard []model.LeaderboardEntry) {
	update := model.LeaderboardUpdate{
		Type:        "leaderboard_update",
		Leaderboard: leaderboard,
	}

	message, err := json.Marshal(update)
	if err != nil {
		log.Printf("Error marshaling leaderboard update: %v", err)
		return
	}

	s.hubsMutex.RLock()
	hub, exists := s.hubs[quizID]
	s.hubsMutex.RUnlock()

	if exists {
		hub.broadcast <- message
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered for quiz %s. Total: %d", h.quizID, len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered for quiz %s. Total: %d", h.quizID, len(h.clients))
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
