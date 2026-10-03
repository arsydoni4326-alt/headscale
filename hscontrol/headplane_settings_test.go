package hscontrol

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juanfont/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, InitHeadplaneSettings(db))
	return db
}

func TestEncryptDecryptAPIKey(t *testing.T) {
	apiKey := "test-api-key-12345"
	sessionToken := "test-session-token-abcdef"

	// Encrypt
	encrypted, nonce, salt, err := encryptAPIKey(apiKey, sessionToken)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEmpty(t, nonce)
	assert.NotEmpty(t, salt)

	// Decrypt
	decrypted, err := decryptAPIKey(encrypted, nonce, salt, sessionToken)
	require.NoError(t, err)
	assert.Equal(t, apiKey, decrypted)
}

func TestEncryptAPIKey_EmptyToken(t *testing.T) {
	_, _, _, err := encryptAPIKey("test-key", "")
	assert.ErrorIs(t, err, errSessionTokenRequired)
}

func TestDecryptAPIKey_WrongToken(t *testing.T) {
	apiKey := "test-api-key"
	sessionToken := "correct-token"

	encrypted, nonce, salt, err := encryptAPIKey(apiKey, sessionToken)
	require.NoError(t, err)

	// Try to decrypt with wrong token
	_, err = decryptAPIKey(encrypted, nonce, salt, "wrong-token")
	assert.ErrorIs(t, err, errDecryptionFailed)
}

func TestGetSettings_NotExists(t *testing.T) {
	db := setupTestDB(t)

	settings, err := GetSettings(db)
	assert.ErrorIs(t, err, errSettingsNotFound)
	assert.Nil(t, settings)
}

func TestUpdateSettings_Create(t *testing.T) {
	db := setupTestDB(t)

	err := UpdateSettings(db, "encrypted", "nonce", "salt", "dark", "John Doe")
	require.NoError(t, err)

	// Verify settings were created
	settings, err := GetSettings(db)
	require.NoError(t, err)
	assert.Equal(t, "encrypted", settings.APIKeyEncrypted)
	assert.Equal(t, "nonce", settings.APIKeyNonce)
	assert.Equal(t, "salt", settings.APIKeySalt)
	assert.Equal(t, "dark", settings.Theme)
	assert.Equal(t, "John Doe", settings.ProfileName)
}

func TestUpdateSettings_Update(t *testing.T) {
	db := setupTestDB(t)

	// Create initial settings
	err := UpdateSettings(db, "encrypted1", "nonce1", "salt1", "light", "Jane")
	require.NoError(t, err)

	// Update settings
	err = UpdateSettings(db, "encrypted2", "nonce2", "salt2", "dark", "Jane Doe")
	require.NoError(t, err)

	// Verify settings were updated (not duplicated)
	settings, err := GetSettings(db)
	require.NoError(t, err)
	assert.Equal(t, "encrypted2", settings.APIKeyEncrypted)
	assert.Equal(t, "nonce2", settings.APIKeyNonce)
	assert.Equal(t, "salt2", settings.APIKeySalt)
	assert.Equal(t, "dark", settings.Theme)
	assert.Equal(t, "Jane Doe", settings.ProfileName)

	// Verify only one row exists
	var count int64
	db.Model(&HeadplaneSettings{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestHandleGetSettings_Unauthorized(t *testing.T) {
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test-password",
		},
	}
	auth := NewHeadplaneAuth(cfg)

	app := &Headscale{
		headplaneAuth: auth,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/headplane/settings", nil)
	rec := httptest.NewRecorder()

	app.HandleGetSettings(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandleGetSettings_NotExists(t *testing.T) {
	db := setupTestDB(t)

	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test-password",
		},
	}
	auth := NewHeadplaneAuth(cfg)

	// Create a valid session
	token, err := auth.generateToken()
	require.NoError(t, err)
	auth.mu.Lock()
	auth.sessions[token] = &headplaneSession{
		Token:     token,
		CreatedAt: testNow(),
		ExpiresAt: testNow().Add(sessionDuration),
	}
	auth.mu.Unlock()

	// Mock state with DB
	mockState := &mockStateForSettings{db: db}
	app := &Headscale{
		headplaneAuth: auth,
		state:         mockState,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/headplane/settings", nil)
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	app.HandleGetSettings(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp SettingsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "", resp.APIKey)
	assert.Equal(t, "light", resp.Theme)
	assert.Equal(t, "", resp.ProfileName)
}

func TestHandleGetSettings_Success(t *testing.T) {
	db := setupTestDB(t)

	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test-password",
		},
	}
	auth := NewHeadplaneAuth(cfg)

	// Create a valid session
	token, err := auth.generateToken()
	require.NoError(t, err)
	auth.mu.Lock()
	auth.sessions[token] = &headplaneSession{
		Token:     token,
		CreatedAt: testNow(),
		ExpiresAt: testNow().Add(sessionDuration),
	}
	auth.mu.Unlock()

	// Store encrypted settings
	apiKey := "test-api-key-12345"
	encrypted, nonce, salt, err := encryptAPIKey(apiKey, token)
	require.NoError(t, err)
	err = UpdateSettings(db, encrypted, nonce, salt, "dark", "Test User")
	require.NoError(t, err)

	// Mock state with DB
	mockState := &mockStateForSettings{db: db}
	app := &Headscale{
		headplaneAuth: auth,
		state:         mockState,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/headplane/settings", nil)
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	app.HandleGetSettings(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp SettingsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, apiKey, resp.APIKey)
	assert.Equal(t, "dark", resp.Theme)
	assert.Equal(t, "Test User", resp.ProfileName)
}

func TestHandleUpdateSettings_Success(t *testing.T) {
	db := setupTestDB(t)

	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test-password",
		},
	}
	auth := NewHeadplaneAuth(cfg)

	// Create a valid session
	token, err := auth.generateToken()
	require.NoError(t, err)
	auth.mu.Lock()
	auth.sessions[token] = &headplaneSession{
		Token:     token,
		CreatedAt: testNow(),
		ExpiresAt: testNow().Add(sessionDuration),
	}
	auth.mu.Unlock()

	// Mock state with DB
	mockState := &mockStateForSettings{db: db}
	app := &Headscale{
		headplaneAuth: auth,
		state:         mockState,
	}

	apiKey := "new-api-key"
	theme := "dark"
	profileName := "John Doe"

	reqBody := UpdateSettingsRequest{
		APIKey:      &apiKey,
		Theme:       &theme,
		ProfileName: &profileName,
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/settings", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	app.HandleUpdateSettings(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]bool
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp["success"])

	// Verify settings were stored
	settings, err := GetSettings(db)
	require.NoError(t, err)
	assert.NotEmpty(t, settings.APIKeyEncrypted)
	assert.Equal(t, "dark", settings.Theme)
	assert.Equal(t, "John Doe", settings.ProfileName)

	// Verify decryption works
	decrypted, err := decryptAPIKey(settings.APIKeyEncrypted, settings.APIKeyNonce, settings.APIKeySalt, token)
	require.NoError(t, err)
	assert.Equal(t, apiKey, decrypted)
}

