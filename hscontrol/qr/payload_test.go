package qr

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistrationPayload_Valid(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)

	payload, err := NewRegistrationPayload(authID, serverURL, expiresAt)

	require.NoError(t, err)
	assert.Equal(t, "headscale-registration", payload.Type)
	assert.Equal(t, "1", payload.Version)
	assert.Equal(t, authID, payload.AuthID)
	assert.Equal(t, serverURL, payload.ServerURL)
	assert.Equal(t, expiresAt, payload.ExpiresAt)
}

func TestNewRegistrationPayload_EmptyAuthID(t *testing.T) {
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)

	payload, err := NewRegistrationPayload("", serverURL, expiresAt)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "auth_id cannot be empty")
}

func TestNewRegistrationPayload_EmptyServerURL(t *testing.T) {
	authID := "hskey-abc123"
	expiresAt := time.Now().Add(24 * time.Hour)

	payload, err := NewRegistrationPayload(authID, "", expiresAt)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "server_url cannot be empty")
}

func TestNewRegistrationPayload_ZeroExpiry(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"

	payload, err := NewRegistrationPayload(authID, serverURL, time.Time{})

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "expires_at cannot be zero")
}

func TestNewRegistrationPayload_PastExpiry(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(-24 * time.Hour)

	payload, err := NewRegistrationPayload(authID, serverURL, expiresAt)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "expires_at must be in the future")
}

func TestPayload_ToJSON(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	payload, err := NewRegistrationPayload(authID, serverURL, expiresAt)
	require.NoError(t, err)

	jsonData, err := payload.ToJSON()
	require.NoError(t, err)

	// Verify JSON structure
	var decoded map[string]interface{}
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "headscale-registration", decoded["type"])
	assert.Equal(t, "1", decoded["version"])
	assert.Equal(t, authID, decoded["auth_id"])
	assert.Equal(t, serverURL, decoded["server_url"])
	assert.Equal(t, "2027-01-01T00:00:00Z", decoded["expires_at"])
}

func TestParseRegistrationPayload_Valid(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.NoError(t, err)
	assert.Equal(t, "headscale-registration", payload.Type)
	assert.Equal(t, "1", payload.Version)
	assert.Equal(t, "hskey-abc123", payload.AuthID)
	assert.Equal(t, "https://headscale.example.com", payload.ServerURL)
	assert.Equal(
		t,
		time.Date(2026, 10, 8, 7, 0, 0, 0, time.UTC),
		payload.ExpiresAt,
	)
}

func TestParseRegistrationPayload_InvalidJSON(t *testing.T) {
	jsonData := []byte(`{invalid json}`)

	payload, err := ParseRegistrationPayload(jsonData)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "failed to parse payload")
}

func TestParseRegistrationPayload_InvalidType(t *testing.T) {
	jsonData := []byte(`{
		"type": "wrong-type",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "invalid payload type")
}

func TestParseRegistrationPayload_InvalidVersion(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "2",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "unsupported payload version")
}

func TestParseRegistrationPayload_EmptyAuthID(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "",
		"server_url": "https://headscale.example.com",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "auth_id cannot be empty")
}

func TestParseRegistrationPayload_EmptyServerURL(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "server_url cannot be empty")
}

func TestParseRegistrationPayload_ZeroExpiry(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "0001-01-01T00:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "expires_at cannot be zero")
}

func TestPayload_RoundTrip(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	// Create original payload
	original, err := NewRegistrationPayload(authID, serverURL, expiresAt)
	require.NoError(t, err)

	// Serialize to JSON
	jsonData, err := original.ToJSON()
	require.NoError(t, err)

	// Parse back
	parsed, err := ParseRegistrationPayload(jsonData)
	require.NoError(t, err)

	// Verify all fields match
	assert.Equal(t, original.Type, parsed.Type)
	assert.Equal(t, original.Version, parsed.Version)
	assert.Equal(t, original.AuthID, parsed.AuthID)
	assert.Equal(t, original.ServerURL, parsed.ServerURL)
	assert.Equal(t, original.ExpiresAt, parsed.ExpiresAt)
}
