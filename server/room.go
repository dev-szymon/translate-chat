package main

import (
	"fmt"
)

type Room struct {
	id          string
	name        string
	users       map[*User]bool
	broadcastCh chan *BroadcastTrasncript
	joinCh      chan *User
	leaveCh     chan *User
}

func (r *Room) run() {
	for {
		select {
		case user := <-r.joinCh:
			r.users[user] = true
			user.currentRoom = r
			user.conn.Write([]byte(fmt.Sprintf("roomId:%s", r.id)))
		case transcript := <-r.broadcastCh:
			for user := range r.users {
				go user.sendTranslation(transcript)
			}
		case user := <-r.leaveCh:
			user.currentRoom = nil
			delete(r.users, user)
		}
	}
}
