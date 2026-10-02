# Phase 13 Agent Instructions: Simple Password Login for Headplane

## Overview

This document archives the parallel implementation tasks for Phase 13. All three tasks have been completed and merged into main.

---

## Task 1: Backend Password Authentication

**Agent:** Backend Auth Logic  
**Worktree:** `/home/denny/Project/headscale-password-backend`  
**Branch:** `feature/password-backend`  
**Status:** ✅ Complete

### Objective
Implement password-based authentication for Headplane web UI that:
- Does NOT require database schema changes
- Stores password in configuration file or environment variable
- Provides an alternative to API key authentication
- Is scoped specifically to Headplane (does not grant general API access)

### Completed Implementation
- Added `POST /api/v1/headplane/login` endpoint
- Session management with 24-hour expiry
- Rate limiting (5 attempts per minute per IP)
- Constant-time password comparison
- Configuration via `headplane.password` in config.yaml or `HEADSCALE_HEADPLANE_PASSWORD` env var
- Comprehensive unit tests

### Files Changed
- `hscontrol/headplane_auth.go` (new)
- `hscontrol/headplane_auth_test.go` (new)
- `hscontrol/types/config.go`
- `hscontrol/app.go`
- `config-example.yaml`

---

## Task 2: Documentation

**Agent:** Documentation  
**Worktree:** `/home/denny/Project/headscale-password-docs`  
**Branch:** `feature/password-docs`  
**Status:** ✅ Complete

### Objective
Update all relevant documentation to cover password-based login for Headplane, including setup, configuration, security, and usage.

### Completed Documentation
- Created `docs/usage/authentication.md` - Complete authentication guide
- Created `docs/troubleshooting.md` - Troubleshooting guide for auth issues
- Updated `docs/ref/configuration.md` - Password configuration reference
- Updated `config-example.yaml` - Comprehensive inline documentation
- Updated `README.md` - Authentication section

### Coverage
- Configuration via config.yaml and environment variable
- Security best practices (HTTPS, strong passwords, rate limiting)
- Login API endpoint documentation
- Session management details
- Troubleshooting common issues
- Clear distinction between password auth (web UI) and API keys (automation)

---

## Task 3: Frontend Password Login UI

**Agent:** Frontend Login UI  
**Worktree:** `/home/denny/Project/headscale-password-frontend`  
**Branch:** `feature/password-frontend`  
**Status:** ✅ Complete

### Objective
Update Headplane login form to support password-based authentication alongside the existing API key method.

### Completed Implementation
- Password/API Key toggle in login form (password is default)
- Password field with show/hide toggle
- Backend API integration (`POST /api/v1/headplane/login`)
- Session token management
- Error handling (401 invalid, 429 rate limit, 503 not configured)
- Comprehensive unit and E2E tests
- Preserved API key login functionality

### Files Changed (in headplane submodule)
- `headplane/app/server/headscale/api/transport.ts`
- `headplane/app/server/headscale/api/index.ts`
- `headplane/app/server/web/auth.ts`
- `headplane/app/server/db/schema.ts`
- `headplane/app/routes/auth/login/action.ts`
- `headplane/app/routes/auth/login/page.tsx`
- `headplane/tests/unit/auth/auth-service.test.ts`
- `headplane/tests/unit/auth/password-login-action.test.ts` (new)
- `headplane/tests/e2e/login.spec.ts`

---

## Security Considerations

- Password transmitted over HTTPS only (secure cookies)
- Rate limiting prevents brute-force attacks (5 attempts per minute per IP)
- Constant-time password comparison prevents timing attacks
- Session tokens: 256-bit random, 24-hour expiry
- Automatic session cleanup (5-minute interval)
- Password never serialized in JSON output

---

## Backward Compatibility

- API key authentication fully preserved and functional
- No database migrations required
- Password authentication is optional
- Existing deployments continue working without changes

---

## Integration Status

All three branches merged into `main` successfully:
1. Backend (feature/password-backend)
2. Docs (feature/password-docs)
3. Frontend (feature/password-frontend)

---

## References

- Backend: `hscontrol/headplane_auth.go`
- Documentation: `docs/usage/authentication.md`
- Frontend: See `PASSWORD_LOGIN_SUMMARY.md`
- Roadmap: `ROADMAP.md` Phase 13
