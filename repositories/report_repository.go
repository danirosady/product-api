package repositories

import (
	"database/sql"
	"product-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetReport(startInclusive time.Time, endExclusive time.Time) (*models.ReportResponse, error) {
	totalsQuery := `
		SELECT
			COALESCE(COUNT(*), 0) AS total_transactions,
			COALESCE(SUM(t.total_amount), 0) AS total_revenue,
			COALESCE((
				SELECT SUM(td.quantity)
				FROM transaction_details td
				JOIN transactions t2 ON t2.id = td.transaction_id
				WHERE t2.created_at >= $1 AND t2.created_at < $2
			), 0) AS total_items_sold
		FROM transactions t
		WHERE t.created_at >= $1 AND t.created_at < $2
	`

	var totalTransactions int64
	var totalRevenue int64
	var totalItemsSold int64
	if err := repo.db.QueryRow(totalsQuery, startInclusive, endExclusive).Scan(&totalTransactions, &totalRevenue, &totalItemsSold); err != nil {
		return nil, err
	}

	bestSellerQuery := `
		SELECT
			p.name,
			COALESCE(SUM(td.quantity), 0) AS qty_terjual
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.name
		ORDER BY qty_terjual DESC, p.name ASC
		LIMIT 1
	`

	var bestName string
	var bestQty int64
	err := repo.db.QueryRow(bestSellerQuery, startInclusive, endExclusive).Scan(&bestName, &bestQty)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		bestName = ""
		bestQty = 0
	}

	return &models.ReportResponse{
		StartTime:      startInclusive,
		EndTime:        endExclusive,
		TotalRevenue:   int(totalRevenue),
		TotalTransaksi: int(totalTransactions),
		ProdukTerlaris: models.ProdukTerlaris{
			Nama:       bestName,
			QtyTerjual: int(bestQty),
		},
	}, nil
}
