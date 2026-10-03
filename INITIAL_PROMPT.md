# Initial Prompt for Phase 13c Task 1

You are working on **Phase 13c Task 1: Backend User Model & Auth** for the Headscale project.

## Your Task

Implement the foundational multi-user backend infrastructure by:
1. Creating a `headplane_users` table for storing multiple user accounts
2. Implementing user registration (admin-only)
3. Updating the login flow to accept username + password
4. Migrating existing single-user password to a default admin user
5. Adding comprehensive tests

## Important Instructions

- **Read `AGENT_INSTRUCTIONS.md` in this worktree first** for complete details
- Work **only** in this worktree: `/home/denny/Project/headscale-multiuser-backend-auth`
- You are on branch: `feature/multiuser-backend-auth`
- Base commit: `880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a`
- Follow existing code conventions in `hscontrol/`
- Read existing auth code in `hscontrol/headplane_auth.go` before starting
- Add comprehensive unit tests for all new functionality
- **Do not merge** when complete - report completion and wait for review

## Key Context

- Phase 13a (single-user password auth) is complete
- Phase 13b (single-user settings) is complete
- This task is the foundation for all other Phase 13c work
- Other agents are blocked until this is merged

## Getting Started

1. Read `AGENT_INSTRUCTIONS.md` in full
2. Review existing code in `hscontrol/headplane_auth.go`
3. Check the database schema in `hscontrol/db/`
4. Start implementation following the steps in AGENT_INSTRUCTIONS.md

Good luck!
