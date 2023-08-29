package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"os/exec"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
)

func ConvertBytesToFlac(f []byte) ([]byte, error) {
	b, err := io.ReadAll(bytes.NewReader(f))
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

func ConvertFileToFlac(file multipart.File) ([]byte, error) {
	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error while reading file: ", err)
		return nil, err
	}

	f, err := ConvertBytesToFlac(b)
	if err != nil {
		fmt.Println("Error converting file: ", err)
		return nil, err
	}
	return f, nil
}

type TranslationResponse struct {
	Translation    string `json:"translation"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

type TranscriptResponse struct {
	Transcript string  `json:"transcript"`
	Confidence float32 `json:"confidence"`
}

func TranscribeFlacFile(ctx context.Context, sourceLang string, flacFile []byte) (*TranscriptResponse, error) {
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

	tr := &TranscriptResponse{
		Transcript: transcript.Transcript,
		Confidence: transcript.Confidence,
	}

	return tr, nil
}

func TranslateTranscript(ctx context.Context, sourceLang, targetLang, transcript string) (*TranslationResponse, error) {
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
		Contents:           []string{transcript},
		TargetLanguageCode: targetLang,
		SourceLanguageCode: sourceLang,
		Parent:             fmt.Sprintf("projects/%v", projectId),
	})

	if err != nil {
		fmt.Println("Translation failed:", err)
		return nil, err
	}

	tr := &TranslationResponse{
		Translation:    translationResponse.Translations[0].TranslatedText,
		TargetLanguage: targetLang,
		SourceLanguage: sourceLang,
	}
	return tr, nil
}

var adjectives = []string{
	"silly",
	"witty",
	"bubbly",
	"zany",
	"cheeky",
	"goofy",
	"wacky",
	"whimsical",
	"quirky",
	"hilarious",
	"bizarre",
	"absurd",
	"lively",
	"funky",
	"zesty",
	"jolly",
	"playful",
	"ridiculous",
	"quizzical",
	"eccentric",
}

var nouns = []string{
	"banana",
	"octopus",
	"tornado",
	"pajamas",
	"cupcake",
	"disco-ball",
	"narwhal",
	"pillow",
	"llama",
	"pickle",
	"bubblegum",
	"sock-puppet",
	"chicken",
	"sushi",
	"unicorn",
	"hotdog",
	"gummy-bear",
	"pineapple",
	"hippo",
	"marshmallow",
}

func GenerateRoomName() string {
	a1 := rand.Intn(19)
	a2 := rand.Intn(19)
	n := rand.Intn(19)
	return fmt.Sprintf("%s-%s-%s", adjectives[a1], adjectives[a2], nouns[n])
}