func TestHandleUpdateSettings_PartialUpdate(t *testing.T) {
	db := setupTestDB(t)

	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test-password",
		},
	}
	auth := NewHeadplaneAuth(cfg)

	// Create a valid session
	token, err := auth.generateToken()
	require.NoError(t, err)
	auth.mu.Lock()
	auth.sessions[token] = &headplaneSession{
		Token:     token,
		CreatedAt: testNow(),
		ExpiresAt: testNow().Add(sessionDuration),
	}
	auth.mu.Unlock()

	// Create initial settings
	apiKey := "initial-key"
	encrypted, nonce, salt, err := encryptAPIKey(apiKey, token)
	require.NoError(t, err)
	err = UpdateSettings(db, encrypted, nonce, salt, "light", "Initial Name")
	require.NoError(t, err)

	// Mock state with DB
	mockState := &mockStateForSettings{db: db}
	app := &Headscale{
		headplaneAuth: auth,
		state:         mockState,
	}

	// Update only theme
	theme := "dark"
	reqBody := UpdateSettingsRequest{
		Theme: &theme,
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/settings", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	app.HandleUpdateSettings(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify only theme changed, API key and profile name remain
	settings, err := GetSettings(db)
	require.NoError(t, err)
	assert.Equal(t, "dark", settings.Theme)
	assert.Equal(t, "Initial Name", settings.ProfileName)

	// Verify API key still decrypts correctly
	decrypted, err := decryptAPIKey(settings.APIKeyEncrypted, settings.APIKeyNonce, settings.APIKeySalt, token)
	require.NoError(t, err)
	assert.Equal(t, apiKey, decrypted)
}

func TestHandleChangePassword_Success(t *testing.T) {
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "old-password",
		},
	}
	auth := NewHeadplaneAuth(cfg)

	// Create a valid session
	token, err := auth.generateToken()
	require.NoError(t, err)
	auth.mu.Lock()
	auth.sessions[token] = &headplaneSession{
		Token:     token,
		CreatedAt: testNow(),
		ExpiresAt: testNow().Add(sessionDuration),
	}
	auth.mu.Unlock()

	app := &Headscale{
		cfg:           cfg,
		headplaneAuth: auth,
	}

	reqBody := ChangePasswordRequest{
		CurrentPassword: "old-password",
		NewPassword:     "new-password",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/change-password", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	app.HandleChangePassword(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]bool
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp["success"])

	// Verify password was changed
	assert.Equal(t, "new-password", cfg.Headplane.Password)
}

func TestHandleChangePassword_InvalidCurrent(t *testing.T) {
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "correct-password",
		},
	}
	auth := NewHeadplaneAuth(cfg)

	// Create a valid session
	token, err := auth.generateToken()
	require.NoError(t, err)
	auth.mu.Lock()
	auth.sessions[token] = &headplaneSession{
		Token:     token,
		CreatedAt: testNow(),
		ExpiresAt: testNow().Add(sessionDuration),
	}
	auth.mu.Unlock()

	app := &Headscale{
		cfg:           cfg,
		headplaneAuth: auth,
	}

	reqBody := ChangePasswordRequest{
		CurrentPassword: "wrong-password",
		NewPassword:     "new-password",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/change-password", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	app.HandleChangePassword(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Verify password was NOT changed
	assert.Equal(t, "correct-password", cfg.Headplane.Password)
}

func TestHandleChangePassword_EmptyNew(t *testing.T) {
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "old-password",
		},
	}
	auth := NewHeadplaneAuth(cfg)

	// Create a valid session
	token, err := auth.generateToken()
	require.NoError(t, err)
	auth.mu.Lock()
	auth.sessions[token] = &headplaneSession{
		Token:     token,
		CreatedAt: testNow(),
		ExpiresAt: testNow().Add(sessionDuration),
	}
	auth.mu.Unlock()

	app := &Headscale{
		cfg:           cfg,
		headplaneAuth: auth,
	}

	reqBody := ChangePasswordRequest{
		CurrentPassword: "old-password",
		NewPassword:     "",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/change-password", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()

	app.HandleChangePassword(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// Mock state for testing
type mockStateForSettings struct {
	db *gorm.DB
}

func (m *mockStateForSettings) DB() *gorm.DB {
	return m.db
}

// Helper to get a consistent test time
func testNow() time.Time {
	return time.Date(2026, 10, 3, 4, 32, 0, 0, time.UTC)
}
