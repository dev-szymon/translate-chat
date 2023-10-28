package ports

import (
	"context"

	"github.com/dev-szymon/translate-chat/server/internal/core/domain"
)

type TranslateServicePort interface {
	TranscribeAudio(ctx context.Context, sourceLang string, flacFile []byte) (*domain.Transcript, error)
	TranslateText(ctx context.Context, sourceLang, targetLang, text string) (*domain.Translation, error)
}
