package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dev-szymon/translate-chat/server/lib"
	"github.com/google/uuid"
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

func (u *User) sendNewMessage(bp *BroadcastPayload) {
	ctx := context.Background()
	m := &Message{
		Id:         uuid.NewString(),
		SenderId:   bp.sender.Id,
		Transcript: bp.Transcript,
		Confidence: bp.Confidence,
	}

	if bp.sender.Language != u.Language {
		translation, err := lib.TranslateTranscript(ctx, bp.sender.Language, u.Language, bp.Transcript)
		if err != nil {
			fmt.Println("Error translating file: ", err)
			u.sendEvent(errorEvent, &ErrorPayload{Message: "Error while translating message", Error: err.Error()})
			return
		}
		m.Translation = &translation.Translation
	}

	p := &NewMessagePayload{
		Message: m,
	}

	u.sendEvent(newMessageEvent, p)
}

func parseEvent(eventType string, payload interface{}) (*Event, error) {
	switch eventType {
	case userJoinedEvent:
		p, ok := payload.(*UserJoinedPayload)
		if !ok {
			return nil, fmt.Errorf("error parsing event %s with payload: %v", eventType, payload)
		}
		b, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}
		e := &Event{
			Type:    eventType,
			Payload: b}
		return e, nil
	case newMessageEvent:
		p, ok := payload.(*NewMessagePayload)
		if !ok {
			return nil, fmt.Errorf("error parsing event %s with payload: %v", eventType, payload)
		}
		b, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}
		e := &Event{
			Type:    eventType,
			Payload: b}
		return e, nil
	case errorEvent:
		p, ok := payload.(*ErrorPayload)
		if !ok {
			return nil, fmt.Errorf("error parsing event %s with payload: %v", eventType, payload)
		}
		b, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}
		e := &Event{
			Type:    eventType,
			Payload: b}
		return e, nil
	default:
		return nil, fmt.Errorf("eventType not recognised: %v", eventType)
	}
}

func (u *User) sendEvent(eventType string, payload interface{}) {
	event, err := parseEvent(eventType, payload)
	if err != nil {
		log.Fatalf("Someting went wrong parsing event: %v", err)
	}

	b, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Something went wrong sending event to user: %v", err)
	}

	u.conn.Write(b)
}
