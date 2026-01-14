package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: true}
}

func nullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// Seed popola il database con dati di esempio utilizzando gofakeit
func Seed(ctx context.Context, database interface {
	WithTransaction(context.Context, func(Querier) error) error
}) error {
	return database.WithTransaction(ctx, func(q Querier) error {
		fmt.Println("Pulizia database...")
		if err := q.DeleteAllRegistrations(ctx); err != nil {
			return fmt.Errorf("errore pulizia registrazioni: %w", err)
		}
		if err := q.DeleteAllConferences(ctx); err != nil {
			return fmt.Errorf("errore pulizia conferenze: %w", err)
		}
		if err := q.DeleteAllUsers(ctx); err != nil {
			return fmt.Errorf("errore pulizia utenti: %w", err)
		}

		fmt.Println("Inizio seeding con i miei dati specifici...")

		hashedPassword, err := HashPassword("password")
		if err != nil {
			return fmt.Errorf("errore hashing password: %w", err)
		}
		u := CreateUserParams{
			Email:    "marco@marcointroini.it",
			Password: hashedPassword,
			Name:     "marco",
			Nickname: nullString("marco"),
			City:     nullString("Lissone"),
			Bio:      nullString("Bio di Marco"),
		}

		user, err := q.CreateUser(ctx, u)
		if err != nil {
			return fmt.Errorf("errore creazione utente %s: %w", u.Email, err)
		}
		fmt.Printf("Utente creato: %s\n", user.Email)

		fmt.Println("Inizio seeding con dati casuali...")

		numUsers := 10
		createdUsers := make([]User, 0, numUsers)
		for i := 0; i < numUsers; i++ {
			rawPassword := gofakeit.Password(true, true, true, true, false, 12)
			hashedPassword, err := HashPassword(rawPassword)
			if err != nil {
				return fmt.Errorf("errore hashing password random: %w", err)
			}
			u := CreateUserParams{
				Email:    gofakeit.Email(),
				Password: hashedPassword,
				Name:     gofakeit.Name(),
				Nickname: nullString(gofakeit.Username()),
				City:     nullString(gofakeit.City()),
				Bio:      nullString(gofakeit.Sentence(10)),
			}

			user, err := q.CreateUser(ctx, u)
			if err != nil {
				return fmt.Errorf("errore creazione utente %s: %w", u.Email, err)
			}
			createdUsers = append(createdUsers, user)
			fmt.Printf("Utente creato: %s\n", user.Email)
		}

		numConferences := 5
		createdConfs := make([]Conference, 0, numConferences)
		for i := 0; i < numConferences; i++ {
			addr := gofakeit.Address()
			website := gofakeit.URL()
			c := CreateConferenceParams{
				Title:     gofakeit.Company() + " Conf " + fmt.Sprint(2025+i),
				Date:      gofakeit.DateRange(time.Now(), time.Now().AddDate(1, 0, 0)),
				Location:  addr.City + ", " + addr.Country,
				Website:   nullString(website),
				Latitude:  nullFloat64(gofakeit.Latitude()),
				Longitude: nullFloat64(gofakeit.Longitude()),
			}

			conf, err := q.CreateConference(ctx, c)
			if err != nil {
				return fmt.Errorf("errore creazione conferenza %s: %w", c.Title, err)
			}
			createdConfs = append(createdConfs, conf)
			fmt.Printf("Conferenza creata: %s\n", conf.Title)
		}

		roles := []string{"attendee", "organizer", "speaker", "volunteer"}
		for _, user := range createdUsers {
			if len(createdConfs) == 0 {
				break
			}
			numRegs := rand.Intn(3) + 1
			perm := rand.Perm(len(createdConfs))
			for i := 0; i < numRegs && i < len(createdConfs); i++ {
				conf := createdConfs[perm[i]]
				r := RegisterUserToConferenceParams{
					UserID:       user.ID,
					ConferenceID: conf.ID,
					Role:         roles[rand.Intn(len(roles))],
					Notes:        nullString(gofakeit.Phrase()),
					NeedsRide:    nullBool(gofakeit.Bool()),
					HasCar:       nullBool(gofakeit.Bool()),
				}

				reg, err := q.RegisterUserToConference(ctx, r)
				if err != nil {
					continue
				}
				fmt.Printf("Registrazione creata: User %v -> Conf %v\n", reg.UserID, reg.ConferenceID)
			}
		}

		fmt.Println("Seeding completato con successo!")
		return nil
	})
}
