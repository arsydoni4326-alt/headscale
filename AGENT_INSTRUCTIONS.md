# Agent 1: Backend Settings Implementation

**Task:** Phase 13b Part 1 & 2 — Backend Settings Storage & Endpoints  
**Worktree:** `/home/denny/Project/headscale-settings-backend`  
**Branch:** `feature/settings-backend`  
**Base Commit:** `d95c5eff`

---

## Objective

Implement backend infrastructure for Headplane settings:
- SQLite table for storing user settings (encrypted API key, theme, profile name)
- Settings CRUD endpoints (GET/POST)
- Password change endpoint
- Session enhancement to auto-load stored API key
- Comprehensive unit tests

---

## Context

Phase 13a (simple password auth) is complete. Users can log in with a password, but:
- Must re-enter API key at each login (no storage)
- Cannot change password without editing config
- No theme or profile customization

This task fixes that by adding persistent settings storage.

---

## Architecture Overview

### Settings Storage

**SQLite table** in Headplane's existing database (already used for audit log):

```sql
CREATE TABLE headplane_settings (
  id INTEGER PRIMARY KEY CHECK (id = 1),  -- Single row constraint
  api_key_encrypted TEXT,                 -- AES-256-GCM encrypted
  api_key_nonce TEXT,                     -- Encryption nonce
  theme TEXT DEFAULT 'light',             -- 'light' or 'dark'
  profile_name TEXT,                      -- Optional display name
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Encryption

**API key encryption:** AES-256-GCM with PBKDF2-derived key from session token.

**Rationale:**
- API key is sensitive (grants full Headscale access)
- Session token is already secure (256-bit random)
- PBKDF2 adds key stretching
- AES-GCM provides authenticated encryption

### API Endpoints

1. **GET /api/v1/headplane/settings**
   - Requires: Valid password session token (from Phase 13a)
   - Returns: `{apiKey: string, theme: string, profileName: string}`
   - Decrypts API key before returning

2. **POST /api/v1/headplane/settings**
   - Requires: Valid password session token
   - Body: `{apiKey?: string, theme?: string, profileName?: string}`
   - Encrypts API key before storing
   - Returns: `{success: true}`

3. **POST /api/v1/headplane/change-password**
   - Requires: Valid password session token
   - Body: `{currentPassword: string, newPassword: string}`
   - Validates current password (constant-time comparison)
   - Updates config password in-memory
   - Returns: `{success: true}`

---

## Implementation Tasks

### Task 1: Settings Table Schema

**File:** Create `hscontrol/headplane_settings.go`

1. Define `HeadplaneSettings` struct:
   ```go
   type HeadplaneSettings struct {
       ID            int       `gorm:"primaryKey;check:id = 1"`
       APIKeyEncrypted string `gorm:"column:api_key_encrypted"`
       APIKeyNonce     string `gorm:"column:api_key_nonce"`
       Theme          string `gorm:"default:light"`
       ProfileName    string `gorm:"column:profile_name"`
       UpdatedAt      time.Time
   }
   ```

2. Implement table creation in Headplane DB init
   - Check where Headplane audit table is created
   - Add settings table creation there

### Task 2: Encryption/Decryption Functions

**File:** `hscontrol/headplane_settings.go`

Implement:
- `encryptAPIKey(apiKey string, sessionToken string) (encrypted string, nonce string, error)`
- `decryptAPIKey(encrypted string, nonce string, sessionToken string) (string, error)`

Use:
- `crypto/aes` for AES-256-GCM
- `crypto/cipher` for GCM mode
- `golang.org/x/crypto/pbkdf2` for key derivation
- Salt: fixed per-installation (store in config or generate once)

### Task 3: Settings CRUD

**File:** `hscontrol/headplane_settings.go`

Implement:
- `GetSettings(sessionToken string) (*HeadplaneSettings, error)`
  - Query single row from table
  - Decrypt API key
  - Return struct

- `UpdateSettings(sessionToken string, apiKey *string, theme *string, profileName *string) error`
  - Fetch existing settings (or create if not exists)
  - Update fields (only non-nil pointers)
  - Encrypt API key if provided
  - Save to database

### Task 4: GET /api/v1/headplane/settings Endpoint

**File:** `hscontrol/headplane_settings.go`

Implement handler:
```go
func (b *Headscale) handleHeadplaneGetSettings(w http.ResponseWriter, r *http.Request) {
    // Extract session token from auth header
    // Call GetSettings(token)
    // Return JSON: {apiKey, theme, profileName}
}
```

Register in `hscontrol/app.go`:
```go
r.Get("/api/v1/headplane/settings", b.handleHeadplaneGetSettings)
```

### Task 5: POST /api/v1/headplane/settings Endpoint

**File:** `hscontrol/headplane_settings.go`

Implement handler:
```go
func (b *Headscale) handleHeadplaneUpdateSettings(w http.ResponseWriter, r *http.Request) {
    // Parse JSON body
    // Extract session token
    // Call UpdateSettings(token, ...)
    // Return {success: true}
}
```

Register in `hscontrol/app.go`:
```go
r.Post("/api/v1/headplane/settings", b.handleHeadplaneUpdateSettings)
```

### Task 6: POST /api/v1/headplane/change-password Endpoint

**File:** `hscontrol/headplane_auth.go` (modify existing file)

Implement handler:
```go
func (b *Headscale) handleHeadplaneChangePassword(w http.ResponseWriter, r *http.Request) {
    // Parse {currentPassword, newPassword}
    // Validate current password (constant-time)
    // Update b.Cfg.Headplane.Password in-memory
    // Optional: persist to config file (if writable)
    // Return {success: true}
}
```

Register in `hscontrol/app.go`:
```go
r.Post("/api/v1/headplane/change-password", b.handleHeadplaneChangePassword)
```

### Task 7: Session Enhancement

**File:** `hscontrol/headplane_auth.go`

Modify session validation to:
1. Load stored API key from settings (if exists)
2. Inject into session context
3. Frontend API calls use stored key automatically

**Implementation:**
- In `validatePasswordSession()` or similar:
  ```go
  settings, _ := GetSettings(token)
  if settings != nil && settings.APIKeyEncrypted != "" {
      // Decrypt and inject into context
  }
  ```

### Task 8: Unit Tests

**File:** Create `hscontrol/headplane_settings_test.go`

Test cases:
1. `TestEncryptDecryptAPIKey` — encryption roundtrip
2. `TestGetSettings_NotExists` — empty table returns nil
3. `TestUpdateSettings_Create` — first update creates row
4. `TestUpdateSettings_Update` — subsequent updates modify row
5. `TestGetSettingsEndpoint_Success` — HTTP GET returns settings
6. `TestUpdateSettingsEndpoint_Success` — HTTP POST updates settings
7. `TestChangePassword_ValidCurrent` — password change succeeds
8. `TestChangePassword_InvalidCurrent` — wrong current password rejected
9. `TestSettingsRequireAuth` — unauthenticated requests fail

Use existing test patterns from `hscontrol/headplane_auth_test.go`.

---

## Dependencies

No external dependencies. Use existing libraries:
- `crypto/aes`, `crypto/cipher`
- `golang.org/x/crypto/pbkdf2`
- `gorm.io/gorm`

---

## Files to Create/Modify

### Create:
- `hscontrol/headplane_settings.go` (~250 lines)
- `hscontrol/headplane_settings_test.go` (~200 lines)

### Modify:
- `hscontrol/app.go` (register 3 new endpoints)
- `hscontrol/headplane_auth.go` (password change, session enhancement)
- `hscontrol/headplane_auth_test.go` (update for session enhancement)

---

## Testing

```bash
cd /home/denny/Project/headscale-settings-backend

