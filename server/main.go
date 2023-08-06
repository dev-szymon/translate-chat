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
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
)

func convertWebmToFlac(file []byte) ([]byte, error) {
	cmd := exec.Command("./ffmpeg", "-i", "pipe:0", "-c:a", "flac", "-f", "flac", "-")

	var (
		output bytes.Buffer
		errors bytes.Buffer
	)
	cmd.Stdin = bytes.NewReader(file)
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

func (e *ApiError) Error() string {
	return e.err
}

func enableCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		handler(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

type TranslatedMessage struct {
	Transcript     string  `json:"transcript"`
	Confidence     float32 `json:"confidence"`
	Translation    string  `json:"translation"`
	TargetLanguage string  `json:"target_language"`
	SourceLanguage string  `json:"source_language"`
}

type GoogleApplicationCredentials struct {
	ProjectID string `json:"project_id"`
}

func translateFlacFile(ctx context.Context, sourceLang, targetLang string, flacFile []byte) (*TranslatedMessage, error) {
	speechService, err := speech.NewClient(ctx)
	if err != nil {
		fmt.Println("Failed to initialise speech service:", err)
		return nil, err
	}
	defer speechService.Close()

	transcriptResponse, err := speechService.Recognize(ctx,
		&speechpb.RecognizeRequest{
			Config: &speechpb.RecognitionConfig{
				Encoding:     speechpb.RecognitionConfig_FLAC,
				LanguageCode: sourceLang,
			},
			Audio: &speechpb.RecognitionAudio{
				AudioSource: &speechpb.RecognitionAudio_Content{
					Content: flacFile,
				},
			},
		},
	)
	if err != nil {
		fmt.Println("Failed to recognize speech:", err)
		return nil, err
	}

	var transcript *speechpb.SpeechRecognitionAlternative
	for _, result := range transcriptResponse.Results {
		for _, a := range result.Alternatives {
			if transcript == nil || transcript.Confidence < a.Confidence {
				transcript = a
			}
		}
	}
	if transcript == nil {
		fmt.Println("Transcription not found.")
		return nil, fmt.Errorf("the service could not generate transcription")
	}

	translateService, err := translate.NewTranslationClient(ctx)
	if err != nil {
		fmt.Println("Failed to initialise translation service", err)
		return nil, err
	}
	defer translateService.Close()

	ctxParent := ctx.Value(contextKey{})
	parent, ok := ctxParent.(string)
	if !ok {
		return nil, fmt.Errorf("project parent missing")

	}

	translationResponse, err := translateService.TranslateText(ctx, &translatepb.TranslateTextRequest{
		Contents:           []string{transcript.Transcript},
		TargetLanguageCode: targetLang,
		SourceLanguageCode: sourceLang,
		Parent:             parent,
	})

	if err != nil {
		fmt.Println("Translation failed:", err)
		return nil, err
	}

	tm := &TranslatedMessage{
		Transcript:     transcript.Transcript,
		Confidence:     transcript.Confidence,
		Translation:    translationResponse.Translations[0].TranslatedText,
		TargetLanguage: targetLang,
		SourceLanguage: sourceLang,
	}
	return tm, nil
}

type contextKey struct{}

func handleTranscribe(w http.ResponseWriter, r *http.Request) {
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
	parent := fmt.Sprintf("projects/%v", googleAppCreds.ProjectID)
	ctx := context.WithValue(context.Background(), contextKey{}, parent)

	if r.Method == "POST" {
		targetLang := r.FormValue("targetLanguage")
		if targetLang == "" {
			fmt.Println("Target language missing")
			writeJSON(w, http.StatusUnprocessableEntity, &ApiError{
				err:    "Please provide target language.",
				status: http.StatusUnprocessableEntity,
			})
		}
		sourceLang := r.FormValue("sourceLanguage")
		if sourceLang == "" {
			fmt.Println("Source language missing")
			writeJSON(w, http.StatusUnprocessableEntity, &ApiError{
				err:    "Please provide source language.",
				status: http.StatusUnprocessableEntity,
			})
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Failed to parse form file:", err)
			writeJSON(w, http.StatusUnprocessableEntity, &ApiError{
				err:    "Failed to parse form file",
				status: http.StatusUnsupportedMediaType,
			})
		}
		defer file.Close()

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

		t, err := translateFlacFile(ctx, sourceLang, targetLang, flacFile)
		if err != nil {
			fmt.Println("failed to translate audio file", err)
			writeJSON(w, http.StatusInternalServerError, &ApiError{
				err:    "failed to translate audio file",
				status: http.StatusInternalServerError,
			})
		}

		writeJSON(w, http.StatusOK, t)
	}
}

type HealthCheck struct {
	Time   time.Time     `json:"time"`
	Uptime time.Duration `json:"uptime"`
}

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./google_application_credentials.json")
	serverStart := time.Now()

	http.HandleFunc("/transcribe", enableCors(handleTranscribe))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		fmt.Println("health check", now)
		writeJSON(w, http.StatusOK, &HealthCheck{
			Time:   now,
			Uptime: time.Since(serverStart),
		})
	})

	fmt.Println("Translation service starting on:", "http://localhost:8055")
	log.Fatal(http.ListenAndServe(":8055", nil))
}
