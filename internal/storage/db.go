package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection
func NewDB() (*DB, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

// getDBPath returns the path to the SQLite database
func getDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "trext", "data.db"), nil
}

// migrate creates the database schema
func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS collections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		method TEXT NOT NULL DEFAULT 'GET',
		headers TEXT DEFAULT '{}',
		body TEXT DEFAULT '',
		collection_id INTEGER,
		FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		method TEXT NOT NULL,
		status_code INTEGER,
		duration_ms INTEGER,
		timestamp INTEGER NOT NULL
	);

	-- Insert default collection if not exists
	INSERT OR IGNORE INTO collections (id, name, description) 
	VALUES (1, 'Default', 'Default collection');
	`

	_, err := db.conn.Exec(schema)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// SaveRequest saves or updates a request
func (db *DB) SaveRequest(req *SavedRequest) error {
	if req.ID == 0 {
		result, err := db.conn.Exec(
			`INSERT INTO requests (name, url, method, headers, body, collection_id) 
			 VALUES (?, ?, ?, ?, ?, ?)`,
			req.Name, req.URL, req.Method, req.Headers, req.Body, req.CollectionID,
		)
		if err != nil {
			return err
		}
		req.ID, _ = result.LastInsertId()
	} else {
		_, err := db.conn.Exec(
			`UPDATE requests SET name=?, url=?, method=?, headers=?, body=?, collection_id=? 
			 WHERE id=?`,
			req.Name, req.URL, req.Method, req.Headers, req.Body, req.CollectionID, req.ID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllRequests returns all saved requests
func (db *DB) GetAllRequests() ([]*SavedRequest, error) {
	rows, err := db.conn.Query(
		`SELECT id, name, url, method, headers, body, COALESCE(collection_id, 1) 
		 FROM requests ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*SavedRequest
	for rows.Next() {
		req := &SavedRequest{}
		if err := rows.Scan(&req.ID, &req.Name, &req.URL, &req.Method, &req.Headers, &req.Body, &req.CollectionID); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

// DeleteRequest deletes a request by ID
func (db *DB) DeleteRequest(id int64) error {
	_, err := db.conn.Exec(`DELETE FROM requests WHERE id=?`, id)
	return err
}

// GetCollections returns all collections
func (db *DB) GetCollections() ([]*Collection, error) {
	rows, err := db.conn.Query(`SELECT id, name, description FROM collections ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*Collection
	for rows.Next() {
		c := &Collection{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}
	return collections, nil
}

// SaveCollection saves or updates a collection
func (db *DB) SaveCollection(c *Collection) error {
	if c.ID == 0 {
		result, err := db.conn.Exec(
			`INSERT INTO collections (name, description) VALUES (?, ?)`,
			c.Name, c.Description,
		)
		if err != nil {
			return err
		}
		c.ID, _ = result.LastInsertId()
	} else {
		_, err := db.conn.Exec(
			`UPDATE collections SET name=?, description=? WHERE id=?`,
			c.Name, c.Description, c.ID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddToHistory adds a request to history
func (db *DB) AddToHistory(entry *HistoryEntry) error {
	_, err := db.conn.Exec(
		`INSERT INTO history (url, method, status_code, duration_ms, timestamp) 
		 VALUES (?, ?, ?, ?, ?)`,
		entry.URL, entry.Method, entry.StatusCode, entry.Duration, entry.Timestamp,
	)
	return err
}

// GetHistory returns recent history entries
func (db *DB) GetHistory(limit int) ([]*HistoryEntry, error) {
	rows, err := db.conn.Query(
		`SELECT id, url, method, status_code, duration_ms, timestamp 
		 FROM history ORDER BY timestamp DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*HistoryEntry
	for rows.Next() {
		h := &HistoryEntry{}
		if err := rows.Scan(&h.ID, &h.URL, &h.Method, &h.StatusCode, &h.Duration, &h.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}
