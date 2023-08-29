package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dev-szymon/translate-chat/server/lib"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type apiError struct {
	err    string
	status int
}

func (e *apiError) Error() string {
	return e.err
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

type Server struct {
	id         string
	rooms      map[string]*Room
	users      map[*User]bool
	register   chan *User
	unregister chan *User
}

func NewServer() *Server {
	return &Server{
		id:         uuid.NewString(),
		rooms:      make(map[string]*Room),
		users:      make(map[*User]bool),
		register:   make(chan *User),
		unregister: make(chan *User),
	}
}

func (s *Server) handleTranslateFile() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := context.Background()
		if r.Method == "POST" {
			userId := r.FormValue("userId")
			if userId == "" {
				WriteJSON(w, http.StatusUnprocessableEntity, &apiError{
					err:    "UserId missing",
					status: http.StatusUnprocessableEntity,
				})
			}
			var user *User

			for u := range s.users {
				if u.Id == userId {
					user = u
					break
				}
			}

			file, _, err := r.FormFile("file")
			if err != nil {
				fmt.Println("Failed to parse form file:", err)
				WriteJSON(w, http.StatusInternalServerError, &apiError{
					err:    "Error processing audio file.",
					status: http.StatusInternalServerError,
				})
			}
			defer file.Close()

			flacFile, err := lib.ConvertFileToFlac(file)
			if err != nil {
				fmt.Println("Failed to convert audio file", err)
				WriteJSON(w, http.StatusInternalServerError, &apiError{
					err:    "Error processing audio file.",
					status: http.StatusInternalServerError,
				})
			}

			transcript, err := lib.TranscribeFlacFile(ctx, user.Language, flacFile)
			if err != nil {
				fmt.Println("failed to translate audio file", err)
				WriteJSON(w, http.StatusInternalServerError, &apiError{
					err:    "Error transcribing audio file.",
					status: http.StatusInternalServerError,
				})
			}

			user.currentRoom.broadcastCh <- &BroadcastTrasncript{Transcript: transcript, FromUser: user}

			WriteJSON(w, http.StatusOK, transcript)
		}
	}
}

func (s *Server) serveWS(conn *websocket.Conn) {
	fmt.Println("New connection from: ", conn.RemoteAddr())
	user := &User{
		Id:          uuid.NewString(),
		conn:        conn,
		currentRoom: nil,
		messageCh:   make(chan []byte, 512),
	}
	s.users[user] = true

	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			conn.Write([]byte("wrong message format"))
			continue
		}
		msgBytes := buf[:n]

		var msg Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			conn.Write([]byte("wrong message format"))
		}

		switch msg.Type {
		case joinRoom:
			var room *Room
			if msg.RoomID != "" {
				for id := range s.rooms {
					if id == msg.RoomID {
						room = s.rooms[id]
						break
					}
				}
			}
			if room == nil {
				room = s.createRoom()
			}
			user.Language = msg.Language
			user.Username = msg.Username
			room.joinCh <- user
		case leaveRoom:
			user.currentRoom.leaveCh <- user
		case updateLanguage:
			if msg.Language == "" {
				conn.Write([]byte("Language not found"))
			}
			user.Language = msg.Language //TODO validate
		}
	}
}

func (s *Server) createRoom() *Room {
	room := &Room{
		id:          uuid.NewString(),
		name:        lib.GenerateRoomName(),
		users:       make(map[*User]bool),
		broadcastCh: make(chan *BroadcastTrasncript),
		joinCh:      make(chan *User),
		leaveCh:     make(chan *User),
	}

	go room.run()
	s.rooms[room.id] = room

	return room
}
