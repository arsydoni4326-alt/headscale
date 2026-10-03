# Agent 3: Documentation & Integration Testing

**Task:** Phase 13b Part 5 & 6 — Documentation & Integration Testing  
**Worktree:** `/home/denny/Project/headscale-settings-docs`  
**Branch:** `feature/settings-docs`  
**Base Commit:** `d95c5eff`

---

## Objective

Document Phase 13b settings features and perform integration testing:
- Create comprehensive user documentation
- Update authentication guide with settings references
- Document API endpoints
- Update README and CHANGELOG
- Perform end-to-end integration testing after backend & frontend complete
- Document any issues found during integration

---

## Context

Phase 13b adds a Settings menu to Headplane. Users can:
- Save Headscale API key (no re-entry at login)
- Change Headplane password
- Select theme (light/dark)
- Set profile name

Backend and frontend agents are implementing in parallel. You'll document their work and test integration.

---

## Implementation Tasks

### Task 1: Create Settings Documentation

**File:** Create `docs/usage/settings.md`

**Structure:**
```markdown
# Settings

The Headplane Settings page allows you to customize your experience and manage your Headscale integration.

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

## Security Considerations

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

---

## Files to Create/Modify

### Create:
- `docs/usage/settings.md`
- `docs/ref/api/headplane-settings.md`
- `docs/phase13b-integration-testing.md`

### Modify:
- `docs/usage/authentication.md`
- `README.md`
- `CHANGELOG.md`
- `docs/phase13b-implementation-plan.md` (add results)

---

## Testing

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

---

## Acceptance Criteria

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

---

## Coordination with Other Agents

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

---

## Completion

When done:
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

---

## References

- **Implementation Plan:** `docs/phase13b-implementation-plan.md`
- **Phase 13a Docs:** `docs/usage/authentication.md`
- **Existing Settings Docs:** `docs/ref/configuration.md`
- **ROADMAP:** `ROADMAP.md` Phase 13b
