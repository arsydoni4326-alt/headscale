# Phase 13c Task 4: Frontend Multi-User Login

**Worktree:** `/home/denny/Project/headscale-multiuser-frontend-login`  
**Branch:** `feature/multiuser-frontend-login`  
**Base Commit:** `880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a`

## Objective

Update Headplane login UI for multi-user:
- Add username field to login form
- Update API calls to send username + password
- Handle session for authenticated user
- Show username in UI after login

## Dependencies

**Blocked until Task 1 (Backend User Model & Auth) is merged.**

## Context

Current login (Phase 13a):
- Single password field
- POST /api/v1/headplane/login with `{password}`

New login:
- Username + password fields
- POST /api/v1/headplane/login with `{username, password}`

## Files to Modify

- `headplane/app/routes/auth/login/action.ts` - Update API call
- `headplane/app/routes/auth/login/route.tsx` - Add username field
- `headplane/app/server/headscale/api/` - Update API client
- Tests

## Implementation Steps

1. Add username input field to login form
2. Update form validation for username
3. Update login API call to include username
4. Store username in session/context after login
5. Display username in header/nav
6. Update tests

## UI Changes

**Login Form:**
```
Username: [__________]
Password: [__________]
[ ] Remember me
[Login]
```

## Testing

```bash
cd /home/denny/Project/headscale-multiuser-frontend-login/headplane
pnpm test
pnpm typecheck
```

## Acceptance Criteria

- [ ] Username field in login form
- [ ] Login works with username + password
- [ ] Username displayed after login
- [ ] Tests pass
- [ ] Typecheck passes

## Completion

Push: `git push -u origin feature/multiuser-frontend-login`

**Do not merge** until Task 1 is merged.
