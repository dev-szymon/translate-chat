package chat

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type Room struct {
	Id      string
	Name    string
	usersMu sync.Mutex
	Users   map[string]*User
}

func NewRoom() *Room {
	room := &Room{
		Id:    uuid.NewString(),
		Name:  generateRoomName(),
		Users: make(map[string]*User),
	}
	return room
}

func (r *Room) AddUser(u *User) {
	r.usersMu.Lock()
	defer r.usersMu.Unlock()
	u.CurrentRoom = r
	r.Users[u.Id] = u
}

func (r *Room) RemoveUser(id string) {
	r.usersMu.Lock()
	defer r.usersMu.Unlock()
	delete(r.Users, id)
}

func generateRoomName() string {
	adjectives := []string{
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

	nouns := []string{
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

	a1 := rand.Intn(len(adjectives) - 1)
	a2 := rand.Intn(len(adjectives) - 1)
	n := rand.Intn(len(nouns) - 1)
	return fmt.Sprintf("%s-%s-%s", adjectives[a1], adjectives[a2], nouns[n])
}
