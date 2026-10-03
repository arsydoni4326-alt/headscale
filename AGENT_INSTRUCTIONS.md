# Phase 13b Agent Instructions: Settings Menu Implementation

## Overview

This document coordinated the parallel implementation of Phase 13b (Settings Menu). All agents have completed their work and the branches are ready for integration.

---

## Completion Status

### ✅ Agent 1: Backend Settings Implementation (COMPLETE)

**Branch:** `feature/settings-backend`  
**Completion Document:** `BACKEND_SETTINGS_COMPLETE.md`, `IMPLEMENTATION_COMPLETE.md`

**Deliverables:**
- SQLite `headplane_settings` table with encrypted API key storage
- GET /api/v1/headplane/settings endpoint
- POST /api/v1/headplane/settings endpoint (partial updates)
- POST /api/v1/headplane/change-password endpoint
- AES-256-GCM encryption for stored API key
- Session-derived encryption keys
- Comprehensive unit tests

**Status:** Merged into main (commit 6084bf54)

---

### ⏸️ Agent 2: Frontend Settings UI (BLOCKED)

**Branch:** `feature/settings-frontend`  
**Status:** Architecture violation - cannot merge

**Issue:** The branch converted the headplane Git submodule into a regular directory with all files committed directly to the headscale repository. This violates the documented architecture where Headplane is maintained as a separate submodule.

**Required Action:** Frontend work must be committed to the headplane repository (github.com/arsydoni4326-alt/headplane), and the headscale repository should only update the submodule pointer.

---

### ✅ Agent 3: Documentation & Integration Testing (COMPLETE)

**Branch:** `feature/settings-docs`  
**Completion Document:** `PHASE1_COMPLETE.md`

**Deliverables:**
- User guide: `docs/usage/settings.md`
- API reference: `docs/ref/api/headplane-settings.md`
- Integration test plan: `docs/phase13b-integration-testing.md`
- Updated authentication guide with settings references
- Updated README and CHANGELOG

**Status:** Currently being merged

---

## Phase 13b Implementation Summary

**Phase 13a (Password Auth):** ✅ Complete and merged in v0.35.9-arsydoni4326-alt

**Phase 13b (Settings Menu):**
- Backend: ✅ Complete and merged
- Frontend: ❌ Blocked (architecture violation)
- Documentation: ✅ Complete and merging

---

## Next Steps

1. **Resolve frontend architecture issue:**
   - Move frontend changes to the headplane repository
   - Update headscale's headplane submodule pointer
   - Create new feature branch that properly updates the submodule

2. **Complete Phase 13b integration:**
   - Backend (already merged)
   - Documentation (being merged now)
   - Frontend (pending architecture fix)

3. **Integration testing:**
   - After all merges complete
   - Follow test plan in `docs/phase13b-integration-testing.md`

---

## Architecture Overview (For Reference)

### Settings Storage

**SQLite table** in Headplane's existing database:

```sql
CREATE TABLE headplane_settings (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  api_key_encrypted TEXT,
  api_key_nonce TEXT,
  theme TEXT DEFAULT 'light',
  profile_name TEXT,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Encryption

**API key encryption:** AES-256-GCM with PBKDF2-derived key from session token.

### API Endpoints

1. **GET /api/v1/headplane/settings** - Retrieve and decrypt settings
2. **POST /api/v1/headplane/settings** - Update settings (partial updates)
3. **POST /api/v1/headplane/change-password** - Change password with validation

---

## Implementation Tasks

<<<<<<< HEAD
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
=======
### Task 1: Create Settings Documentation

**File:** Create `docs/usage/settings.md`

**Structure:**
```markdown
# Settings

The Headplane Settings page allows you to customize your experience and manage your Headscale integration.
>>>>>>> feature/settings-docs

## Accessing Settings

After logging in with your password, click "Settings" in the navigation bar.

## Settings Sections

### Account

**Change Password**

[Instructions for changing password]

**Session Information**

[Explain session expiry, logout]

### Integration

**Headscale API Key**

Store your Headscale API key to avoid re-entering it at each login.

1. Obtain API key from Headscale: `headscale apikeys create`
2. Paste into "API Key" field in Settings
3. Click "Save"
4. Your API key is now stored securely (encrypted)

**Security Note:** The API key is encrypted before storage using AES-256-GCM.

### Preferences

**Theme**

Choose between light and dark themes. Your preference is saved automatically.

### Profile

**Display Name**

Set an optional display name for your profile.

## Dependencies

<<<<<<< HEAD
No external dependencies. Use existing libraries:
- `crypto/aes`, `crypto/cipher`
- `golang.org/x/crypto/pbkdf2`
- `gorm.io/gorm`
=======
- API key is encrypted at rest
- Password changes require current password verification
- All settings changes are logged
- Use HTTPS in production

## Troubleshooting

[Common issues and solutions]
```

### Task 2: Update Authentication Guide

**File:** Modify `docs/usage/authentication.md`

Add section:
```markdown
## Settings and API Key Storage

After logging in with your password, you can save your Headscale API key in Settings:

1. Log in to Headplane with your password
2. Navigate to Settings → Integration
3. Enter your Headscale API key
4. Click Save

