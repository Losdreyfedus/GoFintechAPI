package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type Database struct {
	*sql.DB
}

// addDefaultConnectionParams ensures TLS-related params are set for local Docker SQL Server.
// If not explicitly provided, we default to a local-dev friendly configuration.
func addDefaultConnectionParams(connectionString string) string {
	parsed, err := url.Parse(connectionString)
	if err != nil {
		// If parsing fails, just return the original string to avoid making things worse.
		return connectionString
	}

	query := parsed.Query()
	// detect presence of keys case-insensitively
	hasEncrypt := false
	hasTsc := false
	for key := range query {
		lower := strings.ToLower(key)
		if lower == "encrypt" {
			hasEncrypt = true
		}
		if lower == "trustservercertificate" {
			hasTsc = true
		}
	}

	if !hasEncrypt {
		// For local containers a common need is to disable encryption
		query.Set("encrypt", "disable")
	}
	// TrustServerCertificate is only relevant when encrypt=true, but keeping it set
	// avoids surprises if the user flips encrypt later.
	if !hasTsc {
		query.Set("trustservercertificate", "true")
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

// buildMasterConnectionString returns a connection string that targets the master database
// while preserving credentials and host parameters.
func buildMasterConnectionString(connectionString string) string {
	parsed, err := url.Parse(connectionString)
	if err != nil {
		return connectionString
	}
	query := parsed.Query()
	query.Set("database", "master")
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

// extractDatabaseName reads the "database" query parameter from the connection string.
func extractDatabaseName(connectionString string) (string, bool) {
	parsed, err := url.Parse(connectionString)
	if err != nil {
		return "", false
	}
	dbName := parsed.Query().Get("database")
	if strings.TrimSpace(dbName) == "" {
		return "", false
	}
	return dbName, true
}

// ensureDatabaseExists connects to master and creates the target database if it does not exist.
func ensureDatabaseExists(connectionString string, timeout time.Duration) error {
	dbName, ok := extractDatabaseName(connectionString)
	if !ok {
		return nil
	}
	masterConn := buildMasterConnectionString(connectionString)

	deadline := time.Now().Add(timeout)
	var lastErr error
	for attempt := 0; time.Now().Before(deadline); attempt++ {
		masterDB, err := sql.Open("mssql", masterConn)
		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}
		if pingErr := masterDB.Ping(); pingErr != nil {
			lastErr = pingErr
			_ = masterDB.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		// Create database if not exists
		_, err = masterDB.Exec(
			"IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = '" + dbName + "') BEGIN CREATE DATABASE [" + dbName + "] END",
		)
		closeErr := masterDB.Close()
		if err == nil && closeErr == nil {
			return nil
		}
		if err != nil {
			lastErr = err
		} else {
			lastErr = closeErr
		}
		time.Sleep(2 * time.Second)
	}
	if lastErr == nil {
		lastErr = errors.New("ensureDatabaseExists: timeout")
	}
	return lastErr
}

func NewConnection(connectionString string) (*Database, error) {
	// Be resilient to Docker SQL Server startup timing and TLS defaults
	conn := addDefaultConnectionParams(connectionString)

	// Make sure the target database exists (will wait for server readiness)
	_ = ensureDatabaseExists(conn, 60*time.Second)

	// Try connecting with retry/backoff while the container is warming up
	var (
		db  *sql.DB
		err error
	)
	wait := 1 * time.Second
	for attempt := 0; attempt < 15; attempt++ { // ~2 minutes total with backoff
		db, err = sql.Open("mssql", conn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Database not ready yet (attempt %d): %v", attempt+1, err)
		time.Sleep(wait)
		if wait < 10*time.Second {
			wait *= 2
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	var dbName string
	if err := db.QueryRow("SELECT DB_NAME()").Scan(&dbName); err == nil {
		log.Printf("Connected to database: %s\n", dbName)
	} else {
		log.Printf("Could not get database name: %v\n", err)
	}

	log.Println("Database connection established successfully")
	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
