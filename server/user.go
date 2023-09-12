package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dev-szymon/translate-chat/server/lib"
	"golang.org/x/net/websocket"
)

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Language    string `json:"language"`
	currentRoom *Room
	conn        *websocket.Conn
	messageCh   chan []byte
}

func (u *User) sendTranslation(broadcastTranscript *BroadcastPayload) {
	ctx := context.Background()
	p := &TranslatedMessagePayload{
		RoomId:     u.currentRoom.id,
		UserId:     broadcastTranscript.FromUser.Id,
		Transcript: broadcastTranscript.Transcript,
		Confidence: broadcastTranscript.Confidence,
	}

	if broadcastTranscript.FromUser.Language != u.Language {
		translation, err := lib.TranslateTranscript(ctx, broadcastTranscript.FromUser.Language, u.Language, broadcastTranscript.Transcript)
		if err != nil {
			fmt.Println("Error translating file.")
		}
		p.Translation = &translation.Translation
	}

	u.sendEvent(translatedMessage, p)
}

func (u *User) sendEvent(eventType string, payload interface{}) {
	p, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Something went wrong sending event to user: %v", err)
	}
	event, err := json.Marshal(&Event{
		Type:    eventType,
		Payload: p,
	})
	if err != nil {
		log.Fatalf("Something went wrong sending event to user: %v", err)
	}

	u.conn.Write(event)

}
