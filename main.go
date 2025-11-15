package main

import (
	"os"

	"github.com/igromanas/go-final/pkg/db"
	"github.com/igromanas/go-final/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("TODO_PORT")
	dbPath := os.Getenv("TODO_DBFILE")

	db.Init(dbPath)
	server.Run(port)

	defer db.DB.Close()
}
