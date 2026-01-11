package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/marco-introini/conferenze.tech/backend/db"
)

type Server struct {
	db *db.DB
}

func NewServer(database *db.DB) *Server {
	return &Server{db: database}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}

type RegisterRequest struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	Nickname  *string `json:"nickname"`
	City      *string `json:"city"`
	AvatarURL *string `json:"avatarUrl"`
	Bio       *string `json:"bio"`
}

type RegisterResponse struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}

type ConferenceResponse struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Location  string   `json:"location"`
	Website   *string  `json:"website,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type CreateConferenceRequest struct {
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Location  string   `json:"location"`
	Website   *string  `json:"website"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type RegisterToConferenceRequest struct {
	ConferenceID string  `json:"conferenceId"`
	Role         string  `json:"role"`
	Notes        *string `json:"notes"`
	NeedsRide    bool    `json:"needsRide"`
	HasCar       bool    `json:"hasCar"`
}

func generateToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password and name are required", http.StatusBadRequest)
		return
	}

	existingUser, err := s.db.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != uuid.Nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	passwordHash, err := db.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := s.db.CreateUser(ctx, db.CreateUserParams{
		Email:     req.Email,
		Password:  passwordHash,
		Name:      req.Name,
		Nickname:  req.Nickname,
		City:      req.City,
		AvatarUrl: req.AvatarURL,
		Bio:       req.Bio,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{User: user, Token: token})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !db.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{User: user, Token: token})
}

func (s *Server) ListConferences(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	conferences, err := s.db.ListConferences(ctx)
	if err != nil {
		log.Printf("Error listing conferences: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]ConferenceResponse, len(conferences))
	for i, c := range conferences {
		dateStr := c.Date.Format(time.RFC3339)
		response[i] = ConferenceResponse{
			ID:        c.ID.String(),
			Title:     c.Title,
			Date:      dateStr,
			Location:  c.Location,
			Website:   c.Website,
			Latitude:  c.Latitude,
			Longitude: c.Longitude,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) GetConference(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Conference ID required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid conference ID", http.StatusBadRequest)
		return
	}

	conference, err := s.db.GetConferenceByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Conference not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting conference: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	registrations, err := s.db.GetRegistrationsByConference(ctx, id)
	if err != nil {
		log.Printf("Error getting registrations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := ConferenceWithAttendees{
		ID:        conference.ID.String(),
		Title:     conference.Title,
		Date:      conference.Date.Format(time.RFC3339),
		Location:  conference.Location,
		Website:   conference.Website,
		Latitude:  conference.Latitude,
		Longitude: conference.Longitude,
		Attendees: make([]Attendee, 0),
	}

	for _, reg := range registrations {
		attendee := Attendee{
			User: UserResponse{
				ID:       reg.UserID.String(),
				Nickname: reg.Nickname,
				City:     reg.City,
			},
			NeedsRide: reg.NeedsRide,
			HasCar:    reg.HasCar,
		}
		response.Attendees = append(response.Attendees, attendee)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) CreateConference(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req CreateConferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Location == "" || req.Date == "" {
		http.Error(w, "Title, location and date are required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	conference, err := s.db.CreateConference(ctx, db.CreateConferenceParams{
		Title:     req.Title,
		Date:      date,
		Location:  req.Location,
		Website:   req.Website,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		log.Printf("Error creating conference: %v", err)
		http.Error(w, "Failed to create conference", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ConferenceResponse{
		ID:        conference.ID.String(),
		Title:     conference.Title,
		Date:      conference.Date.Format(time.RFC3339),
		Location:  conference.Location,
		Website:   conference.Website,
		Latitude:  conference.Latitude,
		Longitude: conference.Longitude,
	})
}

func (s *Server) RegisterToConference(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req RegisterToConferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil || userID == uuid.Nil {
		http.Error(w, "Valid user ID required", http.StatusBadRequest)
		return
	}

	conferenceID, err := uuid.Parse(req.ConferenceID)
	if err != nil {
		http.Error(w, "Invalid conference ID", http.StatusBadRequest)
		return
	}

	_, err = s.db.GetConferenceByID(ctx, conferenceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Conference not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting conference: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	role := db.UserRole(req.Role)
	if role != db.RoleAttendee && role != db.RoleSpeaker && role != db.RoleVolunteer {
		role = db.RoleAttendee
	}

	registration, err := s.db.RegisterUserToConference(ctx, db.RegisterUserToConferenceParams{
		UserID:       userID,
		ConferenceID: conferenceID,
		Role:         role,
		Notes:        req.Notes,
		NeedsRide:    req.NeedsRide,
		HasCar:       req.HasCar,
	})
	if err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Failed to register to conference", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registration)
}

func (s *Server) GetUserRegistrations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	registrations, err := s.db.GetRegistrationsByUser(ctx, userID)
	if err != nil {
		log.Printf("Error getting registrations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]RegistrationResponse, len(registrations))
	for i, reg := range registrations {
		response[i] = RegistrationResponse{
			ID:                 reg.ID.String(),
			ConferenceID:       reg.ConferenceID.String(),
			ConferenceTitle:    reg.Title,
			ConferenceDate:     reg.Date.Format(time.RFC3339),
			ConferenceLocation: reg.Location,
			Status:             string(reg.Status),
			Role:               string(reg.Role),
			NeedsRide:          reg.NeedsRide,
			HasCar:             reg.HasCar,
			RegisteredAt:       reg.RegisteredAt.Format(time.RFC3339),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type ConferenceWithAttendees struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Date      string     `json:"date"`
	Location  string     `json:"location"`
	Website   *string    `json:"website,omitempty"`
	Latitude  *float64   `json:"latitude,omitempty"`
	Longitude *float64   `json:"longitude,omitempty"`
	Attendees []Attendee `json:"attendees"`
}

type Attendee struct {
	User      UserResponse `json:"user"`
	NeedsRide bool         `json:"needsRide"`
	HasCar    bool         `json:"hasCar"`
}

type UserResponse struct {
	ID       string  `json:"id"`
	Nickname *string `json:"nickname,omitempty"`
	City     *string `json:"city,omitempty"`
}

type RegistrationResponse struct {
	ID                 string `json:"id"`
	ConferenceID       string `json:"conferenceId"`
	ConferenceTitle    string `json:"conferenceTitle"`
	ConferenceDate     string `json:"conferenceDate"`
	ConferenceLocation string `json:"conferenceLocation"`
	Status             string `json:"status"`
	Role               string `json:"role"`
	NeedsRide          bool   `json:"needsRide"`
	HasCar             bool   `json:"hasCar"`
	RegisteredAt       string `json:"registeredAt"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Run(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", s.Register)
	mux.HandleFunc("/api/login", s.Login)
	mux.HandleFunc("/api/conferences", s.ListConferences)
	mux.HandleFunc("/api/conference", s.GetConference)
	mux.HandleFunc("/api/conferences/create", s.CreateConference)
	mux.HandleFunc("/api/register-to-conference", s.RegisterToConference)
	mux.HandleFunc("/api/my-registrations", s.GetUserRegistrations)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	})

	handler := loggingMiddleware(corsMiddleware(mux))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, handler)
}
