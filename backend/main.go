package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"web-scraper/backend/db"
	"web-scraper/backend/db/pg"
	"web-scraper/backend/pipeline"
	"web-scraper/backend/pipeline/handlers"
	"web-scraper/backend/pipeline/handlers/scraper"
	"web-scraper/backend/pipeline/handlers/sink/DB"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sashabaranov/go-openai"
)

type App struct {
	db       db.DB
	aiClient *openai.Client
}

func (app *App) runPipeline() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing OPENAI_API_KEY")
		return
	}

	p := &pipeline.Pipeline{}
	stopAtId, err := app.getLatestBill()
	if err != nil {
		log.Fatal("Failed to get latest bill ID from database", err)
	}

	scraperStage := scraper.Scraper{StopAtId: stopAtId}
	// aiStage := AI.NewAgent(app.aiClient)
	dbSink := DB.NewDBSink(app.db)

	p.RunPipeline(scraperStage, []handlers.Transformer{}, []handlers.Sink{dbSink})
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

func (app *App) getBillsHandler(w http.ResponseWriter, r *http.Request) {
	bills, err := app.db.GetBills(context.Background())
	if err != nil {
		http.Error(w, "Failed to get bills", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bills); err != nil {
		http.Error(w, "Failed to encode bills to JSON", http.StatusInternalServerError)
	}
}

func (app *App) getBillHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, "Bad request: Bill ID not provided", http.StatusBadRequest)
		return
	}
	billIDStr := pathParts[2]

	billID, err := strconv.Atoi(billIDStr)
	if err != nil {
		http.Error(w, "Invalid bill ID format", http.StatusBadRequest)
		return
	}

	bill, err := app.db.GetBill(context.Background(), billID)
	if err != nil {
		fmt.Printf("Error getting bill from database: %v\n", err)
		http.Error(w, "Failed to retrieve bill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(bill); err != nil {
		fmt.Printf("Error encoding bill to JSON: %v\n", err)
		http.Error(w, "Failed to encode bill to JSON", http.StatusInternalServerError)
	}
}

// Cron job: runs pipeline with DB sink only
func (app *App) cronJob() {
	fmt.Println("⏰ Cron job triggered")
	app.runPipeline()
}

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
	app := &App{db: p, aiClient: aiClient}

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	c := cron.New()
	_, _ = c.AddFunc("@every 6h", app.cronJob)
	c.Start()
	defer c.Stop()

	// Set up the HTTP server and register the handler.
	// Pass the application instance to the handler.
	http.HandleFunc("/get-latest-bill", app.getLatestBillHandler)
	http.HandleFunc("/view", app.getBillsHandler)
	http.HandleFunc("/bills/", app.getBillHandler)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
