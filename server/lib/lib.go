package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
)

func ConvertFileToFlac(file multipart.File) ([]byte, error) {
	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error while reading file: ", err)
		return nil, err
	}
	cmd := exec.Command("./ffmpeg", "-i", "pipe:0", "-c:a", "flac", "-f", "flac", "-")

	var (
		output bytes.Buffer
		errors bytes.Buffer
	)
	cmd.Stdin = bytes.NewReader(b)
	cmd.Stdout = &output
	cmd.Stderr = &errors

	err = cmd.Run()
	if err != nil {
		log.Printf("FFmpeg stderr: %s", errors.String())
		return nil, err
	}

	return output.Bytes(), nil
}

type TranslatedMessage struct {
	Transcript     string  `json:"transcript"`
	Confidence     float32 `json:"confidence"`
	Translation    string  `json:"translation"`
	TargetLanguage string  `json:"target_language"`
	SourceLanguage string  `json:"source_language"`
}

func TranslateFlacFile(ctx context.Context, sourceLang, targetLang string, flacFile []byte) (*TranslatedMessage, error) {
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

	projectId := os.Getenv("GOOGLE_PROJECT_ID")
	if projectId == "" {
		return nil, fmt.Errorf("google project id environmental variable missing")
	}

	translationResponse, err := translateService.TranslateText(ctx, &translatepb.TranslateTextRequest{
		Contents:           []string{transcript.Transcript},
		TargetLanguageCode: targetLang,
		SourceLanguageCode: sourceLang,
		Parent:             fmt.Sprintf("projects/%v", projectId),
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
