package adapters

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dev-szymon/translate-chat/server/internal/core/domain"
)

const (
	errorEvent      = "error"
	joinRoom        = "join-room"
	userJoinedEvent = "user-joined"
	newMessageEvent = "new-message"
	leaveRoom       = "leave-room"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Message struct {
	Id          string  `json:"id"`
	Transcript  string  `json:"transcript"`
	Confidence  float32 `json:"confidence"`
	Translation *string `json:"translation"`
	SenderId    string  `json:"senderId"`
}

type CurrentRoom struct {
	Id    string         `json:"id"`
	Name  string         `json:"name"`
	Users []*domain.User `json:"users"`
}

type UserJoinedPayload struct {
	NewUser *domain.User `json:"newUser"`
	Room    *CurrentRoom `json:"room"`
}

type NewMessagePayload struct {
	Message *Message `json:"message"`
}

type LeaveRoomPayload struct {
	RoomId string `json:"roomId"`
}
type UserLeftRoomMessage struct {
	UserId string `json:"userId"`
}

type ErrorPayload struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type JoinRoomPayload struct {
	Username string `json:"username"`
	Language string `json:"language"`
	RoomId   string `json:"roomId"`
}

type EventService struct {
}

func NewEventService() *EventService {
	return &EventService{}
}

func (e *EventService) EncodeEvent(eventType string, payload interface{}) []byte {
	event, err := parseEvent(eventType, payload)
	if err != nil {
		log.Fatalf("Someting went wrong parsing event: %v", err)
	}

	b, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Something went wrong sending event to user: %v", err)
	}

	return b
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
