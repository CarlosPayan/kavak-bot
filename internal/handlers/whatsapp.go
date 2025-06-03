// internal/handlers/whatsapp.go
package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"

	"carlospayan/agent-comercial-ai/internal/catalog"
	"carlospayan/agent-comercial-ai/internal/config"
	"carlospayan/agent-comercial-ai/internal/llm"
	"carlospayan/agent-comercial-ai/internal/metrics"
	"carlospayan/agent-comercial-ai/internal/store"
)

func WhatsAppHandler(cfg *config.Config, kavakInfo string, cat *catalog.Catalog) http.HandlerFunc {

	client := llm.NewClient(cfg.OpenAI.APIKey)

	return func(w http.ResponseWriter, r *http.Request) {
		whatsappStart := time.Now()
		if err := r.ParseForm(); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}

		userBody := r.FormValue("Body")
		from := r.FormValue("From")
		if strings.TrimSpace(userBody) == "" || strings.TrimSpace(from) == "" {
			http.Error(w, "parameters 'Body' or 'From' missing", http.StatusBadRequest)
			return
		}

		sid := from

		history := store.GetHistory(sid)
		if len(history) == 0 {
			store.AppendMessage(sid, openai.ChatCompletionMessage{
				Role:    "system",
				Content: SYSTEM_INSTRUCCTIONS,
			})
			store.AppendMessage(sid, openai.ChatCompletionMessage{
				Role: "assistant",
				Content: fmt.Sprintf(
					"Información de Kavak (propuesta de valor y sucursales):\n%s",
					kavakInfo,
				),
			})
		}

		catStart := time.Now()
		autos, err := cat.Search(r.Context(), userBody, 3)
		catLatency := time.Since(catStart)
		metrics.CatLatency.Observe(float64(catLatency.Milliseconds()))
		if err != nil {
			http.Error(w, fmt.Sprintf("error searching in catalog: %v", err), http.StatusInternalServerError)
			return
		}

		var recs []string
		for i, a := range autos {
			recs = append(recs, fmt.Sprintf(
				"%d) %s %s %s (%d) – Precio: %.2f MXN, Kilometraje: %d km",
				i+1, a.Make, a.Model, a.Version, a.Year, a.Price, a.KM,
			))
		}
		bloqueRecs := "Nuevas recomendaciones (top-3) basadas en tu mensaje:\n" +
			strings.Join(recs, "\n")

		store.AppendMessage(sid, openai.ChatCompletionMessage{
			Role:    "assistant",
			Content: bloqueRecs,
		})

		store.AppendMessage(sid, openai.ChatCompletionMessage{
			Role:    "user",
			Content: userBody,
		})

		llmStart := time.Now()
		history = store.GetHistory(sid)
		answer, err := client.Chat(r.Context(), history)
		llmLatency := time.Since(llmStart)
		metrics.LLMLatency.Observe(float64(llmLatency.Milliseconds()))
		if err != nil {
			http.Error(w, fmt.Sprintf("error calling LLM: %v", err), http.StatusInternalServerError)
			return
		}

		store.AppendMessage(sid, openai.ChatCompletionMessage{
			Role:    "assistant",
			Content: answer,
		})

		w.Header().Set("Content-Type", "application/xml")
		responseXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <Message>%s</Message>
</Response>`, escapeForXML(answer))
		w.Write([]byte(responseXML))

		whatsappLatency := time.Since(whatsappStart)
		metrics.WhatsappHandlerLatency.Observe(float64(whatsappLatency.Milliseconds()))
	}

}

func escapeForXML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}
