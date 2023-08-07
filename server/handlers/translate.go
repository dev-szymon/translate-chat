package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dev-szymon/translate-chat/server/lib"
)

func HandleTranslateFile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method == "POST" {
		targetLang := r.FormValue("targetLanguage")
		if targetLang == "" {
			fmt.Println("Target language missing")
			writeJSON(w, http.StatusUnprocessableEntity, &apiError{
				err:    "Please provide target language.",
				status: http.StatusUnprocessableEntity,
			})
		}
		sourceLang := r.FormValue("sourceLanguage")
		if sourceLang == "" {
			fmt.Println("Source language missing")
			writeJSON(w, http.StatusUnprocessableEntity, &apiError{
				err:    "Please provide source language.",
				status: http.StatusUnprocessableEntity,
			})
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Failed to parse form file:", err)
			writeJSON(w, http.StatusUnprocessableEntity, &apiError{
				err:    "Failed to parse form file",
				status: http.StatusUnsupportedMediaType,
			})
		}
		defer file.Close()

		flacFile, err := lib.ConvertFileToFlac(file)
		if err != nil {
			fmt.Println("Failed to convert audio file", err)
			writeJSON(w, http.StatusUnprocessableEntity, &apiError{
				err:    "Failed to convert audio file",
				status: http.StatusUnprocessableEntity,
			})
		}

		t, err := lib.TranslateFlacFile(ctx, sourceLang, targetLang, flacFile)
		if err != nil {
			fmt.Println("failed to translate audio file", err)
			writeJSON(w, http.StatusInternalServerError, &apiError{
				err:    "failed to translate audio file",
				status: http.StatusInternalServerError,
			})
		}

		writeJSON(w, http.StatusOK, t)
	}
}
