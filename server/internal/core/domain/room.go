package domain

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type broadcastMessage struct {
	Transcript *Transcript
	Sender     *User
}

type Room struct {
	Id      string
	Name    string
	usersMu sync.Mutex
	Users   map[string]*User

	broadcastCh chan *broadcastMessage
	joinCh      chan *User
	leaveCh     chan *User
}

func NewRoom() *Room {
	room := &Room{
		Id:          uuid.NewString(),
		Name:        generateRoomName(),
		Users:       make(map[string]*User),
		broadcastCh: make(chan *broadcastMessage),
		joinCh:      make(chan *User),
		leaveCh:     make(chan *User),
	}
	go room.ReadPump()
	go room.WritePump()
	return room
}

func (r *Room) AddUser(u *User) {
	r.usersMu.Lock()
	defer r.usersMu.Unlock()
	u.CurrentRoom = r
	r.Users[u.Id] = u
	r.joinCh <- u
}

func (r *Room) ReadPump() {
	for {
		select {
		case newUser := <-r.joinCh:
			fmt.Println(newUser)
			usersInRoom := make([]*User, len(r.Users))
			for _, u := range r.Users {
				usersInRoom = append(usersInRoom, u)
			}

			for u := range r.Users {
				//TODO send event
				fmt.Println(u)
				// go u.sendEvent(userJoinedEvent, &UserJoinedPayload{
				// 	NewUser: newUser,
				// 	Room: &CurrentRoom{
				// 		Id:    r.Id,
				// 		Name:  r.Name,
				// 		Users: usersInRoom,
				// 	},
				// })
			}

		case user := <-r.leaveCh:
			r.usersMu.Lock()
			user.CurrentRoom = nil
			delete(r.Users, user.Id)
			r.usersMu.Unlock()
		}
	}
}

func (r *Room) WritePump() {
	for transcript := range r.broadcastCh {
		for user := range r.Users {
			//TODO
			fmt.Println(user, transcript)
			// go user.sendNewMessage(transcript)
		}
	}
}

func (r *Room) BroadcastMessage(transcript *Transcript, sender *User) {
	r.broadcastCh <- &broadcastMessage{Transcript: transcript, Sender: sender}
}

var adjectives = []string{
	"silly",
	"witty",
	"bubbly",
	"zany",
	"cheeky",
	"goofy",
	"wacky",
	"whimsical",
	"quirky",
	"hilarious",
	"bizarre",
	"absurd",
	"lively",
	"funky",
	"zesty",
	"jolly",
	"playful",
	"ridiculous",
	"quizzical",
	"eccentric",
}

var nouns = []string{
	"banana",
	"octopus",
	"tornado",
	"pajamas",
	"cupcake",
	"disco-ball",
	"narwhal",
	"pillow",
	"llama",
	"pickle",
	"bubblegum",
	"sock-puppet",
	"chicken",
	"sushi",
	"unicorn",
	"hotdog",
	"gummy-bear",
	"pineapple",
	"hippo",
	"marshmallow",
}

func generateRoomName() string {
	a1 := rand.Intn(len(adjectives) - 1)
	a2 := rand.Intn(len(adjectives) - 1)
	n := rand.Intn(len(nouns) - 1)
	return fmt.Sprintf("%s-%s-%s", adjectives[a1], adjectives[a2], nouns[n])
}
