package main

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

const magicByte = "f"

type Server struct {
	rooms      map[string]*Room
	users      map[*User]bool
	register   chan *User
	unregister chan *User
}

func NewServer() *Server {
	return &Server{
		rooms:      make(map[string]*Room),
		users:      make(map[*User]bool),
		register:   make(chan *User),
		unregister: make(chan *User),
	}
}

func (s *Server) serveWS(conn *websocket.Conn) {
	user := &User{
		server:      s,
		conn:        conn,
		currentRoom: nil,
		messageCh:   make(chan []byte, 256),
	}
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				conn.Write([]byte("file size exceeded"))
				break
			}
			conn.Write([]byte("wrong message format"))
			continue
		}
		msgBytes := buf[:n]

		if string(msgBytes[:1]) == magicByte {
			fmt.Println("magic byte available")
			// TODO handle case with magic byte
		}

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
				conn.Write([]byte(fmt.Sprintf("Created new room: %s", room.id)))
			}
			room.joinCh <- user
		case leaveRoom:
			user.currentRoom.leaveCh <- user
		case updateLang:
			if msg.Lang == "" {
				conn.Write([]byte("Language not found"))
			}
			user.lang = msg.Lang //TODO validate
		case message:
			user.currentRoom.broadcastCh <- []byte(msg.Value)
		}

		switch true {

		default:
			user.currentRoom.broadcastCh <- msgBytes
		}
	}
}

func (s *Server) createRoom() *Room {
	room := NewRoom()
	s.rooms[room.id] = room
	return room
}
