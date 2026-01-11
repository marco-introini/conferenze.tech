package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func NewQuerier(db DBTX) Querier {
	return &querier{db: db}
}

type querier struct {
	db DBTX
}

func (q *querier) DeleteAllRegistrations(ctx context.Context) error {
	query := `DELETE FROM conference_registrations`
	_, err := q.db.ExecContext(ctx, query)
	return err
}

func (q *querier) DeleteAllConferences(ctx context.Context) error {
	query := `DELETE FROM conferences`
	_, err := q.db.ExecContext(ctx, query)
	return err
}

func (q *querier) DeleteAllUsers(ctx context.Context) error {
	query := `DELETE FROM users`
	_, err := q.db.ExecContext(ctx, query)
	return err
}

func (q *querier) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	query := `
		INSERT INTO users (email, password, name, nickname, city, avatar_url, bio)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at
	`
	var user User
	err := q.db.QueryRowContext(ctx, query,
		arg.Email, arg.Password, arg.Name, arg.Nickname, arg.City, arg.AvatarURL, arg.Bio,
	).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Nickname,
		&user.City, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (q *querier) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	query := `SELECT id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at FROM users WHERE id = $1`
	var user User
	err := q.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Nickname,
		&user.City, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (q *querier) GetUserByEmail(ctx context.Context, email string) (User, error) {
	query := `SELECT id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at FROM users WHERE email = $1`
	var user User
	err := q.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Nickname,
		&user.City, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (q *querier) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	query := `
		UPDATE users SET
			name = COALESCE($2, name),
			nickname = COALESCE($3, nickname),
			city = COALESCE($4, city),
			avatar_url = COALESCE($5, avatar_url),
			bio = COALESCE($6, bio),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at
	`
	var user User
	err := q.db.QueryRowContext(ctx, query,
		arg.ID, arg.Name, arg.Nickname, arg.City, arg.AvatarURL, arg.Bio,
	).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Nickname,
		&user.City, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (q *querier) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	query := `
		UPDATE users SET password = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at
	`
	var user User
	err := q.db.QueryRowContext(ctx, query, arg.ID, arg.Password).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Nickname,
		&user.City, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (q *querier) CreateConference(ctx context.Context, arg CreateConferenceParams) (Conference, error) {
	query := `
		INSERT INTO conferences (title, date, location, website, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, date, location, website, latitude, longitude, created_at, updated_at
	`
	var conf Conference
	err := q.db.QueryRowContext(ctx, query,
		arg.Title, arg.Date, arg.Location, arg.Website, arg.Latitude, arg.Longitude,
	).Scan(
		&conf.ID, &conf.Title, &conf.Date, &conf.Location, &conf.Website,
		&conf.Latitude, &conf.Longitude, &conf.CreatedAt, &conf.UpdatedAt,
	)
	return conf, err
}

func (q *querier) GetConferenceByID(ctx context.Context, id uuid.UUID) (Conference, error) {
	query := `SELECT id, title, date, location, website, latitude, longitude, created_at, updated_at FROM conferences WHERE id = $1`
	var conf Conference
	err := q.db.QueryRowContext(ctx, query, id).Scan(
		&conf.ID, &conf.Title, &conf.Date, &conf.Location, &conf.Website,
		&conf.Latitude, &conf.Longitude, &conf.CreatedAt, &conf.UpdatedAt,
	)
	return conf, err
}

func (q *querier) ListConferences(ctx context.Context) ([]Conference, error) {
	query := `SELECT id, title, date, location, website, latitude, longitude, created_at, updated_at FROM conferences ORDER BY date DESC`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Conference
	for rows.Next() {
		var i Conference
		if err := rows.Scan(
			&i.ID, &i.Title, &i.Date, &i.Location, &i.Website,
			&i.Latitude, &i.Longitude, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *querier) ListUpcomingConferences(ctx context.Context) ([]Conference, error) {
	query := `SELECT id, title, date, location, website, latitude, longitude, created_at, updated_at FROM conferences WHERE date >= NOW() ORDER BY date ASC`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Conference
	for rows.Next() {
		var i Conference
		if err := rows.Scan(
			&i.ID, &i.Title, &i.Date, &i.Location, &i.Website,
			&i.Latitude, &i.Longitude, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *querier) ListConferencesByLocation(ctx context.Context, location string) ([]Conference, error) {
	query := `SELECT id, title, date, location, website, latitude, longitude, created_at, updated_at FROM conferences WHERE location ILIKE $1 ORDER BY date DESC`
	rows, err := q.db.QueryContext(ctx, query, "%"+location+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Conference
	for rows.Next() {
		var i Conference
		if err := rows.Scan(
			&i.ID, &i.Title, &i.Date, &i.Location, &i.Website,
			&i.Latitude, &i.Longitude, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *querier) UpdateConference(ctx context.Context, arg UpdateConferenceParams) (Conference, error) {
	query := `
		UPDATE conferences SET
			title = COALESCE($2, title),
			date = COALESCE($3, date),
			location = COALESCE($4, location),
			website = COALESCE($5, website),
			latitude = COALESCE($6, latitude),
			longitude = COALESCE($7, longitude),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, date, location, website, latitude, longitude, created_at, updated_at
	`
	var conf Conference
	err := q.db.QueryRowContext(ctx, query,
		arg.ID, arg.Title, arg.Date, arg.Location, arg.Website, arg.Latitude, arg.Longitude,
	).Scan(
		&conf.ID, &conf.Title, &conf.Date, &conf.Location, &conf.Website,
		&conf.Latitude, &conf.Longitude, &conf.CreatedAt, &conf.UpdatedAt,
	)
	return conf, err
}

func (q *querier) DeleteConference(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM conferences WHERE id = $1`
	_, err := q.db.ExecContext(ctx, query, id)
	return err
}

func (q *querier) RegisterUserToConference(ctx context.Context, arg RegisterUserToConferenceParams) (ConferenceRegistration, error) {
	query := `
		INSERT INTO conference_registrations (user_id, conference_id, role, notes, needs_ride, has_car)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, conference_id, status, role, notes, needs_ride, has_car, registered_at, cancelled_at
	`
	var i ConferenceRegistration
	err := q.db.QueryRowContext(ctx, query,
		arg.UserID, arg.ConferenceID, arg.Role, arg.Notes, arg.NeedsRide, arg.HasCar,
	).Scan(
		&i.ID, &i.UserID, &i.ConferenceID, &i.Status, &i.Role,
		&i.Notes, &i.NeedsRide, &i.HasCar, &i.RegisteredAt, &i.CancelledAt,
	)
	return i, err
}

func (q *querier) GetRegistration(ctx context.Context, arg GetRegistrationParams) (ConferenceRegistration, error) {
	query := `SELECT id, user_id, conference_id, status, role, notes, needs_ride, has_car, registered_at, cancelled_at FROM conference_registrations WHERE user_id = $1 AND conference_id = $2`
	var i ConferenceRegistration
	err := q.db.QueryRowContext(ctx, query, arg.UserID, arg.ConferenceID).Scan(
		&i.ID, &i.UserID, &i.ConferenceID, &i.Status, &i.Role,
		&i.Notes, &i.NeedsRide, &i.HasCar, &i.RegisteredAt, &i.CancelledAt,
	)
	return i, err
}

func (q *querier) GetRegistrationsByConference(ctx context.Context, conferenceID uuid.UUID) ([]ConferenceWithUser, error) {
	query := `
		SELECT r.id, r.user_id, r.conference_id, r.status, r.role, r.notes, r.needs_ride, r.has_car, r.registered_at, r.cancelled_at,
		       u.email, u.name, u.nickname, u.city, u.avatar_url
		FROM conference_registrations r
		JOIN users u ON u.id = r.user_id
		WHERE r.conference_id = $1
	`
	rows, err := q.db.QueryContext(ctx, query, conferenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ConferenceWithUser
	for rows.Next() {
		var i ConferenceWithUser
		if err := rows.Scan(
			&i.ID, &i.UserID, &i.ConferenceID, &i.Status, &i.Role, &i.Notes, &i.NeedsRide, &i.HasCar, &i.RegisteredAt, &i.CancelledAt,
			&i.Email, &i.UserName, &i.Nickname, &i.UserCity, &i.AvatarURL,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *querier) GetRegistrationsByUser(ctx context.Context, userID uuid.UUID) ([]RegistrationWithConference, error) {
	query := `
		SELECT r.id, r.user_id, r.conference_id, r.status, r.role, r.notes, r.needs_ride, r.has_car, r.registered_at, r.cancelled_at,
		       c.title, c.date, c.location, c.website
		FROM conference_registrations r
		JOIN conferences c ON c.id = r.conference_id
		WHERE r.user_id = $1
	`
	rows, err := q.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RegistrationWithConference
	for rows.Next() {
		var i RegistrationWithConference
		if err := rows.Scan(
			&i.ID, &i.UserID, &i.ConferenceID, &i.Status, &i.Role, &i.Notes, &i.NeedsRide, &i.HasCar, &i.RegisteredAt, &i.CancelledAt,
			&i.Title, &i.ConfDate, &i.Location, &i.Website,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *querier) UpdateRegistrationStatus(ctx context.Context, arg UpdateRegistrationStatusParams) (ConferenceRegistration, error) {
	query := `
		UPDATE conference_registrations SET status = $2
		WHERE id = $1
		RETURNING id, user_id, conference_id, status, role, notes, needs_ride, has_car, registered_at, cancelled_at
	`
	var i ConferenceRegistration
	err := q.db.QueryRowContext(ctx, query, arg.ID, arg.Status).Scan(
		&i.ID, &i.UserID, &i.ConferenceID, &i.Status, &i.Role,
		&i.Notes, &i.NeedsRide, &i.HasCar, &i.RegisteredAt, &i.CancelledAt,
	)
	return i, err
}

func (q *querier) CancelRegistration(ctx context.Context, arg CancelRegistrationParams) (ConferenceRegistration, error) {
	query := `
		UPDATE conference_registrations SET status = 'cancelled', cancelled_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, conference_id, status, role, notes, needs_ride, has_car, registered_at, cancelled_at
	`
	var i ConferenceRegistration
	err := q.db.QueryRowContext(ctx, query, arg.ID).Scan(
		&i.ID, &i.UserID, &i.ConferenceID, &i.Status, &i.Role,
		&i.Notes, &i.NeedsRide, &i.HasCar, &i.RegisteredAt, &i.CancelledAt,
	)
	return i, err
}

func (q *querier) DeleteRegistration(ctx context.Context, arg DeleteRegistrationParams) error {
	query := `DELETE FROM conference_registrations WHERE user_id = $1 AND conference_id = $2`
	_, err := q.db.ExecContext(ctx, query, arg.UserID, arg.ConferenceID)
	return err
}

func (q *querier) ListUsersNeedingRide(ctx context.Context, conferenceID uuid.UUID) ([]CarpoolUser, error) {
	query := `
		SELECT u.id, u.email, u.password, u.name, u.nickname, u.city, u.avatar_url, u.bio, u.created_at, u.updated_at,
		       c.title, c.location, r.notes
		FROM conference_registrations r
		JOIN users u ON u.id = r.user_id
		JOIN conferences c ON c.id = r.conference_id
		WHERE r.conference_id = $1 AND r.needs_ride = TRUE AND r.status != 'cancelled'
	`
	rows, err := q.db.QueryContext(ctx, query, conferenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CarpoolUser
	for rows.Next() {
		var i CarpoolUser
		if err := rows.Scan(
			&i.ID, &i.Email, &i.Password, &i.Name, &i.Nickname, &i.City, &i.AvatarURL, &i.Bio, &i.CreatedAt, &i.UpdatedAt,
			&i.ConferenceTitle, &i.ConferenceLocation, &i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *querier) ListUsersOfferingRide(ctx context.Context, conferenceID uuid.UUID) ([]CarpoolUser, error) {
	query := `
		SELECT u.id, u.email, u.password, u.name, u.nickname, u.city, u.avatar_url, u.bio, u.created_at, u.updated_at,
		       c.title, c.location, r.notes
		FROM conference_registrations r
		JOIN users u ON u.id = r.user_id
		JOIN conferences c ON c.id = r.conference_id
		WHERE r.conference_id = $1 AND r.has_car = TRUE AND r.status != 'cancelled'
	`
	rows, err := q.db.QueryContext(ctx, query, conferenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CarpoolUser
	for rows.Next() {
		var i CarpoolUser
		if err := rows.Scan(
			&i.ID, &i.Email, &i.Password, &i.Name, &i.Nickname, &i.City, &i.AvatarURL, &i.Bio, &i.CreatedAt, &i.UpdatedAt,
			&i.ConferenceTitle, &i.ConferenceLocation, &i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *querier) GetConferenceStats(ctx context.Context, id uuid.UUID) (ConferenceStats, error) {
	query := `
		SELECT
			c.id, c.title,
			COUNT(r.id) as total_registrations,
			COUNT(r.id) FILTER (WHERE r.status = 'registered') as confirmed_count,
			COUNT(r.id) FILTER (WHERE r.needs_ride = TRUE AND r.status != 'cancelled') as needing_ride_count,
			COUNT(r.id) FILTER (WHERE r.has_car = TRUE AND r.status != 'cancelled') as offering_ride_count
		FROM conferences c
		LEFT JOIN conference_registrations r ON r.conference_id = c.id
		WHERE c.id = $1
		GROUP BY c.id, c.title
	`
	var i ConferenceStats
	err := q.db.QueryRowContext(ctx, query, id).Scan(
		&i.ID, &i.Title, &i.TotalRegistrations, &i.ConfirmedCount, &i.NeedingRideCount, &i.OfferingRideCount,
	)
	return i, err
}
