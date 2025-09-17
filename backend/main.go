package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"web-scraper/backend/pipeline"
	"web-scraper/backend/scraper"

	"github.com/joho/godotenv"
)

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	startURL := fmt.Sprintf("https://opinion.lawmaking.go.kr/gcom/nsmLmSts/out?pageIndex=%s", page)

	in := make(chan scraper.Bill)
	out := make(chan scraper.Bill)

	// 1. Scrape
	go scraper.ScrapeBills(startURL, in)

	// 2. Pipeline with GPT
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		http.Error(w, "Missing OPENAI_API_KEY", http.StatusInternalServerError)
		return
	}
	pipe := pipeline.NewPipeline(apiKey)
	go pipe.RunPipeline(in, out)

	// 3. Return as JSON
	// Stream results in real-time
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	w.Write([]byte("[")) // open JSON array
	first := true

	for bill := range out {
		if !first {
			w.Write([]byte(",")) // separate JSON objects
		}
		first = false
		if err := enc.Encode(bill); err != nil {
			log.Println("encode error:", err)
			break
		}
		w.(http.Flusher).Flush() // 🔥 send chunk immediately to client
	}

	w.Write([]byte("]")) // close JSON array
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
