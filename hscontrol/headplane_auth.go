package hscontrol

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/juanfont/headscale/hscontrol/db"
	"github.com/juanfont/headscale/hscontrol/types"
	"github.com/rs/zerolog/log"
)

const (
	headplaneSessionCookieName = "headscale_headplane_session"
	sessionTokenLength         = 32 // 256 bits
	sessionDuration            = 24 * time.Hour
	maxLoginAttempts           = 5
	loginAttemptWindow         = 1 * time.Minute
)

// headplaneSession represents an authenticated Headplane session.
type headplaneSession struct {
	Token     string
	UserID    uint
	Username  string
	IsAdmin   bool
	CreatedAt time.Time
	ExpiresAt time.Time
}

// loginAttempt tracks failed login attempts for rate limiting.
type loginAttempt struct {
	Count        int
	FirstAttempt time.Time
}

// HeadplaneAuth handles password-based authentication for Headplane.
type HeadplaneAuth struct {
	cfg           *types.Config
	db            *db.HSDatabase
	sessions      map[string]*headplaneSession // token -> session
	loginAttempts map[string]*loginAttempt     // IP -> attempt record
	mu            sync.RWMutex
}

// NewHeadplaneAuth creates a new Headplane authentication handler.
func NewHeadplaneAuth(cfg *types.Config, hsdb *db.HSDatabase) *HeadplaneAuth {
	auth := &HeadplaneAuth{
		cfg:           cfg,
		db:            hsdb,
		sessions:      make(map[string]*headplaneSession),
		loginAttempts: make(map[string]*loginAttempt),
	}

	go auth.cleanupLoop()

	return auth
}

func (h *HeadplaneAuth) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		h.cleanup()
	}
}

func (h *HeadplaneAuth) cleanup() {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now()
	for token, session := range h.sessions {
		if now.After(session.ExpiresAt) {
			delete(h.sessions, token)
		}
	}
	for ip, attempt := range h.loginAttempts {
		if now.Sub(attempt.FirstAttempt) > loginAttemptWindow {
			delete(h.loginAttempts, ip)
		}
	}
}

func (h *HeadplaneAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	clientIP := r.RemoteAddr
	if h.isRateLimited(clientIP) {
		http.Error(w, "Too many attempts", http.StatusTooManyRequests)
		return
	}

	// Get user from database
	user, err := db.GetHeadplaneUserByUsername(h.db.DB, req.Username)
	if err != nil {
		h.recordFailedAttempt(clientIP)
		log.Warn().Str("ip", clientIP).Str("username", req.Username).Msg("Headplane login failed: user not found")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		h.recordFailedAttempt(clientIP)
		log.Warn().Str("ip", clientIP).Str("username", req.Username).Msg("Headplane login failed: invalid password")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	h.clearFailedAttempts(clientIP)
	token, err := h.generateToken()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	session := &headplaneSession{
		Token:     token,
		UserID:    user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	h.mu.Lock()
	h.sessions[token] = session
	h.mu.Unlock()

	log.Info().Str("ip", clientIP).Str("username", user.Username).Msg("Headplane login success")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      token,
		"expires_at": session.ExpiresAt.Unix(),
		"username":   user.Username,
		"is_admin":   user.IsAdmin(),
	})
}

// VerifySession checks if a session token is valid.
func (h *HeadplaneAuth) VerifySession(token string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	session, exists := h.sessions[token]
	if !exists || time.Now().After(session.ExpiresAt) {
		return false
	}
	return true
}

// GetSession retrieves the session for a token.
func (h *HeadplaneAuth) GetSession(token string) (*headplaneSession, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	session, exists := h.sessions[token]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil, false
	}
	return session, true
}

func (h *HeadplaneAuth) isRateLimited(ip string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	attempt, exists := h.loginAttempts[ip]
	return exists && attempt.Count >= maxLoginAttempts && time.Since(attempt.FirstAttempt) <= loginAttemptWindow
}

func (h *HeadplaneAuth) recordFailedAttempt(ip string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	attempt, exists := h.loginAttempts[ip]
	if !exists || time.Since(attempt.FirstAttempt) > loginAttemptWindow {
		h.loginAttempts[ip] = &loginAttempt{Count: 1, FirstAttempt: time.Now()}
	} else {
		attempt.Count++
	}
}

func (h *HeadplaneAuth) clearFailedAttempts(ip string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.loginAttempts, ip)
}

func (h *HeadplaneAuth) generateToken() (string, error) {
	b := make([]byte, sessionTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
