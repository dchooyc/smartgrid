package main

import (
	"context"
	"database/sql"
	"os"
	"strconv"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// defaultDatabaseURL is used when DATABASE_URL is not set.
const defaultDatabaseURL = "postgres://localhost:5432/smartgrid?sslmode=disable"

var db *sql.DB

// openDatabase prepares the connection pool. It does not connect yet, so the
// server still starts when PostgreSQL is down.
func openDatabase() (*sql.DB, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = defaultDatabaseURL
	}
	return sql.Open("postgres", url)
}

// querier is satisfied by both the connection pool and a transaction.
type querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// loadExistingAssets returns the fields of every stored asset that
// validation needs. Their line is 0 because they are not in the upload.
func loadExistingAssets(ctx context.Context, q querier) ([]asset, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT asset_id, COALESCE(parent_asset_id, ''), asset_type, asset_name, operational_status
		FROM assets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []asset
	for rows.Next() {
		var a asset
		if err := rows.Scan(&a.id, &a.parentID, &a.assetType, &a.name, &a.status); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

// insertAssets inserts the assets in the order given. A plain INSERT fails on
// an existing asset_id, so stored assets are never overwritten.
func insertAssets(ctx context.Context, tx *sql.Tx, assets []asset) error {
	insert, err := tx.PrepareContext(ctx, `
		INSERT INTO assets (
			asset_id, parent_asset_id, asset_type, asset_name, operational_status,
			voltage_kv, rating_kva, manufacturer, model, serial_number, commissioned_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`)
	if err != nil {
		return err
	}
	defer insert.Close()

	for _, a := range assets {
		_, err := insert.ExecContext(ctx,
			a.id, nullable(a.parentID), a.assetType, a.name, a.status,
			nullable(formatNumber(a.voltageKV)), nullable(formatNumber(a.ratingKVA)),
			nullable(a.manufacturer), nullable(a.model), nullable(a.serialNumber),
			nullable(formatDate(a.commissionedDate)))
		if err != nil {
			return err
		}
	}
	return nil
}

// nullable stores an empty CSV value as NULL.
func nullable(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}

// formatNumber rewrites a validated number in plain decimal form, which
// PostgreSQL always accepts. Empty stays empty.
func formatNumber(value string) string {
	number, ok := parseNumber(value)
	if !ok {
		return ""
	}
	return strconv.FormatFloat(number, 'f', -1, 64)
}

// formatDate rewrites a validated day/month/year date as year-month-day.
// Empty stays empty.
func formatDate(value string) string {
	date, ok := parseDate(value)
	if !ok {
		return ""
	}
	return date.Format("2006-01-02")
}
