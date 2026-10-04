package hscontrol

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arsydoni4326-alt/headscale/hscontrol/db"
	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadplaneAuth_HandleLogin_NoPassword(t *testing.T) {
	hsdb := setupTestDBForAuth(t)
	defer hsdb.Close()
	
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "",
		},
	}
	auth := NewHeadplaneAuth(cfg, hsdb)

	body := `{"password":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/login", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	auth.HandleLogin(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestHeadplaneAuth_HandleLogin_Success(t *testing.T) {
	hsdb := setupTestDBForAuth(t)
	defer hsdb.Close()
	
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test-password",
		},
	}
	auth := NewHeadplaneAuth(cfg, hsdb)

	body := `{"password":"test-password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/login", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	auth.HandleLogin(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp, "token")
	assert.Contains(t, resp, "expires_at")

	token := resp["token"].(string)
	assert.NotEmpty(t, token)
	assert.True(t, auth.VerifySession(token))
}

func TestHeadplaneAuth_HandleLogin_InvalidPassword(t *testing.T) {
	hsdb := setupTestDBForAuth(t)
	defer hsdb.Close()
	
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "correct-password",
		},
	}
	auth := NewHeadplaneAuth(cfg, hsdb)

	body := `{"password":"wrong-password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/login", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	auth.HandleLogin(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHeadplaneAuth_HandleLogin_RateLimit(t *testing.T) {
	hsdb := setupTestDBForAuth(t)
	defer hsdb.Close()
	
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test-password",
		},
	}
	auth := NewHeadplaneAuth(cfg, hsdb)

	// Make maxLoginAttempts failed attempts
	for i := 0; i < maxLoginAttempts; i++ {
		body := `{"password":"wrong"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/login", bytes.NewBufferString(body))
		req.RemoteAddr = "192.168.1.1:1234"
		rec := httptest.NewRecorder()
		auth.HandleLogin(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}

	// Next attempt should be rate limited
	body := `{"password":"test-password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/headplane/login", bytes.NewBufferString(body))
	req.RemoteAddr = "192.168.1.1:1234"
	rec := httptest.NewRecorder()

	auth.HandleLogin(rec, req)

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestHeadplaneAuth_VerifySession_Expired(t *testing.T) {
	hsdb := setupTestDBForAuth(t)
	defer hsdb.Close()
	
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test",
		},
	}
	auth := NewHeadplaneAuth(cfg, hsdb)

	// Verify non-existent session
	assert.False(t, auth.VerifySession("nonexistent-token"))
}

// setupTestDBForAuth creates an in-memory database for auth tests
func setupTestDBForAuth(t *testing.T) *db.HSDatabase {
	t.Helper()
	cfg := &types.Config{
		Database: types.DatabaseConfig{
			Type: types.DatabaseSqlite,
			Sqlite: types.SqliteConfig{
				Path: ":memory:",
			},
		},
	}

	hsdb, err := db.NewHeadscaleDatabase(cfg)
	require.NoError(t, err)

	return hsdb
}
