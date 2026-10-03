package hscontrol

import (
	"testing"
	"time"

	"github.com/juanfont/headscale/hscontrol/db"
	"github.com/juanfont/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSessionSecurity tests session security features.
func TestSessionSecurity(t *testing.T) {
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	cfg := &types.Config{}
	auth := NewHeadplaneAuth(cfg, hsdb)

	user, err := db.CreateHeadplaneUser(hsdb.DB, "secureuser", "securepass", "user")
	require.NoError(t, err)

	t.Run("session tokens are unique", func(t *testing.T) {
		tokens := make(map[string]bool)
		for i := 0; i < 10; i++ {
			token := loginAndGetToken(t, auth, hsdb, "secureuser", "securepass", "user")
			assert.False(t, tokens[token], "duplicate token generated")
			tokens[token] = true
		}
	})

	t.Run("session expires", func(t *testing.T) {
		token := loginAndGetToken(t, auth, hsdb, "secureuser", "securepass", "user")
		assert.True(t, auth.VerifySession(token))

		// Manually expire the session
		auth.mu.Lock()
		session := auth.sessions[token]
		session.ExpiresAt = time.Now().Add(-1 * time.Hour)
		auth.mu.Unlock()

		auth.cleanup()
		assert.False(t, auth.VerifySession(token))
	})

	t.Run("invalid tokens are rejected", func(t *testing.T) {
		assert.False(t, auth.VerifySession("invalid-token"))
		assert.False(t, auth.VerifySession(""))
	})

	t.Run("session contains correct user info", func(t *testing.T) {
		token := loginAndGetToken(t, auth, hsdb, "secureuser", "securepass", "user")

		session, ok := auth.GetSession(token)
		require.True(t, ok)
		assert.Equal(t, user.ID, session.UserID)
		assert.Equal(t, "secureuser", session.Username)
		assert.False(t, session.IsAdmin)
	})
}

// TestPasswordSecurity tests password hashing and validation.
func TestPasswordSecurity(t *testing.T) {
	hsdb := setupTestDB(t)
	defer hsdb.Close()

	t.Run("passwords are hashed", func(t *testing.T) {
		password := "my-secure-password-123"
		user, err := db.CreateHeadplaneUser(hsdb.DB, "hashtest", password, "user")
		require.NoError(t, err)

		assert.NotEqual(t, password, user.PasswordHash)
		assert.NotEmpty(t, user.PasswordHash)
		assert.Regexp(t, `^\$2[aby]\$\d+\$`, user.PasswordHash)
	})

	t.Run("bcrypt produces unique salts", func(t *testing.T) {
		password := "same-password"
		user1, _ := db.CreateHeadplaneUser(hsdb.DB, "hashtest1", password, "user")
		user2, _ := db.CreateHeadplaneUser(hsdb.DB, "hashtest2", password, "user")

		assert.NotEqual(t, user1.PasswordHash, user2.PasswordHash)
	})

	t.Run("password verification works correctly", func(t *testing.T) {
		user, _ := db.CreateHeadplaneUser(hsdb.DB, "checktest", "correct-password", "user")

		assert.True(t, user.CheckPassword("correct-password"))
		assert.False(t, user.CheckPassword("wrong-password"))
		assert.False(t, user.CheckPassword(""))
	})
}

// TestAPIKeyEncryption tests API key encryption and decryption.
func TestAPIKeyEncryption(t *testing.T) {
	apiKey := "test-api-key-1234567890"
	sessionToken := "session-token-abcdef"

	t.Run("encryption produces ciphertext", func(t *testing.T) {
		encrypted, nonce, salt, err := encryptAPIKey(apiKey, sessionToken)
		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)
		assert.NotEmpty(t, nonce)
		assert.NotEmpty(t, salt)
		assert.NotEqual(t, apiKey, encrypted)
	})

	t.Run("decryption recovers original key", func(t *testing.T) {
		encrypted, nonce, salt, err := encryptAPIKey(apiKey, sessionToken)
		require.NoError(t, err)

		decrypted, err := decryptAPIKey(encrypted, nonce, salt, sessionToken)
		require.NoError(t, err)
		assert.Equal(t, apiKey, decrypted)
	})

	t.Run("wrong session token fails", func(t *testing.T) {
		encrypted, nonce, salt, err := encryptAPIKey(apiKey, sessionToken)
		require.NoError(t, err)

		_, err = decryptAPIKey(encrypted, nonce, salt, "wrong-token")
		assert.Error(t, err)
	})

	t.Run("encryption produces unique ciphertext", func(t *testing.T) {
		encrypted1, _, _, _ := encryptAPIKey(apiKey, sessionToken)
		encrypted2, _, _, _ := encryptAPIKey(apiKey, sessionToken)

		assert.NotEqual(t, encrypted1, encrypted2)
	})
}
