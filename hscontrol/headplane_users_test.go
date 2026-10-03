package hscontrol

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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

func TestHandleGetUser(t *testing.T) {
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	cfg := &types.Config{}
	auth := NewHeadplaneAuth(cfg, hsdb)
	app := &Headscale{
		cfg:           cfg,
		headplaneAuth: auth,
	}
	app.state = &mockState{db: hsdb}

	adminToken := loginAndGetToken(t, auth, hsdb, "admin", "adminpass", "admin")
	user, _ := db.CreateHeadplaneUser(hsdb.DB, "testuser", "password", "user")

	tests := []struct {
		name           string
		userID         string
		token          string
		expectedStatus int
	}{
		{
			name:           "get user successfully",
			userID:         "2",
			token:          adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "unauthorized without token",
			userID:         "2",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "user not found",
			userID:         "999",
			token:          adminToken,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/headplane/users/"+tt.userID, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			app.HandleGetUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp HeadplaneUserResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, user.Username, resp.Username)
			}
		})
	}
}

func TestHandleUpdateUser(t *testing.T) {
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	cfg := &types.Config{}
	auth := NewHeadplaneAuth(cfg, hsdb)
	app := &Headscale{
		cfg:           cfg,
		headplaneAuth: auth,
	}
	app.state = &mockState{db: hsdb}

	adminToken := loginAndGetToken(t, auth, hsdb, "admin", "adminpass", "admin")
	user, _ := db.CreateHeadplaneUser(hsdb.DB, "testuser", "password", "user")

	t.Run("update username successfully", func(t *testing.T) {
		reqBody := UpdateUserRequest{Username: "newusername"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/headplane/users/2", bytes.NewReader(body))
		req.Header.Set("Authorization", adminToken)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		app.HandleUpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp HeadplaneUserResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "newusername", resp.Username)
	})

	t.Run("update role successfully", func(t *testing.T) {
		reqBody := UpdateUserRequest{Role: "admin"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/headplane/users/2", bytes.NewReader(body))
		req.Header.Set("Authorization", adminToken)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		app.HandleUpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandleDeleteUser_LastAdminProtection(t *testing.T) {
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	cfg := &types.Config{}
	auth := NewHeadplaneAuth(cfg, hsdb)
	app := &Headscale{
		cfg:           cfg,
		headplaneAuth: auth,
	}
	app.state = &mockState{db: hsdb}

	adminToken := loginAndGetToken(t, auth, hsdb, "admin", "adminpass", "admin")

	// Try to delete the only admin (should fail)
	req := httptest.NewRequest("DELETE", "/api/v1/headplane/users/1", nil)
	req.Header.Set("Authorization", adminToken)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	app.HandleDeleteUser(w, req)

	// Should fail because admin can't delete themselves
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "own account")
}

func TestHandleDeleteUser_RegularUser(t *testing.T) {
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	cfg := &types.Config{}
	auth := NewHeadplaneAuth(cfg, hsdb)
	app := &Headscale{
		cfg:           cfg,
		headplaneAuth: auth,
	}
	app.state = &mockState{db: hsdb}

	adminToken := loginAndGetToken(t, auth, hsdb, "admin", "adminpass", "admin")
	user, _ := db.CreateHeadplaneUser(hsdb.DB, "testuser", "password", "user")

	// Admin can delete regular user
	req := httptest.NewRequest("DELETE", "/api/v1/headplane/users/2", nil)
	req.Header.Set("Authorization", adminToken)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	app.HandleDeleteUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify user is deleted
	_, err := db.GetHeadplaneUserByID(hsdb.DB, user.ID)
	assert.Equal(t, db.ErrHeadplaneUserNotFound, err)
}

func loginAndGetToken(t *testing.T, auth *HeadplaneAuth, hsdb *db.HSDatabase, username, password, role string) string {
	// Create user first
	db.CreateHeadplaneUser(hsdb.DB, username, password, role)

	reqBody := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/headplane/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	auth.HandleLogin(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["token"].(string)
}

// mockState is a minimal State mock for testing
type mockState struct {
	db *db.HSDatabase
}

func (m *mockState) DB() *db.HSDatabase {
	return m.db
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
