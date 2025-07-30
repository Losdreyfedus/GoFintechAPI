package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type Database struct {
	*sql.DB
}

func NewConnection(connectionString string) (*Database, error) {
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Bağlanılan veritabanı adını logla
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
