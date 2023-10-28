package ports

import (
	"net/http"

	"github.com/dev-szymon/translate-chat/server/internal/core/domain"
)

type HttpServerPort interface {
	HandlePostFile(chat *domain.Chat, ts TranslateServicePort) http.HandlerFunc
	HandleWS(chat *domain.Chat, e EventService) http.Handler
}
