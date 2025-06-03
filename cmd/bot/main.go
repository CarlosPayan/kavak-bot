package main

import (
	"log"
	"net/http"

	"carlospayan/agent-comercial-ai/internal/catalog"
	"carlospayan/agent-comercial-ai/internal/config"
	"carlospayan/agent-comercial-ai/internal/handlers"
	"carlospayan/agent-comercial-ai/internal/utils"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "carlospayan/agent-comercial-ai/internal/metrics"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	kavakURL := "https://www.kavak.com/mx/blog/sedes-de-kavak-en-mexico"
	content, err := utils.FetchKavakInfo(kavakURL)
	if err != nil {
		log.Fatalf("couldn't extract info from Kavak: %v", err)
	}

	cat, err := catalog.NewCatalog(cfg.OpenAI.APIKey, cfg.Catalog.Path)
	if err != nil {
		log.Fatalf("error cargando catálogo: %v", err)
	}

	r := chi.NewRouter()

	r.Get("/qa", handlers.RAGHandler(cfg, content, cat))

	r.Post("/whatsapp", handlers.WhatsAppHandler(cfg, content, cat))

	r.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s…", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
