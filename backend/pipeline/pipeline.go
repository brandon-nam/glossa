package pipeline

import (
	"log"
	"web-scraper/backend/model"
	"web-scraper/backend/pipeline/handlers"
)

// Pipeline with channels
type Pipeline struct{}

// RunPipeline starts processing: scrapeChan -> GPT -> resultChan
func (p *Pipeline) RunPipeline(source handlers.Source, transformers []handlers.Transformer, sink []handlers.Sink) {
	first := make(chan model.Bill)
	prev := first
	log.Println("[Pipeline] Starting pipeline")
	// 1. Start the source
	go source.Run(first)

	// 2. Chain transformers
	for _, t := range transformers {
		next := make(chan model.Bill)
		go t.Transform(prev, next)
		prev = next
	}

	// 3. Send final channel to sink

	for _, sink := range sink {
		go sink.Consume(prev)
	}

	log.Println("[Pipeline] Exit pipeline")
}
