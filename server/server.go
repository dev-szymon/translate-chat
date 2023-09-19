package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/dev-szymon/translate-chat/server/lib"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Server struct {
	id      string
	roomsMu sync.Mutex
	rooms   map[string]*Room
	usersMu sync.Mutex
	users   map[*User]bool
}

func NewServer() *Server {
	return &Server{
		id:    uuid.NewString(),
		rooms: make(map[string]*Room),
		users: make(map[*User]bool),
	}
}

type SuccessResponse struct {
	Message string `json:"message"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleTranslateFile() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := context.Background()
		if r.Method == "POST" {
			encoder := json.NewEncoder(w)
			userId := r.FormValue("userId")
			if userId == "" {
				fmt.Println("id missing")
				encoder.Encode(&ErrorResponse{Error: "user not found"})
				return
			}
			var user *User

			for u := range s.users {
				if u.Id == userId {
					user = u
					break
				}
			}

			if user == nil {
				encoder.Encode(&ErrorResponse{Error: "user not found"})
				return
			}

			file, _, err := r.FormFile("file")
			if err != nil {
				fmt.Println("Failed to parse form file:", err)
				encoder.Encode(&ErrorResponse{Error: "error parsing file"})
				return
			}
			defer file.Close()

			flacFile, err := lib.ConvertFileToFlac(file)
			if err != nil {
				fmt.Println("Failed to convert audio file", err)
				encoder.Encode(&ErrorResponse{Error: "error"})
				return
			}

			transcript, err := lib.TranscribeFlacFile(ctx, user.Language, flacFile)
			if err != nil {
				fmt.Println("failed to translate audio file", err)
				encoder.Encode(&ErrorResponse{Error: "error"})
				return
			}

			user.currentRoom.broadcastCh <- &BroadcastPayload{
				TranscriptResponse: transcript,
				sender:             user,
			}

			encoder.Encode(&SuccessResponse{Message: "ok"})
		}
	}
}

func (s *Server) serveWS(conn *websocket.Conn) {
	fmt.Println("New connection from: ", conn.RemoteAddr())

	user := s.createUser(conn)

	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			user.sendEvent(errorEvent, &ErrorPayload{Message: "Unable to read from buffer", Error: err.Error()})
			continue
		}
		msgBytes := buf[:n]

		var event Event
		err = json.Unmarshal(msgBytes, &event)
		if err != nil {
			user.sendEvent(errorEvent, &ErrorPayload{Message: "Unable to unmarshal event payload", Error: err.Error()})
		}

		switch event.Type {
		case joinRoom:
			var payload JoinRoomPayload
			var room *Room
			err := json.Unmarshal(event.Payload, &payload)
			if err != nil {
				user.sendEvent(errorEvent, &ErrorPayload{Message: "Unable to unmarshal join-room payload", Error: err.Error()})
			}

			room = s.rooms[payload.RoomId]
			if room == nil {
				for _, r := range s.rooms {
					if r.name == payload.RoomId {
						room = r
						break
					}
				}
			}
			if room == nil {
				room = s.createRoom()
			}

			user.Language = payload.Language
			user.Username = payload.Username
			room.joinCh <- user
		case leaveRoom:
			user.currentRoom.leaveCh <- user
		}
	}
}

func (s *Server) createUser(conn *websocket.Conn) *User {
	s.usersMu.Lock()
	defer s.usersMu.Unlock()

	user := &User{
		Id:          uuid.NewString(),
		conn:        conn,
		currentRoom: nil,
		messageCh:   make(chan []byte, 512),
	}
	s.users[user] = true

	return user
}

func (s *Server) createRoom() *Room {
	s.roomsMu.Lock()
	defer s.roomsMu.Unlock()
	room := &Room{
		id:          uuid.NewString(),
		name:        lib.GenerateRoomName(),
		users:       make(map[*User]bool),
		broadcastCh: make(chan *BroadcastPayload),
		joinCh:      make(chan *User),
		leaveCh:     make(chan *User),
	}

	go room.run()
	s.rooms[room.id] = room

	return room
}
