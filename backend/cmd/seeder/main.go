package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/marco-introini/conferenze.tech/backend/db"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/conferenzetech?sslmode=disable"
	}

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("impossibile connettersi al database: %v", err)
	}
	defer sqlDB.Close()

	queries := db.New(sqlDB)
	database := db.WrapDB(queries)

	ctx := context.Background()
	if err := db.Seed(ctx, database); err != nil {
		log.Fatalf("errore durante il seeding: %v", err)
	}

	log.Println("Database popolato con successo!")
}
