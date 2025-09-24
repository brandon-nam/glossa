package pg

import (
	"context"
	"fmt"
	"web-scraper/backend/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGWriter struct {
	conn *pgxpool.Pool
}

func NewPGWriter(pool *pgxpool.Pool) *PGWriter {
	return &PGWriter{
		conn: pool,
	}
}

// InsertBill inserts a new bill record into the database and returns the generated ID.
func (d *PGWriter) InsertBill(ctx context.Context, bill model.Bill) (int, error) {
	var id int
	err := d.conn.QueryRow(
		ctx,
		`INSERT INTO bills (name, proposers, main_text, summary, categories, detail_url, bill_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		bill.Name,
		bill.Proposers,
		bill.MainText,
		bill.Summary,
		bill.Categories,
		bill.DetailURL,
		bill.BillId,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert bill: %w", err)
	}
	return id, nil
}

// GetBill retrieves a single bill record by its ID.
func (d *PGWriter) GetBill(ctx context.Context, billId int) (model.Bill, error) {
	var bill model.Bill
	err := d.conn.QueryRow(
		ctx,
		`SELECT id, name, proposers, department, parliamentary_status, resolution_status, main_text, summary, categories FROM assembly_bill WHERE bill_id = $1`,
		billId,
	).Scan(
		&bill.Id,
		&bill.Name,
		&bill.Proposers,
		&bill.Department,
		&bill.ParliamentaryStatus,
		&bill.ResolutionStatus,
		&bill.MainText,
		&bill.Summary,
		&bill.Categories,
	)
	if err != nil {
		return model.Bill{}, fmt.Errorf("failed to get bill with id %d: %w", billId, err)
	}
	return bill, nil
}

func (d *PGWriter) GetLatestBill(ctx context.Context) (model.Bill, error) {
	var bill model.Bill
	err := d.conn.QueryRow(ctx,
		`SELECT bill_id
		FROM bills
		ORDER BY bill_id DESC
		LIMIT 1`,
	).Scan(&bill.BillId)
	return bill, err
}

// UpdateBill updates an existing bill record in the database.
func (d *PGWriter) UpdateBill(ctx context.Context, bill model.Bill) error {
	_, err := d.conn.Exec(
		ctx,
		`UPDATE bills SET name = $1, proposers = $2, main_text = $3, summary = $4, categories = $5 WHERE id = $6`,
		bill.Name,
		bill.Proposers,
		bill.MainText,
		bill.Summary,
		bill.Categories,
		bill.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update bill with id %d: %w", bill.Id, err)
	}
	return nil
}

// DeleteBill deletes a bill record from the database by its ID.
func (d *PGWriter) DeleteBill(ctx context.Context, id int) error {
	_, err := d.conn.Exec(ctx, `DELETE FROM bills WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete bill with id %d: %w", id, err)
	}
	return nil
}
