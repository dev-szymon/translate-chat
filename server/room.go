package main

import (
	"github.com/dev-szymon/translate-chat/server/lib"
	"golang.org/x/net/websocket"
)

type User struct {
	server      *Server
	currentRoom *Room
	lang        string
	conn        *websocket.Conn
	messageCh   chan []byte
}

type Room struct {
	id          string
	users       map[*User]bool
	broadcastCh chan []byte
	joinCh      chan *User
	leaveCh     chan *User
}

func NewRoom() *Room {
	id := lib.GeneratePseudoID()
	room := &Room{
		id:          id,
		users:       make(map[*User]bool),
		broadcastCh: make(chan []byte),
		joinCh:      make(chan *User),
		leaveCh:     make(chan *User),
	}
	go room.run()
	return room
}

func (r *Room) run() {
	for {
		select {
		case user := <-r.joinCh:
			r.users[user] = true
			user.currentRoom = r
		case msg := <-r.broadcastCh:
			for usr := range r.users {
				usr.conn.Write(msg)
			}
		case usr := <-r.leaveCh:
			usr.currentRoom = nil
			delete(r.users, usr)
		}
	}
}
