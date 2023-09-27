package main

import (
	"sync"

	"github.com/dev-szymon/translate-chat/server/lib"
)

type BroadcastPayload struct {
	*lib.TranscriptResponse
	sender *User
}
type Room struct {
	id          string
	name        string
	usersMu     sync.Mutex
	users       map[*User]bool
	broadcastCh chan *BroadcastPayload
	joinCh      chan *User
	leaveCh     chan *User
}

func (r *Room) run() {
	for {
		select {
		case newUser := <-r.joinCh:
			r.usersMu.Lock()
			r.users[newUser] = true
			r.usersMu.Unlock()
			newUser.currentRoom = r

			usersInRoom := []*User{}
			for u := range r.users {
				usersInRoom = append(usersInRoom, u)
			}

			for u := range r.users {
				go u.sendEvent(userJoinedEvent, &UserJoinedPayload{
					NewUser: newUser,
					Room: &CurrentRoom{
						Id:    r.id,
						Name:  r.name,
						Users: usersInRoom,
					},
				})
			}

		case transcript := <-r.broadcastCh:
			for user := range r.users {
				go user.sendNewMessage(transcript)
			}

		case user := <-r.leaveCh:
			r.usersMu.Lock()
			user.currentRoom = nil
			delete(r.users, user)
			r.usersMu.Unlock()
		}
	}
}
