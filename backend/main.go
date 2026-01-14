package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/marco-introini/conferenze.tech/backend/db"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = DefaultDSN
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("impossibile connettersi al database: %v", err)
	}
	defer sqlDB.Close()

	queries := db.New(sqlDB)
	server := NewServer(queries)
	if err := server.Run(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
