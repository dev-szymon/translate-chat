package ports

type EventService interface {
	EncodeEvent(eventType string, payload interface{}) []byte
}
