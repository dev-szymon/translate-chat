package translate

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dev-szymon/translate-chat/server/internal/core/chat"
)

const MOCK_PREFIX = "[mock] "

type DebugTranslateService struct{}

func NewDebugTranslateService() *DebugTranslateService {
	return &DebugTranslateService{}
}

func (ts *DebugTranslateService) TranscribeAudio(ctx context.Context, sourceLang string, file []byte) (*chat.Transcript, error) {
	sourceLangMessages, ok := mockMessages[sourceLang]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", sourceLang)
	}
	randomIndex := rand.Intn(len(sourceLangMessages) - 1)
	randomConfidence := rand.Float32()

	transcript := &chat.Transcript{
		Text:       fmt.Sprintf("%s%s", MOCK_PREFIX, sourceLangMessages[randomIndex]),
		Confidence: randomConfidence,
		SourceLang: sourceLang,
	}

	time.Sleep(50 * time.Millisecond)

	return transcript, nil
}

func (ts *DebugTranslateService) TranslateText(ctx context.Context, sourceLang, targetLang, text string) (*chat.Translation, error) {
	sourceLangMessages, ok := mockMessages[sourceLang]
	if !ok {
		return nil, fmt.Errorf("source language unrecognized: %s", sourceLang)
	}

	trimmedText := strings.TrimPrefix(text, MOCK_PREFIX)
	var translationMessageIndex int
	for i, m := range sourceLangMessages {
		if trimmedText == m {
			translationMessageIndex = i
		}
	}

	targetLangMessages, ok := mockMessages[targetLang]
	if !ok {
		return nil, fmt.Errorf("target language unsupported: %s", targetLang)
	}
	if translationMessageIndex > len(targetLangMessages)-1 {
		return nil, fmt.Errorf("target language message index out of bounds: %d", translationMessageIndex)
	}

	translation := &chat.Translation{
		Text:       fmt.Sprintf("%s%s", MOCK_PREFIX, targetLangMessages[translationMessageIndex]),
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}
	time.Sleep(100 * time.Millisecond)
	return translation, nil
}

var mockMessages = map[string][]string{
	"pl-PL": {
		"Cześć, jak ci minął dzień?",
		"Ostatnio oglądałeś jakieś dobre filmy?",
		"Jaka jest twoja ulubiona kuchnia?",
		"Masz jakieś ekscytujące plany na weekend?",
		"Właśnie skończyłem czytać fantastyczną książkę, chcesz polecić coś?",
		"Mój dzień był jak dotąd dobry, dzięki za pytanie!",
		"Tak, niedawno obejrzałem 'Diunę' i była niesamowita!",
		"Uwielbiam kuchnię włoską, szczególnie makaron i pizzę.",
		"Planuję pójść na wędrówkę z przyjaciółmi w ten weekend.",
		"Oczywiście! Chętnie przyjmę rekomendację książki. Co przeczytałeś?",
	},
	"en-US": {
		"Hey, how's your day going?",
		"Have you watched any good movies lately?",
		"What's your favorite type of cuisine?",
		"Do you have any exciting plans for the weekend?",
		"I just finished reading a fantastic book, want a recommendation?",
		"My day's been good so far, thanks for asking!",
		"Yes, I recently watched 'Dune' and it was amazing!",
		"I love Italian cuisine, especially pasta and pizza.",
		"I'm planning to go hiking with some friends this weekend.",
		"Absolutely! I'd love a book recommendation. What did you read?",
	},
	"es-ES": {
		"Hola, ¿cómo va tu día?",
		"¿Has visto alguna película buena recientemente?",
		"¿Cuál es tu tipo de cocina favorita?",
		"¿Tienes planes emocionantes para el fin de semana?",
		"Acabo de terminar de leer un libro fantástico, ¿quieres una recomendación?",
		"Mi día ha ido bien hasta ahora, ¡gracias por preguntar!",
		"Sí, recientemente vi 'Dune' y fue asombrosa.",
		"Me encanta la cocina italiana, especialmente la pasta y la pizza.",
		"Estoy planeando hacer senderismo con algunos amigos este fin de semana.",
		"¡Por supuesto! Me encantaría una recomendación de libro. ¿Qué has leído?",
	},
	"it-IT": {
		"Ciao, come sta andando la tua giornata?",
		"Hai visto di recente dei bei film?",
		"Qual è il tuo tipo di cucina preferito?",
		"Hai dei piani eccitanti per il fine settimana?",
		"Ho appena finito di leggere un libro fantastico, vuoi una raccomandazione?",
		"La mia giornata è stata finora buona, grazie per aver chiesto!",
		"Sì, di recente ho visto 'Dune' ed è stato incredibile!",
		"Amo la cucina italiana, soprattutto la pasta e la pizza.",
		"Sto pianificando di fare escursionismo con degli amici questo fine settimana.",
		"Assolutamente! Mi piacerebbe una raccomandazione di un libro. Cosa hai letto?",
	},
	"de-DE": {
		"Hallo, wie läuft dein Tag?",
		"Hast du in letzter Zeit gute Filme gesehen?",
		"Was ist deine Lieblingsküche?",
		"Hast du aufregende Pläne für das Wochenende?",
		"Ich habe gerade ein fantastisches Buch gelesen, möchtest du eine Empfehlung?",
		"Mein Tag war bisher gut, danke der Nachfrage!",
		"Ja, ich habe kürzlich 'Dune' gesehen, und er war großartig!",
		"Ich liebe die italienische Küche, besonders Pasta und Pizza.",
		"Ich plane, am Wochenende mit einigen Freunden wandern zu gehen.",
		"Klar! Ich würde gerne eine Buchempfehlung geben. Was hast du gelesen?",
	},
}
