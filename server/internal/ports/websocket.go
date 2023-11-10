package ports

import (
	"net/http"
)

type WebsocketServerPort interface {
	HandleWS(ts TranslateServicePort) http.HandlerFunc
}
