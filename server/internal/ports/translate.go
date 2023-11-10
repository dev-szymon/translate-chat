package ports

import (
	"context"

	"github.com/dev-szymon/translate-chat/server/internal/core/chat"
)

type TranslateServicePort interface {
	TranscribeAudio(ctx context.Context, sourceLang string, file []byte) (*chat.Transcript, error)
	TranslateText(ctx context.Context, sourceLang, targetLang, text string) (*chat.Translation, error)
}