# Run backend tests
go test ./hscontrol -v -run TestHeadplane

# Run all tests
go test ./...

# Check coverage
go test -cover ./hscontrol
```

---

## Acceptance Criteria

- [ ] `headplane_settings` SQLite table created
- [ ] API key encrypted with AES-256-GCM
- [ ] `GET /api/v1/headplane/settings` returns decrypted settings
- [ ] `POST /api/v1/headplane/settings` updates and encrypts
- [ ] `POST /api/v1/headplane/change-password` validates and changes password
- [ ] Session loads stored API key automatically
- [ ] All unit tests pass
- [ ] No regressions in Phase 13a functionality

---

## Coordination with Other Agents

### With Frontend Agent:
- **API Contract:** Share OpenAPI spec or endpoint signatures
- **Mock Endpoints:** Frontend can mock these endpoints initially
- **Integration Testing:** After both complete, test end-to-end

### With Docs Agent:
- **API Documentation:** Provide endpoint details for docs
- **Security Notes:** Explain encryption approach

---

## Completion

When done:
1. Run all tests and verify they pass
2. Commit your changes
3. Push branch: `git push -u origin feature/settings-backend`
4. Create summary: list files changed, tests added, commit SHA
5. Report completion with merge readiness

**Do not merge yet.** Integration testing happens after all agents complete.

---

## References

- **Implementation Plan:** `docs/phase13b-implementation-plan.md`
- **Phase 13a Auth:** `hscontrol/headplane_auth.go`
- **Existing Tests:** `hscontrol/headplane_auth_test.go`
- **ROADMAP:** `ROADMAP.md` Phase 13b
