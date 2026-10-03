# Initial Prompt for Phase 13c Task 4

You are working on **Phase 13c Task 4: Frontend Multi-User Login** for the Headscale project.

## Your Task

Update the Headplane login UI to support multi-user authentication:
1. Add username input field to login form
2. Update login API call to send username + password
3. Store username in session/context after login
4. Display username in the UI after successful login
5. Update all tests

## Important Instructions

- **Read `AGENT_INSTRUCTIONS.md` in this worktree first** for complete details
- Work **only** in this worktree: `/home/denny/Project/headscale-multiuser-frontend-login`
- You are on branch: `feature/multiuser-frontend-login`
- Base commit: `880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a`
- **BLOCKED:** This task depends on Task 1 (Backend User Model & Auth) being merged first
- Review existing login code in `headplane/app/routes/auth/login/`
- Follow existing form patterns and validation
- **Do not merge** when complete - report completion and wait for review

## Key Context

- Current login only has password field (single-user from Phase 13a)
- Backend now expects `{username, password}` instead of just `{password}`
- Need to update both UI and API client

## Dependencies

**⚠️ BLOCKED until Task 1 is merged.**

Task 1 provides:
- Updated POST /api/v1/headplane/login endpoint
- Accepts `{username, password}` payload
- Returns session with user information

## Getting Started

1. **Wait for Task 1 to be merged** to `dev`
2. Read `AGENT_INSTRUCTIONS.md` in full
3. Review `headplane/app/routes/auth/login/route.tsx`
4. Check `headplane/app/server/headscale/api/` for API client
5. Start implementation following the steps in AGENT_INSTRUCTIONS.md

Good luck!
