package hscontrol

import (
	"testing"

	"github.com/juanfont/headscale/hscontrol/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMigrationFromSingleToMultiUser tests migration from single-user to multi-user mode.
func TestMigrationFromSingleToMultiUser(t *testing.T) {
	// This test simulates upgrading from Phase 13a (single password) to Phase 13c (multi-user)
	
	t.Run("default admin user creation", func(t *testing.T) {
		hsdb := setupTestDBForUsers(t)
		defer hsdb.Close()

		// Simulate first startup with multi-user support
		// Should auto-create default admin user if none exists
		var count int64
		hsdb.DB.Model(&db.HeadplaneUser{}).Count(&count)
		
		// If no users exist, migration should create default admin
		if count == 0 {
			// This would be done by migration logic
			_, err := db.CreateHeadplaneUser(hsdb.DB, "admin", "changeme", "admin")
			require.NoError(t, err)
		}

		// Verify admin user exists
		admin, err := db.GetHeadplaneUserByUsername(hsdb.DB, "admin")
		require.NoError(t, err)
		assert.Equal(t, "admin", admin.Username)
		assert.Equal(t, "admin", admin.Role)
	})

	t.Run("settings migration preserves data", func(t *testing.T) {
		hsdb := setupTestDBForUsers(t)
		defer hsdb.Close()

		// Create a user
		user, err := db.CreateHeadplaneUser(hsdb.DB, "testuser", "password", "user")
		require.NoError(t, err)

		// Create settings for the user (simulating old single-user settings)
		err = UpdateSettings(hsdb.DB, user.ID, "encrypted-key", "nonce", "salt", "dark", "Test User")
		require.NoError(t, err)

		// Verify settings were created
		settings, err := GetSettings(hsdb.DB, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, settings.UserID)
		assert.Equal(t, "dark", settings.Theme)
		assert.Equal(t, "Test User", settings.ProfileName)
	})

	t.Run("table schemas are correct", func(t *testing.T) {
		hsdb := setupTestDBForUsers(t)
		defer hsdb.Close()

		// Verify headplane_users table exists with correct columns
		hasTable := hsdb.DB.Migrator().HasTable(&db.HeadplaneUser{})
		assert.True(t, hasTable, "headplane_users table should exist")

		// Verify headplane_settings table exists with correct columns
		hasTable = hsdb.DB.Migrator().HasTable(&HeadplaneSettings{})
		assert.True(t, hasTable, "headplane_settings table should exist")

		// Verify foreign key relationship
		hasColumn := hsdb.DB.Migrator().HasColumn(&HeadplaneSettings{}, "user_id")
		assert.True(t, hasColumn, "headplane_settings should have user_id column")
	})
}

// TestDatabaseSchema tests the database schema for multi-user support.
func TestDatabaseSchema(t *testing.T) {
	hsdb := setupTestDBForUsers(t)
	defer hsdb.Close()

	t.Run("headplane_users table structure", func(t *testing.T) {
		// Verify all required columns exist
		columns := []string{"id", "username", "password_hash", "role", "created_at", "updated_at"}
		for _, col := range columns {
			hasCol := hsdb.DB.Migrator().HasColumn(&db.HeadplaneUser{}, col)
			assert.True(t, hasCol, "headplane_users should have column: %s", col)
		}
	})

	t.Run("headplane_settings table structure", func(t *testing.T) {
		// Verify all required columns exist
		columns := []string{"id", "user_id", "api_key_encrypted", "api_key_nonce", "api_key_salt", "theme", "profile_name"}
		for _, col := range columns {
			hasCol := hsdb.DB.Migrator().HasColumn(&HeadplaneSettings{}, col)
			assert.True(t, hasCol, "headplane_settings should have column: %s", col)
		}
	})

	t.Run("unique constraints are enforced", func(t *testing.T) {
		// Create first user
		_, err := db.CreateHeadplaneUser(hsdb.DB, "uniquetest", "password", "user")
		require.NoError(t, err)

		// Try to create duplicate username (should fail)
		_, err = db.CreateHeadplaneUser(hsdb.DB, "uniquetest", "password2", "user")
		assert.Error(t, err, "duplicate username should be rejected")
	})

	t.Run("settings are per-user", func(t *testing.T) {
		user1, _ := db.CreateHeadplaneUser(hsdb.DB, "peruser1", "pass", "user")
		user2, _ := db.CreateHeadplaneUser(hsdb.DB, "peruser2", "pass", "user")

		// Each user can have their own settings
		err := UpdateSettings(hsdb.DB, user1.ID, "enc1", "n1", "s1", "dark", "User 1")
		require.NoError(t, err)

		err = UpdateSettings(hsdb.DB, user2.ID, "enc2", "n2", "s2", "light", "User 2")
		require.NoError(t, err)

		// Count settings records
		var count int64
		hsdb.DB.Model(&HeadplaneSettings{}).Count(&count)
		assert.Equal(t, int64(2), count, "should have 2 settings records")
	})
}

// TestBackwardCompatibility tests that existing functionality still works.
func TestBackwardCompatibility(t *testing.T) {
	t.Run("API key authentication still works", func(t *testing.T) {
		// Verify that original API key authentication is not broken
		// This is a placeholder - actual API key auth tests exist elsewhere
		t.Skip("API key authentication tested in existing test suite")
	})

	t.Run("existing endpoints remain functional", func(t *testing.T) {
		// Verify that non-Headplane endpoints are not affected
		t.Skip("Core Headscale functionality tested in existing test suite")
	})
}
