-- name: DeleteAllRegistrations :exec
DELETE FROM conference_registrations;

-- name: DeleteAllConferences :exec
DELETE FROM conferences;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: CreateUser :one
INSERT INTO users (email, password, name, nickname, city, avatar_url, bio)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET
    name = COALESCE($2, name),
    nickname = COALESCE($3, nickname),
    city = COALESCE($4, city),
    avatar_url = COALESCE($5, avatar_url),
    bio = COALESCE($6, bio),
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at;

-- name: UpdateUserPassword :one
UPDATE users SET password = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, email, password, name, nickname, city, avatar_url, bio, created_at, updated_at;

-- name: CreateConference :one
INSERT INTO conferences (title, date, location, website, latitude, longitude)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, title, date, location, website, latitude, longitude, created_at, updated_at;

-- name: GetConferenceByID :one
SELECT id, title, date, location, website, latitude, longitude, created_at, updated_at FROM conferences WHERE id = $1;

-- name: ListConferences :many
SELECT id, title, date, location, website, latitude, longitude, created_at, updated_at FROM conferences ORDER BY date DESC;

-- name: ListUpcomingConferences :many
SELECT id, title, date, location, website, latitude, longitude, created_at, updated_at FROM conferences WHERE date >= NOW() ORDER BY date ASC;

-- name: ListConferencesByLocation :many
SELECT id, title, date, location, website, latitude, longitude, created_at, updated_at FROM conferences WHERE location ILIKE $1 ORDER BY date DESC;

-- name: UpdateConference :one
UPDATE conferences SET
    title = COALESCE($2, title),
    date = COALESCE($3, date),
    location = COALESCE($4, location),
    website = COALESCE($5, website),
    latitude = COALESCE($6, latitude),
    longitude = COALESCE($7, longitude),
    updated_at = NOW()
WHERE id = $1
RETURNING id, title, date, location, website, latitude, longitude, created_at, updated_at;

-- name: DeleteConference :exec
DELETE FROM conferences WHERE id = $1;

-- name: RegisterUserToConference :one
INSERT INTO conference_registrations (user_id, conference_id, role, notes, needs_ride, has_car)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, conference_id, status, role, notes, needs_ride, has_car, registered_at, cancelled_at;

-- name: GetRegistration :one
SELECT id, user_id, conference_id, status, role, notes, needs_ride, has_car, registered_at, cancelled_at FROM conference_registrations WHERE user_id = $1 AND conference_id = $2;

-- name: GetRegistrationsByConference :many
SELECT r.id, r.user_id, r.conference_id, r.status, r.role, r.notes, r.needs_ride, r.has_car, r.registered_at, r.cancelled_at,
       u.email, u.name, u.nickname, u.city, u.avatar_url
FROM conference_registrations r
JOIN users u ON u.id = r.user_id
WHERE r.conference_id = $1;

-- name: GetRegistrationsByUser :many
SELECT r.id, r.user_id, r.conference_id, r.status, r.role, r.notes, r.needs_ride, r.has_car, r.registered_at, r.cancelled_at,
       c.title, c.date, c.location, c.website
FROM conference_registrations r
JOIN conferences c ON c.id = r.conference_id
WHERE r.user_id = $1;

-- name: UpdateRegistrationStatus :one
UPDATE conference_registrations SET status = $2
WHERE id = $1
RETURNING id, user_id, conference_id, status, role, notes, needs_ride, has_car, registered_at, cancelled_at;

-- name: CancelRegistration :one
UPDATE conference_registrations SET status = 'cancelled', cancelled_at = NOW()
WHERE id = $1
RETURNING id, user_id, conference_id, status, role, notes, needs_ride, has_car, registered_at, cancelled_at;

-- name: DeleteRegistration :exec
DELETE FROM conference_registrations WHERE user_id = $1 AND conference_id = $2;

-- name: ListUsersNeedingRide :many
SELECT u.id, u.email, u.password, u.name, u.nickname, u.city, u.avatar_url, u.bio, u.created_at, u.updated_at,
       c.title, c.location, r.notes
FROM conference_registrations r
JOIN users u ON u.id = r.user_id
JOIN conferences c ON c.id = r.conference_id
WHERE r.conference_id = $1 AND r.needs_ride = TRUE AND r.status != 'cancelled';

-- name: ListUsersOfferingRide :many
SELECT u.id, u.email, u.password, u.name, u.nickname, u.city, u.avatar_url, u.bio, u.created_at, u.updated_at,
       c.title, c.location, r.notes
FROM conference_registrations r
JOIN users u ON u.id = r.user_id
JOIN conferences c ON c.id = r.conference_id
WHERE r.conference_id = $1 AND r.has_car = TRUE AND r.status != 'cancelled';

-- name: GetConferenceStats :one
SELECT
    c.id, c.title,
    COUNT(r.id) as total_registrations,
    COUNT(r.id) FILTER (WHERE r.status = 'registered') as confirmed_count,
    COUNT(r.id) FILTER (WHERE r.needs_ride = TRUE AND r.status != 'cancelled') as needing_ride_count,
    COUNT(r.id) FILTER (WHERE r.has_car = TRUE AND r.status != 'cancelled') as offering_ride_count
FROM conferences c
LEFT JOIN conference_registrations r ON r.conference_id = c.id
WHERE c.id = $1
GROUP BY c.id, c.title;

-- Token management queries (user_tokens table must be present in schema.sql)
-- name: CreateToken :one
INSERT INTO user_tokens (user_id, token_hash)
VALUES ($1, $2)
RETURNING id, user_id, token_hash, created_at, last_used_at, revoked;

-- name: GetTokensByUser :many
SELECT id, user_id, token_hash, created_at, last_used_at, revoked
FROM user_tokens
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetTokenByHash :one
SELECT id, user_id, token_hash, created_at, last_used_at, revoked
FROM user_tokens
WHERE token_hash = $1;

-- name: DeleteToken :exec
DELETE FROM user_tokens WHERE id = $1;

-- name: RevokeToken :one
UPDATE user_tokens SET revoked = true WHERE id = $1 RETURNING id, user_id, token_hash, created_at, last_used_at, revoked;
