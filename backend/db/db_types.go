package db

import (
	"context"
	"web-scraper/backend/model"
)

type DB interface {
	InsertBill(ctx context.Context, bill model.Bill) (int, error)
	GetBill(ctx context.Context, id int) (model.Bill, error)
	GetLatestBill(ctx context.Context) (model.Bill, error)
	UpdateBill(ctx context.Context, bill model.Bill) error
	DeleteBill(ctx context.Context, id int) error
}
