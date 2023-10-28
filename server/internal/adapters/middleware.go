package adapters

import (
	"net/http"
	"os"
)

func EnableCors(handler http.HandlerFunc) http.HandlerFunc {
	webUrl := os.Getenv("WEB_URL")

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", webUrl)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		handler(w, r)
	}
}
