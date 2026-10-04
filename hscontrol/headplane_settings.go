package hscontrol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/arsydoni4326-alt/headscale/hscontrol/db"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
)

const (
	pbkdf2Iterations = 100000
	pbkdf2KeyLen     = 32 // AES-256
	pbkdf2SaltLen    = 32
	aesGCMNonceSize  = 12
)

var (
	errSettingsNotFound       = errors.New("settings not found")
	errInvalidEncryptedData   = errors.New("invalid encrypted data")
	errDecryptionFailed       = errors.New("decryption failed")
	errSessionTokenRequired   = errors.New("session token required for encryption")
	errInvalidCurrentPassword = errors.New("invalid current password")
)

// HeadplaneSettings stores per-user settings for Headplane UI.
type HeadplaneSettings struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"uniqueIndex;not null"`
	APIKeyEncrypted string    `gorm:"column:api_key_encrypted"`
	APIKeyNonce     string    `gorm:"column:api_key_nonce"`
	APIKeySalt      string    `gorm:"column:api_key_salt"` // PBKDF2 salt
	Theme           string    `gorm:"default:light"`
	ProfileName     string    `gorm:"column:profile_name"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM.
func (HeadplaneSettings) TableName() string {
	return "headplane_settings"
}

// SettingsResponse is the JSON response for GET /api/v1/headplane/settings.
type SettingsResponse struct {
	APIKey      string `json:"apiKey"`
	Theme       string `json:"theme"`
	ProfileName string `json:"profileName"`
}

// UpdateSettingsRequest is the JSON request for POST /api/v1/headplane/settings.
type UpdateSettingsRequest struct {
	APIKey      *string `json:"apiKey,omitempty"`
	Theme       *string `json:"theme,omitempty"`
	ProfileName *string `json:"profileName,omitempty"`
}

// ChangePasswordRequest is the JSON request for POST /api/v1/headplane/change-password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// encryptAPIKey encrypts an API key using AES-256-GCM with a key derived from the session token.
func encryptAPIKey(apiKey, sessionToken string) (encrypted, nonce, salt string, err error) {
	if sessionToken == "" {
		return "", "", "", errSessionTokenRequired
	}

	// Generate random salt for PBKDF2
	saltBytes := make([]byte, pbkdf2SaltLen)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", "", "", err
	}

	// Derive encryption key from session token using PBKDF2
	key := pbkdf2.Key([]byte(sessionToken), saltBytes, pbkdf2Iterations, pbkdf2KeyLen, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", "", err
	}

	// Generate nonce
	nonceBytes := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonceBytes); err != nil {
		return "", "", "", err
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonceBytes, []byte(apiKey), nil)

	return base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(nonceBytes),
		base64.StdEncoding.EncodeToString(saltBytes),
		nil
}

// decryptAPIKey decrypts an API key using AES-256-GCM.
func decryptAPIKey(encrypted, nonce, salt, sessionToken string) (string, error) {
	if sessionToken == "" {
		return "", errSessionTokenRequired
	}

	// Decode base64 inputs
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", errInvalidEncryptedData
	}

	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return "", errInvalidEncryptedData
	}

	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", errInvalidEncryptedData
	}

	// Derive same key from session token
	key := pbkdf2.Key([]byte(sessionToken), saltBytes, pbkdf2Iterations, pbkdf2KeyLen, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonceBytes, ciphertext, nil)
	if err != nil {
		return "", errDecryptionFailed
	}

	return string(plaintext), nil
}

// GetSettings retrieves the settings for a specific user from the database.
// If no settings exist for the user, they are automatically created with defaults.
func GetSettings(db *gorm.DB, userID uint) (*HeadplaneSettings, error) {
	var settings HeadplaneSettings
	err := db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Auto-create settings for the user with defaults
			settings = HeadplaneSettings{
				UserID:      userID,
				Theme:       "light",
				ProfileName: "",
			}
			if err := db.Create(&settings).Error; err != nil {
				return nil, err
			}
			return &settings, nil
		}
		return nil, err
	}
	return &settings, nil
}

// UpdateSettings updates the settings for a specific user in the database.
// If no settings exist for the user, they are created.
func UpdateSettings(db *gorm.DB, userID uint, apiKeyEncrypted, apiKeyNonce, apiKeySalt, theme, profileName string) error {
	// First, check if settings exist
	var existing HeadplaneSettings
	err := db.Where("user_id = ?", userID).First(&existing).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new settings
		settings := HeadplaneSettings{
			UserID:          userID,
			APIKeyEncrypted: apiKeyEncrypted,
			APIKeyNonce:     apiKeyNonce,
			APIKeySalt:      apiKeySalt,
			Theme:           theme,
			ProfileName:     profileName,
		}
		return db.Create(&settings).Error
	}
	
	if err != nil {
		return err
	}
	
	// Update existing settings
	return db.Model(&existing).Updates(map[string]interface{}{
		"api_key_encrypted": apiKeyEncrypted,
		"api_key_nonce":     apiKeyNonce,
		"api_key_salt":      apiKeySalt,
		"theme":             theme,
		"profile_name":      profileName,
	}).Error
}

