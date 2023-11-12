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

type broadcastTranslation struct {
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

	broadcastTranslationCh chan *broadcastTranslation
	broadcastEventCh       chan *OutboundEvent

	joinCh  chan *Client
	leaveCh chan *Client
}

func NewPool(ts ports.TranslateServicePort) *Pool {
	pool := &Pool{
		clients: make(map[*Client]bool),
		room:    chat.NewRoom(),

		broadcastTranslationCh: make(chan *broadcastTranslation),
		broadcastEventCh:       make(chan *OutboundEvent),
		joinCh:                 make(chan *Client),
		leaveCh:                make(chan *Client),
	}

	go pool.Read()
	go pool.Broadcast(ts)
	return pool
}

func (p *Pool) handleAddClient(client *Client) {
	p.clientsMu.Lock()
	p.clients[client] = true
	p.clientsMu.Unlock()

	p.room.AddUser(client.user)

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
}

func (p *Pool) handleRemoveClient(client *Client) {
	p.clientsMu.Lock()
	p.room.RemoveUser(client.user.Id)
	delete(p.clients, client)
	p.clientsMu.Unlock()

}

func (p *Pool) Read() {
	for {
		select {
		case client := <-p.joinCh:
			p.handleAddClient(client)
		case client := <-p.leaveCh:
			p.handleRemoveClient(client)
		}
	}
}

func (p *Pool) Broadcast(ts ports.TranslateServicePort) {
	for b := range p.broadcastTranslationCh {
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
		wg.Add(len(languages) - 1)
		translationsCh := make(chan *chat.Translation, len(languages)-1)

		for l := range languages {
			if l != transcript.SourceLang {
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
		}

		wg.Wait()
		close(translationsCh)

		translations := make(map[string]*chat.Translation)
		for translation := range translationsCh {
			translations[translation.TargetLang] = translation
		}

		for client := range p.clients {
			var translationText *string
			existingTranslation, ok := translations[client.user.Language]
			if ok {
				translationText = &existingTranslation.Text
			}

			client.conn.WriteJSON(&OutboundEvent{
				Type: NEW_MESSAGE_EVENT,
				Payload: &NewMessagePayload{
					Message: &Message{
						Id:          uuid.NewString(),
						Transcript:  transcript.Text,
						Confidence:  transcript.Confidence,
						Translation: translationText,
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
