package translate

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/dev-szymon/translate-chat/server/internal/core/chat"
)

type TranslateService struct{}

func NewTranslateService() *TranslateService {
	return &TranslateService{}
}

func (ts *TranslateService) TranscribeAudio(ctx context.Context, sourceLang string, file []byte) (*chat.Transcript, error) {
	speechService, err := speech.NewClient(ctx)
	if err != nil {
		fmt.Println("Failed to initialise speech service:", err)
		return nil, err
	}
	defer speechService.Close()

	flacFile, err := convertBytesToFlac(file)
	if err != nil {
		fmt.Println("Error converting file to flac", err)
		return nil, err
	}

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

	var bestTranscript *speechpb.SpeechRecognitionAlternative
	for _, result := range transcriptResponse.Results {
		for _, a := range result.Alternatives {
			if bestTranscript == nil || bestTranscript.Confidence < a.Confidence {
				bestTranscript = a
			}
		}
	}
	if bestTranscript == nil {
		fmt.Println("Transcription not found.")
		return nil, fmt.Errorf("the service could not generate transcription")
	}

	transcript := &chat.Transcript{
		Text:       bestTranscript.Transcript,
		Confidence: bestTranscript.Confidence,
	}
	return transcript, nil
}

func (ts *TranslateService) TranslateText(ctx context.Context, sourceLang, targetLang, text string) (*chat.Translation, error) {
	translateService, err := translate.NewTranslationClient(ctx)
	if err != nil {
		fmt.Println("Failed to initialise translation service", err)
		return nil, err
	}
	defer translateService.Close()

	projectId := os.Getenv("GOOGLE_PROJECT_ID")
	if projectId == "" {
		return nil, fmt.Errorf("google project id environmental variable missing")
	}

	translationResponse, err := translateService.TranslateText(ctx, &translatepb.TranslateTextRequest{
		Contents:           []string{text},
		TargetLanguageCode: targetLang,
		SourceLanguageCode: sourceLang,
		Parent:             fmt.Sprintf("projects/%v", projectId),
	})

	if err != nil {
		fmt.Println("Translation failed:", err)
		return nil, err
	}

	translation := &chat.Translation{
		Text:       translationResponse.Translations[0].TranslatedText,
		TargetLang: targetLang,
		SourceLang: sourceLang,
	}
	return translation, nil
}

func convertBytesToFlac(file []byte) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-c:a", "flac", "-f", "flac", "-")

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
