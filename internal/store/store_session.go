package store

import (
	"sync"

	"carlospayan/agent-comercial-ai/internal/catalog"

	"github.com/sashabaranov/go-openai"
)

var (
	mu           sync.Mutex
	messageHist  = make(map[string][]openai.ChatCompletionMessage)
	lastCarStore = make(map[string]catalog.Car)
)

func GetHistory(sessionID string) []openai.ChatCompletionMessage {
	mu.Lock()
	defer mu.Unlock()
	if hist, ok := messageHist[sessionID]; ok {
		return hist
	}
	messageHist[sessionID] = []openai.ChatCompletionMessage{}
	return messageHist[sessionID]
}

func AppendMessage(sessionID string, msg openai.ChatCompletionMessage) {
	mu.Lock()
	defer mu.Unlock()
	messageHist[sessionID] = append(messageHist[sessionID], msg)
}

func SetLastCar(sessionID string, car catalog.Car) {
	mu.Lock()
	defer mu.Unlock()
	lastCarStore[sessionID] = car
}

func GetLastCar(sessionID string) (catalog.Car, bool) {
	mu.Lock()
	defer mu.Unlock()
	car, ok := lastCarStore[sessionID]
	return car, ok
}

func DeleteHistory(sessionID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(messageHist, sessionID)
	delete(lastCarStore, sessionID)
}
