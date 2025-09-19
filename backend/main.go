package main

import (
	"fmt"
	"log"
	"net/http"

	"web-scraper/backend/pipeline"
	"web-scraper/backend/pipeline/handlers"
	"web-scraper/backend/pipeline/handlers/scraper"
	"web-scraper/backend/pipeline/handlers/writer"

	"github.com/joho/godotenv"
)

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	// 2. Pipeline with GPT
	// apiKey := os.Getenv("OPENAI_API_KEY")
	// if apiKey == "" {
	// 	http.Error(w, "Missing OPENAI_API_KEY", http.StatusInternalServerError)
	// 	return
	// }

	p := &pipeline.Pipeline{}
	scraperStage := scraper.Scraper{}
	// aiStage := AI.NewAgent(apiKey)
	jsonSink := &writer.JSONSink{Writer: w}

	p.RunPipeline(scraperStage, []handlers.Transformer{}, jsonSink)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	http.HandleFunc("/scrape", scrapeHandler)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
