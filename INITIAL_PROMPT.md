# Initial Prompt for Phase 13c Task 3

You are working on **Phase 13c Task 3: Backend User Management API** for the Headscale project.

## Your Task

Create comprehensive user management API endpoints for admins:
1. List all users (GET /api/v1/headplane/users)
2. Get user by ID (GET /api/v1/headplane/users/:id)
3. Update user (PUT /api/v1/headplane/users/:id)
4. Delete user (DELETE /api/v1/headplane/users/:id)
5. Add admin-only authorization middleware

## Important Instructions

- **Read `AGENT_INSTRUCTIONS.md` in this worktree first** for complete details
- Work **only** in this worktree: `/home/denny/Project/headscale-multiuser-user-mgmt-api`
- You are on branch: `feature/multiuser-user-mgmt-api`
- Base commit: `880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a`
- **BLOCKED:** This task depends on Task 1 (Backend User Model & Auth) being merged first
- Follow existing API patterns in `hscontrol/api/v1/`
- Add comprehensive authorization checks
- **Do not merge** when complete - report completion and wait for review

## Key Context

- Task 1 provides basic user CRUD functions
- This task adds full REST API endpoints
- All endpoints must be admin-only
- Must prevent deleting the last admin user

## Dependencies

**⚠️ BLOCKED until Task 1 is merged.**

Task 1 provides:
- `headplane_users` table and model
- RegisterUser, ListUsers, DeleteUser functions
- Session with user ID and role
- Multi-user authentication

## Getting Started

1. **Wait for Task 1 to be merged** to `dev`
2. Read `AGENT_INSTRUCTIONS.md` in full
3. Review existing API handlers in `hscontrol/api/v1/`
4. Check `hscontrol/headplane_users.go` for available functions
5. Start implementation following the steps in AGENT_INSTRUCTIONS.md

Good luck!
