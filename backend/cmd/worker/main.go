package main

import (
	"log"

	"agora-backend/internal/config"
	"agora-backend/internal/db"
	"agora-backend/internal/llm"
	agoraworkflow "agora-backend/internal/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg := config.LoadConfig()
	temporalClient, err := client.Dial(client.Options{HostPort: cfg.TemporalHost})
	if err != nil {
		log.Fatalf("[Worker] failed to connect to Temporal: %v", err)
	}
	defer temporalClient.Close()
	database, err := db.InitDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("[Worker] failed to connect to PostgreSQL: %v", err)
	}
	defer database.Close()

	w := worker.New(temporalClient, cfg.TemporalTaskQueue, worker.Options{})
	llmClient := llm.NewOpenAICompatibleClient(cfg.LLMBaseURL, cfg.LLMModel, cfg.LLMAPIKey)
	agoraworkflow.Register(w, &agoraworkflow.Activities{DB: database, LLM: llmClient})
	log.Printf("[Worker] polling task queue %q", cfg.TemporalTaskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("[Worker] stopped: %v", err)
	}
}
