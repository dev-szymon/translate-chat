package chat

type Transcript struct {
	Text       string
	Confidence float32
	SourceLang string
}

type Translation struct {
	Text       string
	SourceLang string
	TargetLang string
}
