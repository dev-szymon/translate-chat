package websocket

import (
	"context"
	"fmt"
	"sync"

	"github.com/dev-szymon/translate-chat/server/internal/core/chat"
	"github.com/dev-szymon/translate-chat/server/internal/ports"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type broadcastMessage struct {
	file   []byte
	sender *chat.User
}

type Client struct {
	conn *websocket.Conn
	user *chat.User
}
type Pool struct {
	clientsMu sync.Mutex
	clients   map[*Client]bool
	room      *chat.Room

	broadcastCh chan *broadcastMessage

	joinCh  chan *Client
	leaveCh chan *Client
}

func NewPool(ts ports.TranslateServicePort) *Pool {
	pool := &Pool{
		clients: make(map[*Client]bool),
		room:    chat.NewRoom(),

		broadcastCh: make(chan *broadcastMessage),
		joinCh:      make(chan *Client),
		leaveCh:     make(chan *Client),
	}

	go pool.Read()
	go pool.Broadcast(ts)
	return pool
}

func (p *Pool) addClient(client *Client) {
	p.clientsMu.Lock()
	defer p.clientsMu.Unlock()

	p.clients[client] = true
	p.room.AddUser(client.user)
}

func (p *Pool) removeClient(client *Client) {
	p.clientsMu.Lock()
	defer p.clientsMu.Unlock()

	p.room.RemoveUser(client.user.Id)
}

func (p *Pool) Read() {
	for {
		select {
		case client := <-p.joinCh:
			p.addClient(client)
			usersInRoom := []*RoomUser{}
			for c := range p.clients {
				usersInRoom = append(usersInRoom, &RoomUser{Id: c.user.Id, Username: c.user.Username, Language: c.user.Language})
			}
			payload := &UserJoinedPayload{
				NewUser: &RoomUser{Id: client.user.Id, Username: client.user.Username, Language: client.user.Language},
				Room:    &CurrentRoom{Id: p.room.Id, Name: p.room.Name, Users: usersInRoom},
			}
			for c := range p.clients {
				c.conn.WriteJSON(&OutboundEvent{
					Type:    USER_JOINED_EVENT,
					Payload: payload,
				})
			}

		case client := <-p.leaveCh:
			p.removeClient(client)
		}
	}
}

func (p *Pool) Broadcast(ts ports.TranslateServicePort) {
	for b := range p.broadcastCh {
		transcript, err := ts.TranscribeAudio(context.Background(), b.sender.Language, b.file)
		if err != nil {
			// TODO handle error
			fmt.Println("failed to translate audio file", err)
		}

		languages := make(map[string]bool)
		for client := range p.clients {
			languages[client.user.Language] = true
		}

		var wg sync.WaitGroup
		wg.Add(len(languages))
		translationsCh := make(chan *chat.Translation, len(languages))

		for l := range languages {
			go func(ch chan *chat.Translation, sourceLang, targetLang, text string) {
				translation, err := ts.TranslateText(context.Background(), sourceLang, targetLang, text)
				if err != nil {
					// TODO handle error
					fmt.Printf("Translation error: %+v", err)
				}
				wg.Done()
				ch <- translation
			}(translationsCh, transcript.SourceLang, l, transcript.Text)
		}

		wg.Wait()
		close(translationsCh)

		translations := make(map[string]*chat.Translation)
		for translation := range translationsCh {
			translations[translation.TargetLang] = translation
		}

		for client := range p.clients {
			client.conn.WriteJSON(&OutboundEvent{
				Type: NEW_MESSAGE_EVENT,
				Payload: &NewMessagePayload{
					Message: &Message{
						Id:          uuid.NewString(),
						Transcript:  transcript.Text,
						Confidence:  transcript.Confidence,
						Translation: &translations[client.user.Language].Text,
						Sender: &RoomUser{
							Id:       b.sender.Id,
							Username: b.sender.Username,
							Language: b.sender.Language,
						},
					},
				}})
		}
	}
}
