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

// TestMultiUserRegistrationFlow tests the complete user registration flow.
func TestMultiUserRegistrationFlow(t *testing.T) {
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	cfg := &types.Config{}
	auth := NewHeadplaneAuth(cfg, hsdb)
	app := &Headscale{
		cfg:           cfg,
		headplaneAuth: auth,
	}
	app.state = &mockState{db: hsdb}

	// Create admin user
	admin, err := db.CreateHeadplaneUser(hsdb.DB, "admin", "adminpass", "admin")
	require.NoError(t, err)
	adminToken := loginAndGetToken(t, auth, hsdb, "admin", "adminpass", "admin")

	// Test 1: Admin can register new user
	t.Run("admin can register new user", func(t *testing.T) {
		reqBody := RegisterUserRequest{
			Username: "newuser",
			Password: "newpass123",
			Role:     "user",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/headplane/users", bytes.NewReader(body))
		req.Header.Set("Authorization", adminToken)
		w := httptest.NewRecorder()

		app.HandleRegisterUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp HeadplaneUserResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "newuser", resp.Username)
		assert.Equal(t, "user", resp.Role)
	})

	// Test 2: Regular user cannot register new users
	t.Run("regular user cannot register new users", func(t *testing.T) {
		userToken := loginAndGetToken(t, auth, hsdb, "newuser", "newpass123", "user")

		reqBody := RegisterUserRequest{
			Username: "anotheruser",
			Password: "pass123",
			Role:     "user",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/headplane/users", bytes.NewReader(body))
		req.Header.Set("Authorization", userToken)
		w := httptest.NewRecorder()

		app.HandleRegisterUser(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test 3: Duplicate username should fail
	t.Run("duplicate username should fail", func(t *testing.T) {
		reqBody := RegisterUserRequest{
			Username: "newuser",
			Password: "pass123",
			Role:     "user",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/headplane/users", bytes.NewReader(body))
		req.Header.Set("Authorization", adminToken)
		w := httptest.NewRecorder()

		app.HandleRegisterUser(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// Test 4: Verify admin user exists
	user, err := db.GetHeadplaneUserByUsername(hsdb.DB, "admin")
	require.NoError(t, err)
	assert.Equal(t, admin.ID, user.ID)
}

func stringPtr(s string) *string {
	return &s
}
