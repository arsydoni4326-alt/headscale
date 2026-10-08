package qr

import (
	"encoding/json"
	"fmt"
	"time"
)

// RegistrationPayload represents the JSON payload encoded in a QR code
// for headscale node registration.
type RegistrationPayload struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	AuthID    string    `json:"auth_id"`
	ServerURL string    `json:"server_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewRegistrationPayload creates and validates a new registration payload.
// It enforces non-empty fields and ensures the expiry is in the future.
func NewRegistrationPayload(
	authID, serverURL string,
	expiresAt time.Time,
) (*RegistrationPayload, error) {
	if authID == "" {
		return nil, fmt.Errorf("auth_id cannot be empty")
	}
	if serverURL == "" {
		return nil, fmt.Errorf("server_url cannot be empty")
	}
	if expiresAt.IsZero() {
		return nil, fmt.Errorf("expires_at cannot be zero")
	}
	if expiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("expires_at must be in the future")
	}

	return &RegistrationPayload{
		Type:      "headscale-registration",
		Version:   "1",
		AuthID:    authID,
		ServerURL: serverURL,
		ExpiresAt: expiresAt,
	}, nil
}

// ToJSON serializes the payload to JSON bytes.
func (p *RegistrationPayload) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// ParseRegistrationPayload deserializes and validates a registration payload
// from JSON bytes.
func ParseRegistrationPayload(data []byte) (*RegistrationPayload, error) {
	var payload RegistrationPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	// Validate the parsed payload
	if payload.Type != "headscale-registration" {
		return nil, fmt.Errorf("invalid payload type: %s", payload.Type)
	}
	if payload.Version != "1" {
		return nil, fmt.Errorf("unsupported payload version: %s", payload.Version)
	}
	if payload.AuthID == "" {
		return nil, fmt.Errorf("auth_id cannot be empty")
	}
	if payload.ServerURL == "" {
		return nil, fmt.Errorf("server_url cannot be empty")
	}
	if payload.ExpiresAt.IsZero() {
		return nil, fmt.Errorf("expires_at cannot be zero")
	}
	if payload.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("payload has expired")
	}

	return &payload, nil
}
