# API Documentation

## Base URL
`/api`

## Public Routes (No Authentication Required)

### Register User
- **Endpoint:** `POST /api/register`
- **Description:** Register a new user

### Login
- **Endpoint:** `POST /api/login`
- **Description:** Authenticate user and receive tokens

### List Conferences
- **Endpoint:** `GET /api/conferences`
- **Description:** Retrieve all conferences

### Get Conference Details
- **Endpoint:** `GET /api/conferences/{conference_id}`
- **Description:** Retrieve details for a specific conference

---

## Protected Routes (Authentication Required)

### Create Conference
- **Endpoint:** `POST /api/conferences`
- **Description:** Create a new conference

### Delete Conference
- **Endpoint:** `DELETE /api/conferences/{conference_id}`
- **Description:** Delete a specific conference

### Register to Conference
- **Endpoint:** `POST /api/conferences/{conference_id}/register`
- **Description:** Register the authenticated user to a conference

### Get User Registrations
- **Endpoint:** `GET /api/users/registrations`
- **Description:** Retrieve all conference registrations for the authenticated user

### Unregister from Conference
- **Endpoint:** `DELETE /api/users/registrations/{conference_id}`
- **Description:** Cancel registration to a specific conference

### Get User Profile
- **Endpoint:** `GET /api/users/{user_id}`
- **Description:** Retrieve user profile by ID

### Get Current User
- **Endpoint:** `GET /api/me`
- **Description:** Retrieve the authenticated user's profile

### List Tokens
- **Endpoint:** `GET /api/tokens`
- **Description:** Retrieve all tokens for the authenticated user

### Revoke Token
- **Endpoint:** `POST /api/tokens/revoke`
- **Description:** Revoke a specific token

---

## Health Check

### Health
- **Endpoint:** `GET /health`
- **Description:** Server health check endpoint
