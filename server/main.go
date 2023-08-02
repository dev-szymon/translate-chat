package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
)

func convertWebmToFlac(inputFile []byte) ([]byte, error) {
	cmd := exec.Command("./ffmpeg", "-i", "pipe:0", "-c:a", "flac", "-f", "flac", "-")

	var (
		output bytes.Buffer
		errors bytes.Buffer
	)
	cmd.Stdin = bytes.NewReader(inputFile)
	cmd.Stdout = &output
	cmd.Stderr = &errors

	err := cmd.Run()
	if err != nil {
		log.Printf("FFmpeg stderr: %s", errors.String())
		return nil, err
	}

	return output.Bytes(), nil
}

type ApiError struct {
	err    string
	status int
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)

}

func (e *ApiError) Error() string {
	return e.err
}

type TranscribedMessage struct {
	Transcript  string
	Translation string
}

type GoogleApplicationCredentials struct {
	ProjectID string `json:"project_id"`
}

func transcribeHandler(w http.ResponseWriter, r *http.Request) {
	googleAppCredsFile, err := os.Open(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
		writeJSON(w, http.StatusInternalServerError, &ApiError{
			err:    "Failed to parse form file",
			status: http.StatusInternalServerError,
		})
	}
	c, err := io.ReadAll(googleAppCredsFile)
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
		writeJSON(w, http.StatusInternalServerError, &ApiError{
			err:    "Failed to parse form file",
			status: http.StatusInternalServerError,
		})
	}
	var googleAppCreds GoogleApplicationCredentials
	err = json.Unmarshal(c, &googleAppCreds)
	if err != nil {
		fmt.Println("Failed unmarshal credentials file:", err)
		writeJSON(w, http.StatusInternalServerError, &ApiError{
			err:    "Failed unmarshal credentials file",
			status: http.StatusInternalServerError,
		})
	}

	fmt.Println(googleAppCreds)

	if r.Method == "POST" {
		file, _, err := r.FormFile("file")
		targetLang := "en-US"
		sourceLang := "pl-PL"

		if err != nil {
			fmt.Println("Failed to parse form file:", err)
			writeJSON(w, http.StatusUnprocessableEntity, &ApiError{
				err:    "Failed to parse form file",
				status: http.StatusUnsupportedMediaType,
			})
		}
		defer file.Close()

		ctx := context.Background()

		speechService, err := speech.NewClient(ctx)
		if err != nil {
			fmt.Println("Failed to initialise speech service:", err)
			writeJSON(w, http.StatusInternalServerError, &ApiError{
				err:    "Failed to initialise speech service",
				status: http.StatusInternalServerError,
			})
		}
		defer speechService.Close()

		audioBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Failed to read audio file:", err)
			writeJSON(w, http.StatusUnprocessableEntity, &ApiError{
				err:    "Failed to read audio file",
				status: http.StatusUnprocessableEntity,
			})
		}

		flacFile, err := convertWebmToFlac(audioBytes)
		if err != nil {
			fmt.Println("Failed to convert audio file", err)
			writeJSON(w, http.StatusUnprocessableEntity, &ApiError{
				err:    "Failed to convert audio file",
				status: http.StatusUnprocessableEntity,
			})
		}

		req := &speechpb.RecognizeRequest{
			Config: &speechpb.RecognitionConfig{
				Encoding:     speechpb.RecognitionConfig_FLAC,
				LanguageCode: sourceLang,
			},
			Audio: &speechpb.RecognitionAudio{
				AudioSource: &speechpb.RecognitionAudio_Content{
					Content: flacFile,
				},
			},
		}
		resp, err := speechService.Recognize(ctx, req)
		if err != nil {
			fmt.Println("Failed to recognize speech:", err)
			writeJSON(w, http.StatusInternalServerError, &ApiError{
				err:    "Failed to recognize speech",
				status: http.StatusInternalServerError,
			})
		}

		translateService, err := translate.NewTranslationClient(ctx)
		if err != nil {
			fmt.Println("Failed to initialise translation service", err)
			writeJSON(w, http.StatusInternalServerError, &ApiError{
				err:    "Failed to initialise translation service",
				status: http.StatusInternalServerError,
			})
		}
		defer translateService.Close()

		var transcript *speechpb.SpeechRecognitionAlternative

		for _, result := range resp.Results {
			for _, a := range result.Alternatives {
				if transcript == nil || transcript.Confidence < a.Confidence {
					transcript = a
				}
			}
		}
		if transcript != nil {
			translationResponse, err := translateService.TranslateText(ctx, &translatepb.TranslateTextRequest{
				Contents:           []string{transcript.Transcript},
				TargetLanguageCode: targetLang,
				SourceLanguageCode: sourceLang,
				Parent:             fmt.Sprintf("project/%v", googleAppCreds.ProjectID),
			})
			if err != nil {
				fmt.Println("Translation failed:", err)
				writeJSON(w, http.StatusInternalServerError, &ApiError{
					err:    "Translation failed",
					status: http.StatusInternalServerError,
				})
			}
			fmt.Println(translationResponse)

			writeJSON(w, http.StatusOK, &TranscribedMessage{
				Transcript:  transcript.Transcript,
				Translation: translationResponse.Translations[0].TranslatedText,
			})
		}

	}
}

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./google_application_credentials.json")

	http.HandleFunc("/transcribe", transcribeHandler)

	fmt.Println("Translation service starting on:", "http://localhost:8055")
	log.Fatal(http.ListenAndServe("localhost:8055", nil))
}
