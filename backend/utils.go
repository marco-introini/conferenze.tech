package main

import (
	"database/sql"
	"time"
)

// nullString converts a string pointer to sql.NullString
func nullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// nullFloat64 converts a float64 pointer to sql.NullFloat64
func nullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

// nullBool converts a bool to sql.NullBool
func nullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// stringPtr converts sql.NullString to a string pointer
func stringPtr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// float64Ptr converts sql.NullFloat64 to a float64 pointer
func float64Ptr(f sql.NullFloat64) *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// boolPtr converts sql.NullBool to a bool pointer
func boolPtr(b sql.NullBool) *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// timePtr converts sql.NullTime to a string pointer (RFC3339 format)
func timePtr(t sql.NullTime) *string {
	if !t.Valid {
		return nil
	}
	s := t.Time.Format(time.RFC3339)
	return &s
}
