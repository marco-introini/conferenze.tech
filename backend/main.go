package main

import (
	"context"
	"log"
	"os"

	"github.com/marco-introini/backend-conferenze/db"
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

	// Esempio: lista conferenze
	conferences, err := database.ListUpcomingConferences(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Conferenze in programma:")
	for _, c := range conferences {
		log.Printf("- %s (%s) il %v\n", c.Title, c.Location, c.Date.Format("02/01/2006"))
	}
}
