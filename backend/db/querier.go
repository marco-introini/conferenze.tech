package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Querier interface {
	DeleteAllRegistrations(ctx context.Context) error
	DeleteAllConferences(ctx context.Context) error
	DeleteAllUsers(ctx context.Context) error
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error)
	CreateConference(ctx context.Context, arg CreateConferenceParams) (Conference, error)
	GetConferenceByID(ctx context.Context, id uuid.UUID) (Conference, error)
	ListConferences(ctx context.Context) ([]Conference, error)
	ListUpcomingConferences(ctx context.Context) ([]Conference, error)
	ListConferencesByLocation(ctx context.Context, location string) ([]Conference, error)
	UpdateConference(ctx context.Context, arg UpdateConferenceParams) (Conference, error)
	DeleteConference(ctx context.Context, id uuid.UUID) error
	RegisterUserToConference(ctx context.Context, arg RegisterUserToConferenceParams) (ConferenceRegistration, error)
	GetRegistration(ctx context.Context, arg GetRegistrationParams) (ConferenceRegistration, error)
	GetRegistrationsByConference(ctx context.Context, conferenceID uuid.UUID) ([]ConferenceWithUser, error)
	GetRegistrationsByUser(ctx context.Context, userID uuid.UUID) ([]RegistrationWithConference, error)
	UpdateRegistrationStatus(ctx context.Context, arg UpdateRegistrationStatusParams) (ConferenceRegistration, error)
	CancelRegistration(ctx context.Context, arg CancelRegistrationParams) (ConferenceRegistration, error)
	DeleteRegistration(ctx context.Context, arg DeleteRegistrationParams) error
	ListUsersNeedingRide(ctx context.Context, conferenceID uuid.UUID) ([]CarpoolUser, error)
	ListUsersOfferingRide(ctx context.Context, conferenceID uuid.UUID) ([]CarpoolUser, error)
	GetConferenceStats(ctx context.Context, id uuid.UUID) (ConferenceStats, error)
}

type CreateUserParams struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	Nickname  *string `json:"nickname"`
	City      *string `json:"city"`
	AvatarURL *string `json:"avatarUrl"`
	Bio       *string `json:"bio"`
}

type UpdateUserParams struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name"`
	Nickname  *string   `json:"nickname"`
	City      *string   `json:"city"`
	AvatarURL *string   `json:"avatarUrl"`
	Bio       *string   `json:"bio"`
}

type UpdateUserPasswordParams struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

type CreateConferenceParams struct {
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	Location  string    `json:"location"`
	Website   *string   `json:"website"`
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
}

type UpdateConferenceParams struct {
	ID        uuid.UUID  `json:"id"`
	Title     *string    `json:"title"`
	Date      *time.Time `json:"date"`
	Location  *string    `json:"location"`
	Website   *string    `json:"website"`
	Latitude  *float64   `json:"latitude"`
	Longitude *float64   `json:"longitude"`
}

type RegisterUserToConferenceParams struct {
	UserID       uuid.UUID `json:"userId"`
	ConferenceID uuid.UUID `json:"conferenceId"`
	Role         UserRole  `json:"role"`
	Notes        *string   `json:"notes"`
	NeedsRide    bool      `json:"needsRide"`
	HasCar       bool      `json:"hasCar"`
}

type GetRegistrationParams struct {
	UserID       uuid.UUID `json:"userId"`
	ConferenceID uuid.UUID `json:"conferenceId"`
}

type UpdateRegistrationStatusParams struct {
	ID     uuid.UUID          `json:"id"`
	Status RegistrationStatus `json:"status"`
}

type CancelRegistrationParams struct {
	ID uuid.UUID `json:"id"`
}

type DeleteRegistrationParams struct {
	UserID       uuid.UUID `json:"userId"`
	ConferenceID uuid.UUID `json:"conferenceId"`
}
