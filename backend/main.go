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
	"github.com/sashabaranov/go-openai"
)

type App struct {
	db       db.DB
	aiClient *openai.Client
	cron     cron.Cron
}

func (app *App) runPipeline(w http.ResponseWriter, r *http.Request) {
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
	stopAtId, err := app.getLatestBill()
	if err != nil {
		log.Fatal("Failed to get latest bill ID from database", err)
	}

	scraperStage := scraper.Scraper{StopAtId: stopAtId}
	// aiStage := AI.NewAgent(app.aiClient)
	jsonSink := JSON.NewJSONSink(w)
	dbSink := DB.NewDBSink(app.db)

	p.RunPipeline(scraperStage, []handlers.Transformer{}, []handlers.Sink{jsonSink, dbSink})
}

func (app *App) scrapeHandler(w http.ResponseWriter, r *http.Request) {
	app.runPipeline(w, r)
}

func (app *App) getLatestBill() (int, error) {
	latestbill, err := app.db.GetLatestBill(context.Background())
	if err != nil {
		log.Fatal("Failed to get latest bill")
		return 0, err
	}

	return latestbill.BillId, nil
}

func (app *App) getLatestBillHandler(w http.ResponseWriter, r *http.Request) {
	billId, err := app.getLatestBill()
	if err != nil {
		http.Error(w, "Failed to get latest bill", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d", billId)
}

// // Cron job: runs pipeline with DB sink only
// func (app *App) cronJob() {
// 	fmt.Println("⏰ Cron job triggered")
// 	app.runPipeline(nil, false)
// }

func main() {
	// load env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// db connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	connPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer connPool.Close()
	if err := connPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	fmt.Println("🚀 Successfully connected to the database pool!")
	p := pg.NewPGWriter(connPool)

	// create openai client
	openAiApiKey := os.Getenv("OPENAI_API_KEY")
	aiClient := openai.NewClient(openAiApiKey)

	// pass in required dependencies
	app := &App{db: p, aiClient: aiClient, cron: cron.Cron{}}

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	// Set up the HTTP server and register the handler.
	// Pass the application instance to the handler.
	http.HandleFunc("/get-latest-bill", app.getLatestBillHandler)
	http.HandleFunc("/scrape", app.scrapeHandler)

	// schedule cron job (e.g., every 1 minute)
	// _, err = app.cron.AddFunc("@every 1m", app.cronJob)
	// if err != nil {
	// 	log.Fatalf("Failed to schedule cron job: %v", err)
	// }
	// app.cron.Start()
	// defer app.cron.Stop()

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
