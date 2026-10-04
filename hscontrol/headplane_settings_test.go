package hscontrol

import (
	"testing"
	"time"

	"github.com/juanfont/headscale/hscontrol/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForSettings(t *testing.T) *gorm.DB {
	t.Helper()
	dbConn, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	
	// Create headplane_users table (required for foreign key)
	require.NoError(t, dbConn.AutoMigrate(&db.HeadplaneUser{}))
	
	// Create headplane_settings table
	require.NoError(t, InitHeadplaneSettings(dbConn))
	
	return dbConn
}

func createTestUser(t *testing.T, dbConn *gorm.DB, username, password, role string) *db.HeadplaneUser {
	t.Helper()
	user, err := db.CreateHeadplaneUser(dbConn, username, password, role)
	require.NoError(t, err)
	return user
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

func TestGetSettings_AutoCreate(t *testing.T) {
	dbConn := setupTestDBForSettings(t)
	user := createTestUser(t, dbConn, "testuser", "password", "user")

	// Get settings should auto-create
	settings, err := GetSettings(dbConn, user.ID)
	require.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, user.ID, settings.UserID)
	assert.Equal(t, "light", settings.Theme)
}

func TestPerUserIsolation(t *testing.T) {
	dbConn := setupTestDBForSettings(t)
	user1 := createTestUser(t, dbConn, "user1", "password1", "user")
	user2 := createTestUser(t, dbConn, "user2", "password2", "user")

	// Create settings for user1
	err := UpdateSettings(dbConn, user1.ID, "enc1", "nonce1", "salt1", "dark", "User One")
	require.NoError(t, err)

	// Create settings for user2
	err = UpdateSettings(dbConn, user2.ID, "enc2", "nonce2", "salt2", "light", "User Two")
	require.NoError(t, err)

	// Verify user1's settings
	settings1, err := GetSettings(dbConn, user1.ID)
	require.NoError(t, err)
	assert.Equal(t, user1.ID, settings1.UserID)
	assert.Equal(t, "dark", settings1.Theme)

	// Verify user2's settings
	settings2, err := GetSettings(dbConn, user2.ID)
	require.NoError(t, err)
	assert.Equal(t, user2.ID, settings2.UserID)
	assert.Equal(t, "light", settings2.Theme)

	// Verify both users have separate settings
	var count int64
	dbConn.Model(&HeadplaneSettings{}).Count(&count)
	assert.Equal(t, int64(2), count)
}

// Mock state for testing
type mockStateForSettings struct {
	db *gorm.DB
}

func (m *mockStateForSettings) DB() *db.HSDatabase {
	return &db.HSDatabase{DB: m.db}
}

func testNow() time.Time {
	return time.Date(2026, 10, 3, 4, 32, 0, 0, time.UTC)
}
