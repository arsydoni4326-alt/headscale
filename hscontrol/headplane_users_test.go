package hscontrol

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juanfont/headscale/hscontrol/db"
	"github.com/juanfont/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiUserLogin(t *testing.T) {
	cfg := &types.Config{
		Headplane: types.HeadplaneConfig{
			Password: "test-password",
		},
	}

	hsdb := setupTestDB(t)
	defer hsdb.Close()

	auth := NewHeadplaneAuth(cfg, hsdb)

	// Create a test user
	user, err := db.CreateHeadplaneUser(hsdb.DB, "testuser", "testpass123", "user")
	require.NoError(t, err)
	require.NotNil(t, user)

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
		expectToken    bool
	}{
		{
			name:           "valid credentials",
			username:       "testuser",
			password:       "testpass123",
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name:           "invalid username",
			username:       "wronguser",
			password:       "testpass123",
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name:           "invalid password",
			username:       "testuser",
			password:       "wrongpass",
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]string{
				"username": tt.username,
				"password": tt.password,
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/v1/headplane/login", bytes.NewReader(body))
			w := httptest.NewRecorder()

			auth.HandleLogin(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectToken {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp["token"])
				assert.Equal(t, "testuser", resp["username"])
			}
		})
	}
}

func TestPasswordHashing(t *testing.T) {
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	user, err := db.CreateHeadplaneUser(hsdb.DB, "testuser", "mypassword", "user")
	require.NoError(t, err)

	// Password should be hashed, not stored in plain text
	assert.NotEqual(t, "mypassword", user.PasswordHash)
	assert.NotEmpty(t, user.PasswordHash)

	// CheckPassword should work with correct password
	assert.True(t, user.CheckPassword("mypassword"))

	// CheckPassword should fail with incorrect password
	assert.False(t, user.CheckPassword("wrongpassword"))
}

func TestSessionManagement(t *testing.T) {
	cfg := &types.Config{}
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	auth := NewHeadplaneAuth(cfg, hsdb)

	// Create a test user
	user, err := db.CreateHeadplaneUser(hsdb.DB, "testuser", "testpass", "admin")
	require.NoError(t, err)

	// Simulate login
	reqBody := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/headplane/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	auth.HandleLogin(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	token := resp["token"].(string)
	assert.NotEmpty(t, token)

	// Verify session
	assert.True(t, auth.VerifySession(token))

	// Get session details
	session, ok := auth.GetSession(token)
	assert.True(t, ok)
	assert.Equal(t, user.ID, session.UserID)
	assert.Equal(t, "testuser", session.Username)
	assert.True(t, session.IsAdmin)

	// Invalid token should fail
	assert.False(t, auth.VerifySession("invalid-token"))
}

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *db.HSDatabase {
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