Your API key will be stored securely and reused across sessions. See [Settings](settings.md) for details.
```

### Task 3: Update README

**File:** `README.md`

Update features section:
```markdown
### Authentication

- **Password Authentication** — Simple password-based login for Headplane web UI access
- **Settings Menu** — Manage API key, password, theme, and profile after login
- **API Keys** — Token-based authentication for programmatic API access and automation
```

### Task 4: Update CHANGELOG

**File:** `CHANGELOG.md`

Add to `# Next` section:
```markdown
## Phase 13b — Settings Menu (Single User)

### Features

- **Settings Page** — Centralized settings accessible after password login
  - API Key Management: Store Headscale API key (encrypted at rest, no re-entry)
  - Change Password: Update Headplane password with current password verification
  - Theme Selection: Choose light or dark theme with persistence
  - Profile Name: Set optional display name
- **Theme System** — Light/dark theme support with CSS variables and persistence
- **Session Enhancement** — Automatically load stored API key on login

### Backend

- Add `headplane_settings` SQLite table for settings storage
- Add `GET /api/v1/headplane/settings` endpoint (retrieve settings)
- Add `POST /api/v1/headplane/settings` endpoint (update settings)
- Add `POST /api/v1/headplane/change-password` endpoint
- Implement AES-256-GCM encryption for stored API keys
- Implement PBKDF2 key derivation from session token
- Add comprehensive unit tests for settings and encryption

### Frontend

- Add `/settings/profile` route with four sections:
  - Account (password change, session info)
  - Integration (API key input with masked display)
  - Preferences (theme selector)
  - Profile (display name)
- Implement theme system with CSS variables and `data-theme` attribute
- Add "Settings" link to navigation (visible after password login)
- Add form validation and error handling for all settings operations
- Add unit and E2E tests for settings UI

### Security

- API keys encrypted at rest with AES-256-GCM
- Password change requires current password verification (constant-time)
- CSRF protection on all settings endpoints
- Session token used for encryption key derivation

### Documentation

- Created `docs/usage/settings.md` — comprehensive settings guide
- Updated `docs/usage/authentication.md` — settings reference
- Updated README — features list
```

### Task 5: API Documentation

**File:** Create `docs/ref/api/headplane-settings.md`

**Structure:**
```markdown
# Headplane Settings API

## GET /api/v1/headplane/settings

Retrieve current settings for the authenticated user.

**Authentication:** Password session token required

**Request:**
```http
GET /api/v1/headplane/settings
Authorization: Bearer <session-token>
```

**Response:**
```json
{
  "apiKey": "hs_abc123...",
  "theme": "light",
  "profileName": "John Doe"
}
```

**Status Codes:**
- 200: Success
- 401: Unauthorized (invalid or expired session)

[Document other endpoints similarly]
```

### Task 6: Integration Testing Plan

**File:** Create `docs/phase13b-integration-testing.md`

**Test scenarios:**

1. **API Key Storage Flow**
   - Log in with password
   - Navigate to Settings
   - Save API key
   - Log out
   - Log in again
   - Verify API key loaded (Headscale API calls work)

2. **Password Change Flow**
   - Log in with current password
   - Change password in Settings
   - Log out
   - Log in with new password (should succeed)
   - Try logging in with old password (should fail)

3. **Theme Persistence**
   - Log in
   - Change theme to dark
   - Verify theme applied immediately
   - Log out
   - Log in again
   - Verify dark theme loaded

4. **Profile Name**
   - Set profile name in Settings
   - Verify displayed in UI
   - Log out and log in
   - Verify name persisted

5. **Error Handling**
   - Test invalid API key format
   - Test wrong current password for change
   - Test backend offline scenarios

6. **Security**
   - Verify API key encrypted in database
   - Verify unauthenticated requests rejected
   - Verify CSRF protection active

### Task 7: Perform Integration Testing

**When:** After backend and frontend agents complete their work.

**Setup:**
1. Merge backend branch to integration environment
2. Merge frontend branch to integration environment
3. Build and run Headscale + Headplane

**Execute:**
- Run through all test scenarios from Task 6
- Document results (pass/fail for each)
- Take screenshots for documentation
- Report any bugs found

**Tools:**
```bash
# Backend
cd /home/denny/Project/headscale
go build -o headscale cmd/headscale/main.go
./headscale serve

# Frontend
cd /home/denny/Project/headscale/headplane
pnpm install
pnpm dev

# Manual testing in browser
# Automated E2E tests
pnpm test:e2e
```

### Task 8: Document Known Issues

**File:** Update `docs/phase13b-implementation-plan.md`

Add section at end:
```markdown
## Integration Testing Results

**Date:** [date]
**Tester:** Agent 3

### Test Results

| Scenario | Status | Notes |
|----------|--------|-------|
| API Key Storage | ✅ Pass | |
| Password Change | ✅ Pass | |
| Theme Persistence | ✅ Pass | |
| Profile Name | ✅ Pass | |
| Error Handling | ⚠️ Partial | [describe issues] |
| Security | ✅ Pass | |

### Known Issues

1. [Issue description]
   - Severity: High/Medium/Low
   - Workaround: [if any]
   - Fix required: [yes/no]

### Recommendations

[Any improvements or follow-up work]
```
>>>>>>> feature/settings-docs

