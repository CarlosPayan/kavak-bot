package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CatLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "catalog_search_latency_ms",
		Help:    "Time (ms) for cat.Search",
		Buckets: prometheus.ExponentialBuckets(10, 2, 8),
	})

	LLMLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "openai_chat_latency_ms",
		Help:    "TIme (ms) for OpenAI ChatCompletion",
		Buckets: prometheus.ExponentialBuckets(100, 2, 8),
	})

	QAHandlerLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "qa_latency_ms",
		Help:    "Tome for (ms)  handler /qa",
		Buckets: prometheus.ExponentialBuckets(100, 2, 8),
	})

	WhatsappHandlerLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "whatsapp_request_latency_ms",
		Help:    "Time for (ms)  handler /whatsapp",
		Buckets: prometheus.ExponentialBuckets(100, 2, 8),
	})
)

func init() {
	prometheus.MustRegister(CatLatency, LLMLatency, QAHandlerLatency, WhatsappHandlerLatency)
}
