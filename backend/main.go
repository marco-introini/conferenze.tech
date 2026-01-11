package main

import (
	"log"
	"os"

	"github.com/marco-introini/backend-conferenze/db"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/conferenzetech?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database, err := db.New(dsn)
	if err != nil {
		log.Fatalf("impossibile connettersi al database: %v", err)
	}
	defer database.Close()

	server := NewServer(database)
	if err := server.Run(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