// InitHeadplaneSettings creates the headplane_settings table if it doesn't exist.
func InitHeadplaneSettings(db *gorm.DB) error {
	return db.AutoMigrate(&HeadplaneSettings{})
}

// HandleGetSettings handles GET /api/v1/headplane/settings.
func (h *Headscale) HandleGetSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify session and get user ID
	token := r.Header.Get("Authorization")
	if token == "" {
		// Try cookie
		cookie, err := r.Cookie(headplaneSessionCookieName)
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" || !h.headplaneAuth.VerifySession(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get session to extract user ID
	session, ok := h.headplaneAuth.GetSession(token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get settings from database for this user (auto-creates if not found)
	settings, err := GetSettings(h.state.DB().DB, session.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get settings")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Decrypt API key if present
	var apiKey string
	if settings.APIKeyEncrypted != "" && settings.APIKeyNonce != "" && settings.APIKeySalt != "" {
		decrypted, err := decryptAPIKey(settings.APIKeyEncrypted, settings.APIKeyNonce, settings.APIKeySalt, token)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decrypt API key")
			// Don't fail the request, just return empty API key
			apiKey = ""
		} else {
			apiKey = decrypted
		}
	}

	// Return settings
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SettingsResponse{
		APIKey:      apiKey,
		Theme:       settings.Theme,
		ProfileName: settings.ProfileName,
	})
}

// HandleUpdateSettings handles POST /api/v1/headplane/settings.
func (h *Headscale) HandleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify session and get user ID
	token := r.Header.Get("Authorization")
	if token == "" {
		// Try cookie
		cookie, err := r.Cookie(headplaneSessionCookieName)
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" || !h.headplaneAuth.VerifySession(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get session to extract user ID
	session, ok := h.headplaneAuth.GetSession(token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request
	var req UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get current settings (auto-creates with defaults if not found)
	currentSettings, err := GetSettings(h.state.DB().DB, session.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current settings")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Start with current values
	apiKeyEncrypted := currentSettings.APIKeyEncrypted
	apiKeyNonce := currentSettings.APIKeyNonce
	apiKeySalt := currentSettings.APIKeySalt
	theme := currentSettings.Theme
	profileName := currentSettings.ProfileName

	// Update API key if provided
	if req.APIKey != nil {
		if *req.APIKey == "" {
			// Clear API key
			apiKeyEncrypted = ""
			apiKeyNonce = ""
			apiKeySalt = ""
		} else {
			// Encrypt new API key
			encrypted, nonce, salt, err := encryptAPIKey(*req.APIKey, token)
			if err != nil {
				log.Error().Err(err).Msg("Failed to encrypt API key")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			apiKeyEncrypted = encrypted
			apiKeyNonce = nonce
			apiKeySalt = salt
		}
	}

	// Update theme if provided
	if req.Theme != nil {
		theme = *req.Theme
	}

	// Update profile name if provided
	if req.ProfileName != nil {
		profileName = *req.ProfileName
	}

	// Save to database
	err = UpdateSettings(h.state.DB().DB, session.UserID, apiKeyEncrypted, apiKeyNonce, apiKeySalt, theme, profileName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update settings")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info().Uint("user_id", session.UserID).Msg("Settings updated successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// HandleChangePassword handles POST /api/v1/headplane/change-password.
func (h *Headscale) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify session
	token := r.Header.Get("Authorization")
	if token == "" {
		// Try cookie
		cookie, err := r.Cookie(headplaneSessionCookieName)
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" || !h.headplaneAuth.VerifySession(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get the session to identify the user
	session, ok := h.headplaneAuth.GetSession(token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user, err := db.GetHeadplaneUserByID(h.state.DB().DB, session.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user for password change")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Validate current password
	if !user.CheckPassword(req.CurrentPassword) {
		log.Warn().Msg("Password change failed: invalid current password")
		http.Error(w, "Invalid current password", http.StatusUnauthorized)
		return
	}

	// Validate new password
	if req.NewPassword == "" {
		http.Error(w, "New password cannot be empty", http.StatusBadRequest)
		return
	}

	// Update password in database
	err = db.UpdateHeadplaneUserPassword(h.state.DB().DB, session.UserID, req.NewPassword)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update password")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info().Str("username", session.Username).Msg("Password changed successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
