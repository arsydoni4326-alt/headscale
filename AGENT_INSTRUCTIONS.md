# Agent Instructions: Task 1 — Backend API Key Admin Capability

## Overview

**Task:** Ensure API key-authenticated sessions expose admin capability to the frontend and all user management endpoints enforce admin-only access.

**Branch:** `feature/knowngap-backend-admin`  
**Worktree:** `/home/denny/Project/headscale-knowngap-backend-admin`  
**Base Commit:** `7ec81109f16952e5730837d1d908d36c0df2935a`

## Objective

Close the "Known Gap" in Phase 13c by ensuring:
1. API key authentication sessions correctly expose admin role/capability.
2. All Headplane user management endpoints enforce admin-only access.
3. Backend returns proper authorization errors for non-admin attempts.

## Context

Currently, the user management UI is only accessible when logged in with Headplane credentials. When accessing via API key, the `/admin/users` page shows only Headscale user namespace management, not Headplane dashboard user management.

The backend must:
- Ensure API key sessions include admin capability data.
- Enforce admin-only access on all user management endpoints.
- Return consistent authorization errors.

## Files to Modify

### Primary Files
- `hscontrol/headplane_users.go` — user management endpoints
- `hscontrol/headplane_auth.go` — authentication and authorization logic
- `hscontrol/db/headplane_users.go` — database layer for user queries

### Test Files
- `hscontrol/headplane_users_test.go` — add/extend tests for API key admin flows
- `hscontrol/headplane_auth_test.go` — add/extend auth tests
- `hscontrol/servertest/` — if applicable

## Implementation Requirements

### 1. API Key Session Context
- Review how API key authentication creates session context.
- Ensure admin role/capability is exposed in API responses.
- Verify API key sessions include all necessary user metadata.

### 2. Endpoint Authorization
- All user management endpoints must check admin role:
  - `GET /api/v1/headplane/users` — list users
  - `POST /api/v1/headplane/users` — create user
  - `PUT /api/v1/headplane/users/:id` — update user
  - `DELETE /api/v1/headplane/users/:id` — delete user
- Return `401 Unauthorized` or `403 Forbidden` with clear error messages for non-admin.

### 3. Error Handling
- Use consistent error types (e.g., `ErrHeadplaneUserNotAuthorized`).
- Return JSON error responses with clear messages.
- Log authorization failures appropriately.

### 4. Testing
- Add unit tests for API key admin authorization.
- Test all user management endpoints with:
  - Admin API key (should succeed)
  - Non-admin API key (should fail with 403)
  - No authentication (should fail with 401)
- Test error responses are well-formed JSON.

## Acceptance Criteria

- [ ] API key-authenticated sessions expose admin role/capability
- [ ] All user management endpoints enforce admin-only access
- [ ] Non-admin API key sessions receive proper 403 errors
- [ ] All new/modified code has unit tests
- [ ] All tests pass (`go test ./hscontrol/...`)
- [ ] No regression in existing auth flows
- [ ] Code follows existing patterns and conventions

## Dependencies

- None (this task is independent)

## Constraints

- Work only in this worktree
- Do not modify unrelated auth logic
- Follow existing error handling patterns
- Preserve existing API key and session functionality

## Testing Commands

```bash
# Run all hscontrol tests
go test ./hscontrol/...

# Run specific test files
go test ./hscontrol/headplane_users_test.go
go test ./hscontrol/headplane_auth_test.go

# Run with race detection
go test -race ./hscontrol/...
```

## Documentation References

- `ROADMAP.md` — Phase 13c Known Gap description
- `AGENTS.md` — Project guidelines and architecture
- `CONTRIBUTING.md` — Code conventions
- `docs/ref/api/headplane-settings.md` — API reference
- `docs/usage/authentication.md` — Authentication guide

## Completion Checklist

- [ ] Read all relevant documentation
- [ ] Understand existing API key auth flow
- [ ] Implement admin capability exposure
- [ ] Harden all user management endpoints
- [ ] Add comprehensive tests
- [ ] Run all tests and verify they pass
- [ ] Commit changes with clear message
- [ ] Document any limitations or follow-up needs

## Commit Message Format

```
backend: ensure API key sessions expose admin capability

- Expose admin role/capability in API key session context
- Enforce admin-only access on all user management endpoints
- Add authorization tests for API key admin flows
- Return proper 403 errors for non-admin attempts

Closes Known Gap in Phase 13c (backend component)
```

## Notes

- This task unblocks the frontend tasks (Task 2 and Task 3)
- Frontend needs to check admin capability in session/API key context
- Coordinate with frontend agents on expected API response structure
