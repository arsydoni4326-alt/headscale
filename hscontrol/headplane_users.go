package hscontrol

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/arsydoni4326-alt/headscale/hscontrol/db"
	"github.com/rs/zerolog/log"
)

// HeadplaneUserResponse is the JSON response for user operations.
type HeadplaneUserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"createdAt"`
}

// RegisterUserRequest is the JSON request for user registration.
type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// authContext represents the authenticated context for a request
type authContext struct {
	IsAdmin  bool
	Username string
	UserID   uint
	IsAPIKey bool // true if authenticated via API key, false if via password session
}

// requireAdminAuth checks if the request has valid admin authentication (either password session or API key).
// Returns the authentication context for use by the endpoint.
func (h *Headscale) requireAdminAuth(r *http.Request) (*authContext, error) {
	token := r.Header.Get("Authorization")
	
	// Strip "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	
	if token == "" {
		// Try cookie (password sessions only)
		cookie, err := r.Cookie(headplaneSessionCookieName)
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" {
		return nil, db.ErrHeadplaneUserNotAuthenticated
	}

	// First, try password session authentication
	if session, ok := h.headplaneAuth.GetSession(token); ok {
		if !session.IsAdmin {
			return nil, db.ErrHeadplaneUserForbidden
		}
		return &authContext{
			IsAdmin:  true,
			Username: session.Username,
			UserID:   session.UserID,
			IsAPIKey: false,
		}, nil
	}

	// Second, try API key authentication
	apiKey, err := h.state.DB().AuthenticateAPIKey(token)
	if err == nil && apiKey != nil {
		// API keys are always admin (all-access)
		return &authContext{
			IsAdmin:  true,
			Username: "api-key-admin", // Placeholder username for API key auth
			UserID:   0,                // No user ID for API keys
			IsAPIKey: true,
		}, nil
	}

	// Neither password session nor API key was valid
	return nil, db.ErrHeadplaneUserNotAuthenticated
}

// requirePasswordAdminSession checks if the request has a valid password-authenticated admin session.
// Some endpoints (like password change) require password authentication specifically, not API key authentication.
func (h *Headscale) requirePasswordAdminSession(r *http.Request) (*headplaneSession, error) {
	token := r.Header.Get("Authorization")
	
	// Strip "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	
	if token == "" {
		// Try cookie
		cookie, err := r.Cookie(headplaneSessionCookieName)
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" {
		return nil, db.ErrHeadplaneUserNotAuthenticated
	}

	// Check if it's a password session
	session, ok := h.headplaneAuth.GetSession(token)
	if !ok {
		// Check if it's an API key trying to access a password-only endpoint
		_, err := h.state.DB().AuthenticateAPIKey(token)
		if err == nil {
			// It's a valid API key, but this endpoint requires password auth
			return nil, db.ErrHeadplanePasswordAuthRequired
		}
		return nil, db.ErrHeadplaneUserNotAuthenticated
	}

	if !session.IsAdmin {
		return nil, db.ErrHeadplaneUserForbidden
	}

	return session, nil
}

// HandleRegisterUser handles POST /api/v1/headplane/users (admin-only).
func (h *Headscale) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify admin authentication (password session or API key)
	authCtx, err := h.requireAdminAuth(r)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	// Parse request
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Default to 'user' role if not specified
	if req.Role == "" {
		req.Role = "user"
	}

	// Only allow 'user' or 'admin' roles
	if req.Role != "user" && req.Role != "admin" {
		http.Error(w, "Invalid role. Must be 'user' or 'admin'", http.StatusBadRequest)
		return
	}

	// Create user
	user, err := db.CreateHeadplaneUser(h.state.DB().DB, req.Username, req.Password, req.Role)
	if err != nil {
		if err == db.ErrHeadplaneUserExists {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("Failed to create headplane user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info().
		Str("admin", authCtx.Username).
		Bool("via_api_key", authCtx.IsAPIKey).
		Str("new_user", user.Username).
		Str("role", user.Role).
		Msg("Headplane user created")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HeadplaneUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Unix(),
	})
}

// HandleListUsers handles GET /api/v1/headplane/users (admin-only).
func (h *Headscale) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify admin authentication (password session or API key)
	_, err := h.requireAdminAuth(r)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	// Get all users
	users, err := db.ListHeadplaneUsers(h.state.DB().DB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list headplane users")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	userList := make([]HeadplaneUserResponse, len(users))
	for i, user := range users {
		userList[i] = HeadplaneUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Unix(),
		}
	}

	// Wrap in object with "users" key to match frontend expectations
	response := map[string]interface{}{
		"users": userList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGetUser handles GET /api/v1/headplane/users/:id (admin-only).
func (h *Headscale) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify admin authentication (password session or API key)
	_, err := h.requireAdminAuth(r)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	// Get user ID from URL
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := db.GetHeadplaneUserByID(h.state.DB().DB, uint(userID))
	if err != nil {
		if err == db.ErrHeadplaneUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("Failed to get headplane user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HeadplaneUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Unix(),
	})
}

// UpdateUserRequest is the JSON request for updating a user.
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}

// HandleUpdateUser handles PUT /api/v1/headplane/users/:id (admin-only).
func (h *Headscale) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify admin authentication (password session or API key)
	authCtx, err := h.requireAdminAuth(r)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	// Get user ID from URL
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate role if provided
	if req.Role != "" && req.Role != "user" && req.Role != "admin" {
		http.Error(w, "Invalid role. Must be 'user' or 'admin'", http.StatusBadRequest)
		return
	}

	// Check if attempting to demote self from admin (only applies to password sessions)
	if !authCtx.IsAPIKey && uint(userID) == authCtx.UserID && req.Role == "user" {
		// Check if this would leave no admins
		adminCount, err := db.CountAdminUsers(h.state.DB().DB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to count admin users")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if adminCount <= 1 {
			http.Error(w, "Cannot demote the last admin user", http.StatusBadRequest)
			return
		}
	}

	// Update user
	user, err := db.UpdateHeadplaneUser(h.state.DB().DB, uint(userID), req.Username, req.Role)
	if err != nil {
		if err == db.ErrHeadplaneUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if err == db.ErrHeadplaneUserExists {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("Failed to update headplane user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info().
		Str("admin", authCtx.Username).
		Bool("via_api_key", authCtx.IsAPIKey).
		Str("updated_user", user.Username).
		Uint("user_id", user.ID).
		Msg("Headplane user updated")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HeadplaneUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Unix(),
	})
}

// HandleDeleteUser handles DELETE /api/v1/headplane/users/:id (admin-only).
func (h *Headscale) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify admin authentication (password session or API key)
	authCtx, err := h.requireAdminAuth(r)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	// Get user ID from URL
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Prevent self-deletion (only applies to password sessions)
	if !authCtx.IsAPIKey && uint(userID) == authCtx.UserID {
		http.Error(w, "Cannot delete your own account", http.StatusBadRequest)
		return
	}

	// Get user before deletion for logging and admin check
	user, err := db.GetHeadplaneUserByID(h.state.DB().DB, uint(userID))
	if err != nil {
		if err == db.ErrHeadplaneUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("Failed to get headplane user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prevent deleting the last admin
	if user.IsAdmin() {
		adminCount, err := db.CountAdminUsers(h.state.DB().DB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to count admin users")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if adminCount <= 1 {
			http.Error(w, "Cannot delete the last admin user", http.StatusBadRequest)
			return
		}
	}

	// Delete user
	err = db.DeleteHeadplaneUser(h.state.DB().DB, uint(userID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete headplane user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info().
		Str("admin", authCtx.Username).
		Bool("via_api_key", authCtx.IsAPIKey).
		Str("deleted_user", user.Username).
		Msg("Headplane user deleted")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// handleAuthError writes appropriate HTTP error response based on auth error type.
func (h *Headscale) handleAuthError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	
	switch err {
	case db.ErrHeadplaneUserNotAuthenticated:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "unauthorized",
			"message": "Authentication required. Please log in.",
		})
	case db.ErrHeadplaneUserForbidden:
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "forbidden",
			"message": "Admin privileges required to access this resource.",
		})
	case db.ErrHeadplanePasswordAuthRequired:
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "forbidden",
			"message": "This endpoint requires password authentication. API key authentication is not allowed.",
		})
	default:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "unauthorized",
			"message": "Authentication failed.",
		})
	}
}
