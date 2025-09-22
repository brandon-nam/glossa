package DB

import (
	"context"
	"log"
	"web-scraper/backend/db"
	"web-scraper/backend/model"
)

// sink that writes to DB
type DBSink struct {
	db db.DB
}

func NewDBSink(db db.DB) *DBSink {
	return &DBSink{db: db}
}

func (s *DBSink) Consume(in <-chan model.Bill) {
	for bill := range in {
		_, err := s.db.InsertBill(context.Background(), bill)
		if err != nil {
			log.Printf("failed to insert bill: %v", err)
		}
	}
}
