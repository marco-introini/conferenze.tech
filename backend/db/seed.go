package db

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// Seed popola il database con dati di esempio utilizzando gofakeit
func Seed(ctx context.Context, database *DB) error {
	return database.WithTransaction(ctx, func(q Querier) error {
		fmt.Println("Inizio seeding con dati casuali...")

		// 1. Creazione Utenti
		numUsers := 10
		createdUsers := make([]User, 0, numUsers)
		for i := 0; i < numUsers; i++ {
			u := CreateUserParams{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 12),
				Name:     gofakeit.Name(),
				Nickname: stringPtr(gofakeit.Username()),
				City:     stringPtr(gofakeit.City()),
				Bio:      stringPtr(gofakeit.Sentence(10)),
			}

			user, err := q.CreateUser(ctx, u)
			if err != nil {
				return fmt.Errorf("errore creazione utente %s: %w", u.Email, err)
			}
			createdUsers = append(createdUsers, user)
			fmt.Printf("Utente creato: %s\n", user.Email)
		}

		// 2. Creazione Conferenze
		numConferences := 5
		createdConfs := make([]Conference, 0, numConferences)
		for i := 0; i < numConferences; i++ {
			c := CreateConferenceParams{
				Title:     gofakeit.Company() + " Conf " + fmt.Sprint(2025+i),
				Date:      gofakeit.DateRange(time.Now(), time.Now().AddDate(1, 0, 0)),
				Location:  gofakeit.Address().City + ", " + gofakeit.Address().Country,
				Website:   stringPtr(gofakeit.URL()),
				Latitude:  float64Ptr(gofakeit.Latitude()),
				Longitude: float64Ptr(gofakeit.Longitude()),
			}

			conf, err := q.CreateConference(ctx, c)
			if err != nil {
				return fmt.Errorf("errore creazione conferenza %s: %w", c.Title, err)
			}
			createdConfs = append(createdConfs, conf)
			fmt.Printf("Conferenza creata: %s\n", conf.Title)
		}

		// 3. Creazione Registrazioni
		roles := []UserRole{RoleAttendee, RoleOrganizer, RoleSpeaker, RoleVolunteer}
		for _, user := range createdUsers {
			// Registra ogni utente a 1-3 conferenze casuali
			numRegs := rand.Intn(3) + 1
			perm := rand.Perm(len(createdConfs))
			for i := 0; i < numRegs && i < len(createdConfs); i++ {
				conf := createdConfs[perm[i]]
				r := RegisterUserToConferenceParams{
					UserID:       user.ID,
					ConferenceID: conf.ID,
					Role:         roles[rand.Intn(len(roles))],
					NeedsRide:    gofakeit.Bool(),
					HasCar:       gofakeit.Bool(),
					Notes:        stringPtr(gofakeit.Phrase()),
				}

				reg, err := q.RegisterUserToConference(ctx, r)
				if err != nil {
					// Ignoriamo errori di duplicati (UNIQUE constraint) se capitano
					continue
				}
				fmt.Printf("Registrazione creata: User %s -> Conf %s\n", reg.UserID, reg.ConferenceID)
			}
		}

		fmt.Println("Seeding completato con successo!")
		return nil
	})
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
