package domain

type Message struct {
	Transcription *Transcript
	Translations  map[string]*Translation
	SenderId      string
}
