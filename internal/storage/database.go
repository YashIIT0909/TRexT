package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/YashIIT0909/TRexT/internal/storage/db"
	"github.com/YashIIT0909/TRexT/sql/schemas"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

// DB wraps the PostgreSQL database connection and sqlc queries
type DB struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewDB creates a new database connection and runs migrations
func NewDB() (*DB, error) {
	// Load .env file if it exists (ignores error if file doesn't exist)
	_ = godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required. Create a .env file with DATABASE_URL=postgres://user:password@host:port/dbname?sslmode=disable")
	}

	ctx := context.Background()

	// Create connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run goose migrations using embedded SQL files
	// Goose needs a standard database/sql connection for migrations
	if err := runMigrations(connStr); err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{
		pool:    pool,
		queries: db.New(pool),
	}, nil
}

// runMigrations runs goose migrations using database/sql
func runMigrations(connStr string) error {
	goose.SetBaseFS(schemas.EmbedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Use stdlib adapter for goose (it requires database/sql)
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open sql connection for migrations: %w", err)
	}
	defer sqlDB.Close()

	if err := goose.Up(sqlDB, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Close closes the database connection pool
func (d *DB) Close() {
	d.pool.Close()
}

// Queries returns the sqlc generated queries
func (d *DB) Queries() *db.Queries {
	return d.queries
}

// Pool returns the underlying connection pool
func (d *DB) Pool() *pgxpool.Pool {
	return d.pool
}

// SaveRequest saves or updates a request
func (d *DB) SaveRequest(req *SavedRequest) error {
	ctx := context.Background()

	if req.ID == 0 {
		collectionID := pgtype.Int4{Int32: int32(req.CollectionID), Valid: req.CollectionID > 0}
		result, err := d.queries.CreateRequest(ctx, db.CreateRequestParams{
			Name:         req.Name,
			Url:          req.URL,
			Method:       req.Method,
			Headers:      pgtype.Text{String: req.Headers, Valid: true},
			Body:         pgtype.Text{String: req.Body, Valid: true},
			CollectionID: collectionID,
		})
		if err != nil {
			return err
		}
		req.ID = int64(result.ID)
	} else {
		collectionID := pgtype.Int4{Int32: int32(req.CollectionID), Valid: req.CollectionID > 0}
		err := d.queries.UpdateRequest(ctx, db.UpdateRequestParams{
			Name:         req.Name,
			Url:          req.URL,
			Method:       req.Method,
			Headers:      pgtype.Text{String: req.Headers, Valid: true},
			Body:         pgtype.Text{String: req.Body, Valid: true},
			CollectionID: collectionID,
			ID:           int32(req.ID),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllRequests returns all saved requests
func (d *DB) GetAllRequests() ([]*SavedRequest, error) {
	ctx := context.Background()
	rows, err := d.queries.GetAllRequests(ctx)
	if err != nil {
		return nil, err
	}

	requests := make([]*SavedRequest, len(rows))
	for i, row := range rows {
		requests[i] = &SavedRequest{
			ID:           int64(row.ID),
			Name:         row.Name,
			URL:          row.Url,
			Method:       row.Method,
			Headers:      row.Headers.String,
			Body:         row.Body.String,
			CollectionID: int64(row.CollectionID),
		}
	}
	return requests, nil
}

// DeleteRequest deletes a request by ID
func (d *DB) DeleteRequest(id int64) error {
	return d.queries.DeleteRequest(context.Background(), int32(id))
}

// GetCollections returns all collections
func (d *DB) GetCollections() ([]*Collection, error) {
	ctx := context.Background()
	rows, err := d.queries.GetCollections(ctx)
	if err != nil {
		return nil, err
	}

	collections := make([]*Collection, len(rows))
	for i, row := range rows {
		collections[i] = &Collection{
			ID:          int64(row.ID),
			Name:        row.Name,
			Description: row.Description.String,
		}
	}
	return collections, nil
}

// SaveCollection saves or updates a collection
func (d *DB) SaveCollection(c *Collection) error {
	ctx := context.Background()

	if c.ID == 0 {
		result, err := d.queries.CreateCollection(ctx, db.CreateCollectionParams{
			Name:        c.Name,
			Description: pgtype.Text{String: c.Description, Valid: true},
		})
		if err != nil {
			return err
		}
		c.ID = int64(result.ID)
	} else {
		err := d.queries.UpdateCollection(ctx, db.UpdateCollectionParams{
			Name:        c.Name,
			Description: pgtype.Text{String: c.Description, Valid: true},
			ID:          int32(c.ID),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// AddToHistory adds a request to history
func (d *DB) AddToHistory(entry *HistoryEntry) error {
	ctx := context.Background()
	result, err := d.queries.AddToHistory(ctx, db.AddToHistoryParams{
		Url:        entry.URL,
		Method:     entry.Method,
		StatusCode: pgtype.Int4{Int32: int32(entry.StatusCode), Valid: true},
		DurationMs: pgtype.Int8{Int64: entry.Duration, Valid: true},
		Timestamp:  entry.Timestamp,
	})
	if err != nil {
		return err
	}
	entry.ID = int64(result.ID)
	return nil
}

// GetHistory returns recent history entries
func (d *DB) GetHistory(limit int) ([]*HistoryEntry, error) {
	ctx := context.Background()
	rows, err := d.queries.GetHistory(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	history := make([]*HistoryEntry, len(rows))
	for i, row := range rows {
		history[i] = &HistoryEntry{
			ID:         int64(row.ID),
			URL:        row.Url,
			Method:     row.Method,
			StatusCode: int(row.StatusCode.Int32),
			Duration:   row.DurationMs.Int64,
			Timestamp:  row.Timestamp,
		}
	}
	return history, nil
}

// Ensure stdlib driver is registered
var _ = stdlib.GetDefaultDriver()
