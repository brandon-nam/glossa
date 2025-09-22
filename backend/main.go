package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"web-scraper/backend/db"
	"web-scraper/backend/db/pg"
	"web-scraper/backend/pipeline"
	"web-scraper/backend/pipeline/handlers"
	"web-scraper/backend/pipeline/handlers/scraper"
	"web-scraper/backend/pipeline/handlers/sink/DB"
	"web-scraper/backend/pipeline/handlers/sink/JSON"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

type App struct {
	db   db.DB
	cron cron.Cron
}

func (app *App) scrapeHandler(w http.ResponseWriter, r *http.Request) {
	//add context?
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	// 2. Pipeline with GPT
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		http.Error(w, "Missing OPENAI_API_KEY", http.StatusInternalServerError)
		return
	}

	p := &pipeline.Pipeline{}
	scraperStage := scraper.Scraper{}
	// aiStage := AI.NewAgent(apiKey)
	jsonSink := JSON.NewJSONSink(w)
	dbSink := DB.NewDBSink(app.db)

	p.RunPipeline(scraperStage, []handlers.Transformer{}, []handlers.Sink{jsonSink, dbSink})
}

func main() {
	// load env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// get db url
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// create db pool
	connPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// ensure graceful shutdown
	defer connPool.Close()

	// test connection
	if err := connPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	// print on success
	fmt.Println("🚀 Successfully connected to the database pool!")

	// create pg writer
	p := pg.NewPGWriter(connPool)

	// pass in required dependencies
	app := &App{db: p, cron: cron.Cron{}}

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	// Set up the HTTP server and register the handler.
	// Pass the application instance to the handler.
	http.HandleFunc("/scrape", app.scrapeHandler)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
