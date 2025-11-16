package main

import (
	"log"
	"os"

	"github.com/igromanas/go-final/pkg/db"
	"github.com/igromanas/go-final/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("TODO_PORT")
	dbPath := os.Getenv("TODO_DBFILE")

	err := db.Init(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	if err := server.Run(port); err != nil {
		log.Fatal(err)
	}
}
