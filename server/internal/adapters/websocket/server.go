package websocket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/dev-szymon/translate-chat/server/internal/core/chat"
	"github.com/dev-szymon/translate-chat/server/internal/ports"
	"github.com/gorilla/websocket"
)

type Server struct {
	roomIdToPoolMu sync.Mutex
	roomIdToPool   map[string]*Pool
}

func (s *Server) AddConnectionPool(pool *Pool) {
	s.roomIdToPoolMu.Lock()
	defer s.roomIdToPoolMu.Unlock()
	s.roomIdToPool[pool.room.Id] = pool
}

func NewServer() *Server {
	return &Server{roomIdToPool: make(map[string]*Pool)}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const MAGIC_BYTE = "B"

func (s *Server) HandleWS(ts ports.TranslateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		encoder := json.NewEncoder(w)
		if err != nil {
			encoder.Encode(&struct{ Error string }{Error: "room id missing"})
		}
		fmt.Println("New connection from: ", ws.RemoteAddr())

		user := chat.NewUser()

		client := &Client{
			conn: ws,
			user: user,
		}

		var pool *Pool

		for {
			_, msgBytes, err := ws.ReadMessage()
			if err != nil {
				if err == io.EOF {
					break
				}
				ws.WriteJSON(&OutboundEvent{
					Type: ERROR_EVENT, Payload: &ErrorPayload{Error: "Unable to read from buffer"},
				})
				continue
			}

			if string(msgBytes[0]) == MAGIC_BYTE {
				audioBytes := msgBytes[1:]

				pool := s.roomIdToPool[user.CurrentRoom.Id]
				pool.broadcastCh <- &broadcastMessage{file: audioBytes, sender: user}
			} else {
				var event InboundEvent
				err = json.Unmarshal(msgBytes, &event)
				if err != nil {
					ws.WriteJSON(&OutboundEvent{
						Type: ERROR_EVENT, Payload: &ErrorPayload{Error: "Unable to unmarshal event payload"},
					})
				}

				switch event.Type {
				case JOIN_ROOM_EVENT:
					var payload JoinRoomPayload
					err := json.Unmarshal(event.Payload, &payload)
					if err != nil {
						ws.WriteJSON(&OutboundEvent{
							Type: ERROR_EVENT, Payload: &ErrorPayload{Error: "Unable to unmarshal join-room payload"},
						})
					}

					pool, ok := s.roomIdToPool[payload.RoomId]
					if !ok {
						for _, p := range s.roomIdToPool {
							if p.room.Name == payload.RoomId {
								pool = p
								break
							}
						}
						if pool == nil {
							pool = NewPool(ts)
							s.AddConnectionPool(pool)
						}
					}

					client.user.UpdateUserDetails(payload.Username, payload.Language)
					pool.joinCh <- client
				case LEAVE_ROOM_EVENT:
					pool.leaveCh <- client
				}
			}
		}
	}
}