---

## Files to Create/Modify

### Create:
<<<<<<< HEAD
- `hscontrol/headplane_settings.go` (~250 lines)
- `hscontrol/headplane_settings_test.go` (~200 lines)

### Modify:
- `hscontrol/app.go` (register 3 new endpoints)
- `hscontrol/headplane_auth.go` (password change, session enhancement)
- `hscontrol/headplane_auth_test.go` (update for session enhancement)
=======
- `docs/usage/settings.md`
- `docs/ref/api/headplane-settings.md`
- `docs/phase13b-integration-testing.md`

### Modify:
- `docs/usage/authentication.md`
- `README.md`
- `CHANGELOG.md`
- `docs/phase13b-implementation-plan.md` (add results)
>>>>>>> feature/settings-docs

---

## Testing

<<<<<<< HEAD
```bash
cd /home/denny/Project/headscale-settings-backend

# Run backend tests
go test ./hscontrol -v -run TestHeadplane

# Run all tests
go test ./...

# Check coverage
go test -cover ./hscontrol
```
=======
### Documentation Review
- All docs are clear and accurate
- All code examples work
- All links are valid
- Screenshots match current UI
- Security notes are prominent

### Integration Testing
- All test scenarios pass
- No regressions in Phase 13a
- Backend and frontend work together
- Error handling works correctly
- Security measures verified
>>>>>>> feature/settings-docs

---

## Acceptance Criteria

<<<<<<< HEAD
- [ ] `headplane_settings` SQLite table created
- [ ] API key encrypted with AES-256-GCM
- [ ] `GET /api/v1/headplane/settings` returns decrypted settings
- [ ] `POST /api/v1/headplane/settings` updates and encrypts
- [ ] `POST /api/v1/headplane/change-password` validates and changes password
- [ ] Session loads stored API key automatically
- [ ] All unit tests pass
- [ ] No regressions in Phase 13a functionality
=======
- [ ] `docs/usage/settings.md` created with comprehensive guide
- [ ] `docs/usage/authentication.md` updated with settings reference
- [ ] README features list updated
- [ ] CHANGELOG updated with Phase 13b entries
- [ ] API documentation complete
- [ ] Integration testing performed
- [ ] Test results documented
- [ ] Known issues documented (if any)
- [ ] All documentation links work
- [ ] Screenshots added (if needed)
>>>>>>> feature/settings-docs

---

## Coordination with Other Agents

<<<<<<< HEAD
### With Frontend Agent:
- **API Contract:** Share OpenAPI spec or endpoint signatures
- **Mock Endpoints:** Frontend can mock these endpoints initially
- **Integration Testing:** After both complete, test end-to-end

### With Docs Agent:
- **API Documentation:** Provide endpoint details for docs
- **Security Notes:** Explain encryption approach
=======
### With Backend Agent:
- Get final API endpoint signatures
- Get encryption implementation details
- Get any security notes for docs

### With Frontend Agent:
- Get screenshots of settings page
- Get user flow details
- Get any UI-specific notes for docs

### Integration Testing:
- **Wait for both agents to complete** before testing
- Test their merged work together
- Report any integration issues
- Verify all acceptance criteria met

---

## Timeline

**Phase 1 (Day 1-2):** Documentation drafting (can start immediately)
- Create settings guide structure
- Draft API documentation
- Prepare CHANGELOG entries

**Phase 2 (Day 3-4):** Wait for backend/frontend completion
- Review their implementations
- Update docs with accurate details
- Prepare integration test plan

**Phase 3 (Day 5):** Integration testing
- Set up test environment
- Execute all test scenarios
- Document results
- Report issues
>>>>>>> feature/settings-docs

---

## Completion

When done:
<<<<<<< HEAD
1. Run all tests and verify they pass
2. Commit your changes
3. Push branch: `git push -u origin feature/settings-backend`
4. Create summary: list files changed, tests added, commit SHA
5. Report completion with merge readiness

**Do not merge yet.** Integration testing happens after all agents complete.
=======
1. All documentation complete and reviewed
2. Integration testing passed
3. Commit your changes
4. Push branch: `git push -u origin feature/settings-docs`
5. Create summary: files changed, tests performed, commit SHA
6. Report completion with merge readiness

**Do not merge yet.** Wait for all three branches to be ready, then merge in order:
1. Backend
2. Frontend
3. Docs
>>>>>>> feature/settings-docs

---

## References

- **Implementation Plan:** `docs/phase13b-implementation-plan.md`
<<<<<<< HEAD
- **Phase 13a Auth:** `hscontrol/headplane_auth.go`
- **Existing Tests:** `hscontrol/headplane_auth_test.go`
=======
- **Phase 13a Docs:** `docs/usage/authentication.md`
- **Existing Settings Docs:** `docs/ref/configuration.md`
>>>>>>> feature/settings-docs
- **ROADMAP:** `ROADMAP.md` Phase 13b
