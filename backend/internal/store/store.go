// Package store reads and writes assets in PostgreSQL.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

// ListAssets returns stored assets, ordered by ID. A non-empty assetType
// keeps only assets of that type, a non-empty after keeps only assets whose
// ID comes after it, and a limit above zero caps how many are returned.
func (s *Store) ListAssets(ctx context.Context, assetType, after string, limit int) ([]asset.Stored, error) {
	var conditions []string
	var args []any
	if assetType != "" {
		args = append(args, assetType)
		conditions = append(conditions, fmt.Sprintf("asset_type = $%d", len(args)))
	}
	if after != "" {
		args = append(args, after)
		conditions = append(conditions, fmt.Sprintf("asset_id > $%d", len(args)))
	}

	condition := ""
	if len(conditions) > 0 {
		condition = "WHERE " + strings.Join(conditions, " AND ")
	}
	return s.queryAssets(ctx, condition, limit, args...)
}

// CountAssets counts the stored assets, or only those of a non-empty
// assetType.
func (s *Store) CountAssets(ctx context.Context, assetType string) (asset.Counts, error) {
	counts := asset.Counts{ByType: map[string]int{}, ByStatus: map[string]int{}}

	rows, err := s.db.QueryContext(ctx, `
		SELECT asset_type, operational_status, count(*)
		FROM assets
		WHERE $1 = '' OR asset_type = $1
		GROUP BY asset_type, operational_status`, assetType)
	if err != nil {
		return counts, err
	}
	defer rows.Close()

	for rows.Next() {
		var foundType, status string
		var count int
		if err := rows.Scan(&foundType, &status, &count); err != nil {
			return counts, err
		}
		counts.Total += count
		counts.ByType[foundType] += count
		counts.ByStatus[status] += count
	}
	return counts, rows.Err()
}

// GetAsset returns the asset with the given ID, or nil if there is none.
func (s *Store) GetAsset(ctx context.Context, id string) (*asset.Stored, error) {
	assets, err := s.queryAssets(ctx, "WHERE asset_id = $1", 0, id)
	if err != nil || len(assets) == 0 {
		return nil, err
	}
	return &assets[0], nil
}

// ListChildren returns the assets whose direct parent has the given ID.
func (s *Store) ListChildren(ctx context.Context, parentID string) ([]asset.Stored, error) {
	return s.queryAssets(ctx, "WHERE parent_asset_id = $1", 0, parentID)
}

// queryAssets selects whole assets. The condition is an optional WHERE clause
// and a limit above zero caps the number of rows.
func (s *Store) queryAssets(ctx context.Context, condition string, limit int, args ...any) ([]asset.Stored, error) {
	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", limit)
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT asset_id, parent_asset_id, asset_type, asset_name, operational_status,
			voltage_kv, rating_kva, manufacturer, model, serial_number,
			commissioned_date::text, created_at
		FROM assets `+condition+`
		ORDER BY asset_id `+limitClause, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assets := []asset.Stored{}
	for rows.Next() {
		var a asset.Stored
		err := rows.Scan(&a.ID, &a.ParentID, &a.Type, &a.Name, &a.Status,
			&a.VoltageKV, &a.RatingKVA, &a.Manufacturer, &a.Model, &a.SerialNumber,
			&a.CommissionedDate, &a.CreatedAt)
		if err != nil {
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
