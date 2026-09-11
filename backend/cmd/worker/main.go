package main

import (
	"log"

	"agora-backend/internal/config"
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

	w := worker.New(temporalClient, cfg.TemporalTaskQueue, worker.Options{})
	agoraworkflow.Register(w)
	log.Printf("[Worker] polling task queue %q", cfg.TemporalTaskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("[Worker] stopped: %v", err)
	}
}
