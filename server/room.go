package main

import (
	"fmt"
	"sync"

	"github.com/dev-szymon/translate-chat/server/lib"
)

type BroadcastPayload struct {
	*lib.TranscriptResponse

	FromUser *User `json:"from_user"`
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
		case user := <-r.joinCh:
			r.usersMu.Lock()
			r.users[user] = true
			r.usersMu.Unlock()
			user.currentRoom = r
			usersInRoom := []User{}
			for roomUser := range r.users {
				usersInRoom = append(usersInRoom, *roomUser)
			}

			for userAwaitingMessage := range r.users {
				userAwaitingMessage.sendEvent(userJoined, &UserJoinedRoomPayload{
					UserId:   user.Id,
					Username: user.Username,
					Language: user.Language,
					RoomId:   r.id,
					RoomName: r.name,
					Users:    usersInRoom,
				})
			}

		case transcript := <-r.broadcastCh:
			for user := range r.users {
				fmt.Println(user.Id)
				go user.sendTranslation(transcript)
			}
		case user := <-r.leaveCh:
			r.usersMu.Lock()
			user.currentRoom = nil
			delete(r.users, user)
			r.usersMu.Unlock()
		}
	}
}
