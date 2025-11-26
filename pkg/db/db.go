package db

import (
	"context"
	"database/sql"
	"log"

	"os"

	_ "modernc.org/sqlite"
)

const tableSchema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
	title VARCHAR(255) NOT NULL,
	comment TEXT,
    repeat VARCHAR(128)
	);`

const indexSchema = `CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`

var DB *sql.DB

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("DB connection error: %v", err) // todo
		return err
	}

	// DB.SetMaxOpenConns(1) // sqlite?

	if install {
		_, err = DB.ExecContext(context.Background(), tableSchema)
	}

	if err != nil {
		log.Printf("schema execution error: %v", err)
		return err
	}

	if install {
		_, err = DB.ExecContext(context.Background(), indexSchema)
	}

	if err != nil {
		log.Printf("schema execution error: %v", err)
		return err
	}

	return nil
}
