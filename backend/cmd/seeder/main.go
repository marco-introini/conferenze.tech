package main

import (
	"context"
	"log"
	"os"

	"github.com/marco-introini/conferenze.tech/backend/db"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/conferenzetech?sslmode=disable"
	}

	database, err := db.New(dsn)
	if err != nil {
		log.Fatalf("impossibile connettersi al database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := db.Seed(ctx, database); err != nil {
		log.Fatalf("errore durante il seeding: %v", err)
	}

	log.Println("Database popolato con successo!")
}
