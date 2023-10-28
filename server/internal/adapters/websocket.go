package adapters

import (
	"sync"

	"github.com/dev-szymon/translate-chat/server/internal/core/domain"
	"github.com/gorilla/websocket"
)

type Hub struct {
	rooms map[string]*domain.Room
}

type Room struct {
	clientsMu sync.Mutex
	clients   map[*Client]bool

	broadcastCh chan *domain.Message
}

type Client struct {
	conn        *websocket.Conn
	user        *domain.User
	currentRoom *domain.Room
}
