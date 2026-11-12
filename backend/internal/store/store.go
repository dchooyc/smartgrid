// Package store reads and writes assets in PostgreSQL.
package store

import (
	"context"
	"database/sql"
	"strconv"

	_ "github.com/lib/pq" // PostgreSQL driver

	"smartgrid/backend/internal/asset"
)

// Store is the assets database.
type Store struct {
	db *sql.DB
}

// Open prepares the connection pool. It does not connect yet, so the server
// still starts when PostgreSQL is down.
func Open(url string) (*Store, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Assets returns the fields of every stored asset that validation needs.
// Their line is 0 because they are not in an upload.
func (s *Store) Assets(ctx context.Context) ([]asset.Asset, error) {
	return loadAssets(ctx, s.db)
}

// Tx is a database transaction. Commit saves its changes and Rollback
// discards them.
type Tx struct {
	*sql.Tx
}

// Begin starts a transaction.
func (s *Store) Begin(ctx context.Context) (*Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

// Assets is Store.Assets, read inside the transaction.
func (tx *Tx) Assets(ctx context.Context) ([]asset.Asset, error) {
	return loadAssets(ctx, tx.Tx)
}

// InsertAssets inserts the assets in the order given. A plain INSERT fails on
// an existing asset_id, so stored assets are never overwritten.
func (tx *Tx) InsertAssets(ctx context.Context, assets []asset.Asset) error {
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
			a.ID, nullable(a.ParentID), a.Type, a.Name, a.Status,
			nullable(formatNumber(a.VoltageKV)), nullable(formatNumber(a.RatingKVA)),
			nullable(a.Manufacturer), nullable(a.Model), nullable(a.SerialNumber),
			nullable(formatDate(a.CommissionedDate)))
		if err != nil {
			return err
		}
	}
	return nil
}

// querier is satisfied by both the connection pool and a transaction.
type querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func loadAssets(ctx context.Context, q querier) ([]asset.Asset, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT asset_id, COALESCE(parent_asset_id, ''), asset_type, asset_name, operational_status
		FROM assets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []asset.Asset
	for rows.Next() {
		var a asset.Asset
		if err := rows.Scan(&a.ID, &a.ParentID, &a.Type, &a.Name, &a.Status); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

// nullable stores an empty CSV value as NULL.
func nullable(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}

// formatNumber rewrites a validated number in plain decimal form, which
// PostgreSQL always accepts. Empty stays empty.
func formatNumber(value string) string {
	number, ok := asset.ParseNumber(value)
	if !ok {
		return ""
	}
	return strconv.FormatFloat(number, 'f', -1, 64)
}

// formatDate rewrites a validated day/month/year date as year-month-day.
// Empty stays empty.
func formatDate(value string) string {
	date, ok := asset.ParseDate(value)
	if !ok {
		return ""
	}
	return date.Format("2006-01-02")
}
