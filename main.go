package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/igromanas/go-final/pkg/db"
	"github.com/igromanas/go-final/pkg/server"
)

func setEnvs() {
	// TODO_PORT
	os.Setenv("TODO_PORT", ":7540")

	// TODO_DBFILE
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("%v", err)
	}
	fmt.Println("wd:", wd)
	p, err := filepath.Abs(wd)
	if err != nil {
		log.Printf("%v", err)
	}
	fmt.Println("p:", p)

	dbEnv := filepath.Join(p, "/pkg/db/scheduler.db")
	os.Setenv("TODO_DBFILE", dbEnv)
}

func getDBFile() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "./pkg/db/scheduler.db"
	}
	return dbFile
}

func main() {
	db.Init(getDBFile())
	server.Run()

	// -- time testing --
	// res, err := api.NextDate(time.Now(), "20251110", "d 7")
	// fmt.Printf("res: %s\nerr: %w\n", res, err)
	defer db.DB.Close()
}
