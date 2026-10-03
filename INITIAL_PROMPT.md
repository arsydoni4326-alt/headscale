# Initial Prompt for Phase 13c Task 2

You are working on **Phase 13c Task 2: Backend Per-User Settings** for the Headscale project.

## Your Task

Refactor the settings storage from single-user to per-user by:
1. Adding a `user_id` column to the `headplane_settings` table
2. Updating all settings CRUD functions to filter by user ID
3. Migrating existing single-user settings to the default admin user
4. Adding comprehensive tests for per-user isolation

## Important Instructions

- **Read `AGENT_INSTRUCTIONS.md` in this worktree first** for complete details
- Work **only** in this worktree: `/home/denny/Project/headscale-multiuser-per-user-settings`
- You are on branch: `feature/multiuser-per-user-settings`
- Base commit: `880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a`
- **BLOCKED:** This task depends on Task 1 (Backend User Model & Auth) being merged first
- Review existing settings code in `hscontrol/headplane_settings.go`
- Follow existing database migration patterns
- **Do not merge** when complete - report completion and wait for review

## Key Context

- Phase 13b implemented single-user settings storage
- Settings table currently has a single-row constraint (id=1)
- Need to transition to multi-row with user_id foreign key
- Each user should have completely isolated settings

## Dependencies

**⚠️ BLOCKED until Task 1 is merged.**

Task 1 provides:
- `headplane_users` table
- User ID in session context
- Multi-user authentication

## Getting Started

1. **Wait for Task 1 to be merged** to `dev`
2. Read `AGENT_INSTRUCTIONS.md` in full
3. Review `hscontrol/headplane_settings.go` and `hscontrol/db/db.go`
4. Plan the database migration carefully
5. Start implementation following the steps in AGENT_INSTRUCTIONS.md

Good luck!
