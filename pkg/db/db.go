package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"os"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
	title VARCHAR(255) NOT NULL,
	comment TEXT,
    repeat VARCHAR(128)
	);`

const schema2 = `CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`

// const schema3 = `SELECT name FROM sqlite_master WHERE type='table';`

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

	DB.SetMaxOpenConns(1)
	fmt.Printf("ping: %v\n", DB.Ping())

	if install {
		_, err = DB.ExecContext(context.Background(), schema)
	}

	if err != nil {
		log.Printf("schema execution error: {%v}", err)
		return err
	}

	if install {
		_, err = DB.ExecContext(context.Background(), schema2)
	}

	if err != nil {
		log.Printf("schema execution error: {%v}", err)
		return err
	}

	// check
	// rows, _ := DB.QueryContext(context.Background(), schema3)
	// defer rows.Close()
	// for rows.Next() {
	// 	var name string
	// 	_ = rows.Scan(&name)
	// 	fmt.Println("row:", name)
	// }

	return nil
}
