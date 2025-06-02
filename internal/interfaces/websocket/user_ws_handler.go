package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hinha/echo-casbin-ddd-app/internal/application/dto"
	"github.com/hinha/echo-casbin-ddd-app/internal/application/interfaces"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"github.com/labstack/echo/v4"
)

// @title User WebSocket API
// @version 1.0
// @description WebSocket API for real-time user management
// @BasePath /ws

// UserWSHandler handles WebSocket connections for user management
type UserWSHandler struct {
	userUseCase interfaces.UserUseCase
	clients     map[*websocket.Conn]bool
	broadcast   chan []byte
	register    chan *websocket.Conn
	unregister  chan *websocket.Conn
	shutdown    chan struct{}
	mutex       sync.Mutex
}

// NewUserWSHandler creates a new UserWSHandler
func NewUserWSHandler(userUseCase interfaces.UserUseCase) *UserWSHandler {
	return &UserWSHandler{
		userUseCase: userUseCase,
		clients:     make(map[*websocket.Conn]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
		shutdown:    make(chan struct{}),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// UserListMessage represents a user list message
type UserListMessage struct {
	Users []*UserResponse `json:"users"`
}

// UserResponse represents a user in the response
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// toUserResponse converts a user entity to a user response
func toUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Active:    user.Active,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Start starts the WebSocket handler
func (h *UserWSHandler) Start() {
	go h.run()
}

// run runs the WebSocket handler
func (h *UserWSHandler) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			h.sendUserList(client)
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("Error writing message: %v", err)
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		case <-h.shutdown:
			h.mutex.Lock()
			for client := range h.clients {
				client.Close()
				delete(h.clients, client)
			}
			h.mutex.Unlock()
			return
		}
	}
}

// sendUserList sends the user list to a client
func (h *UserWSHandler) sendUserList(client *websocket.Conn) {
	input := dto.ListUsersInput{
		Page:  1,
		Limit: 100,
	}

	output, err := h.userUseCase.ListUsers(context.Background(), input)
	if err != nil {
		log.Printf("Error getting user list: %v", err)
		return
	}

	users := make([]*UserResponse, len(output.Users))
	for i, user := range output.Users {
		users[i] = toUserResponse(user)
	}

	message := WSMessage{
		Type: "user_list",
		Payload: func() json.RawMessage {
			data, _ := json.Marshal(UserListMessage{Users: users})
			return data
		}(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Error writing message: %v", err)
		client.Close()
		h.mutex.Lock()
		delete(h.clients, client)
		h.mutex.Unlock()
	}
}

// BroadcastUserUpdate broadcasts a user update to all clients
func (h *UserWSHandler) BroadcastUserUpdate() {
	input := dto.ListUsersInput{
		Page:  1,
		Limit: 100,
	}

	output, err := h.userUseCase.ListUsers(context.Background(), input)
	if err != nil {
		log.Printf("Error getting user list: %v", err)
		return
	}

	users := make([]*UserResponse, len(output.Users))
	for i, user := range output.Users {
		users[i] = toUserResponse(user)
	}

	message := WSMessage{
		Type: "user_list",
		Payload: func() json.RawMessage {
			data, _ := json.Marshal(UserListMessage{Users: users})
			return data
		}(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.broadcast <- data
}

// HandleWebSocket handles WebSocket connections
// @Summary Connect to user WebSocket
// @Description Establish a WebSocket connection for real-time user updates
// @Tags websocket
// @Accept json
// @Produce json
// @Success 101 {string} string "Switching Protocols to WebSocket"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /users [get]
func (h *UserWSHandler) HandleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	h.register <- ws

	go func() {
		defer func() {
			h.unregister <- ws
		}()

		// Set up a channel to receive read errors
		errCh := make(chan error, 1)

		// Start a goroutine to read messages
		go func() {
			for {
				_, _, err := ws.ReadMessage()
				if err != nil {
					errCh <- err
					return
				}
			}
		}()

		// Wait for either an error or shutdown signal
		select {
		case err := <-errCh:
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
		case <-h.shutdown:
			// Server is shutting down, no need to log anything
		}
	}()

	return nil
}

// StartPeriodicUpdate starts periodic updates
func (h *UserWSHandler) StartPeriodicUpdate(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				h.BroadcastUserUpdate()
			case <-h.shutdown:
				return
			}
		}
	}()
}

// RegisterRoutes registers the WebSocket routes
func (h *UserWSHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/ws/users", h.HandleWebSocket)
}

// Stop gracefully shuts down the WebSocket handler
func (h *UserWSHandler) Stop() {
	close(h.shutdown)
}
