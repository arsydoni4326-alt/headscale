# Phase 13c Task 1: Backend User Model & Auth

**Worktree:** `/home/denny/Project/headscale-multiuser-backend-auth`  
**Branch:** `feature/multiuser-backend-auth`  
**Base Commit:** `880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a`

## Objective

Implement foundational multi-user backend:
- Add `headplane_users` table for multiple user accounts
- Migrate single-user data to new schema
- Implement user registration (admin-only)
- Implement multi-user login with username/password
- Session management with user ID
- Password hashing with bcrypt
- Comprehensive tests

## Context

- Phase 13a complete: single-user password auth exists
- Phase 13b complete: single-user settings storage exists
- This task enables multiple users with individual credentials
- Must preserve backward compatibility

## Database Schema

```sql
CREATE TABLE headplane_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT DEFAULT 'user',
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

## Key Changes

1. **New table**: `headplane_users` with username, password_hash, role
2. **Migration**: Auto-create admin user from config password on first startup
3. **Login**: Change from `{password}` to `{username, password}`
4. **Session**: Store user ID, not just token
5. **Endpoints**: 
   - POST /api/v1/headplane/register (admin-only)
   - GET /api/v1/headplane/users (admin-only)
   - DELETE /api/v1/headplane/users/:id (admin-only)

## Files to Create

- `hscontrol/db/headplane_users.go` - Schema definition
- `hscontrol/headplane_users.go` - CRUD functions
- `hscontrol/headplane_users_test.go` - Tests

## Files to Modify

- `hscontrol/headplane_auth.go` - Multi-user login logic
- `hscontrol/app.go` - Register endpoints
- `hscontrol/db/db.go` - Auto-migrate table

## Implementation Steps

1. Define HeadplaneUser struct in db/headplane_users.go
2. Add auto-migration in db/db.go
3. Create migration function for single→multi user
4. Implement RegisterUser, ListUsers, DeleteUser
5. Update login handler for username/password
6. Update session to store user ID
7. Add admin-only middleware for user management
8. Write comprehensive tests

## Testing

```bash
cd /home/denny/Project/headscale-multiuser-backend-auth
go test ./hscontrol -v -run TestHeadplane
go test ./...
```

## Acceptance Criteria

- [ ] `headplane_users` table created
- [ ] Bcrypt password hashing (cost 12)
- [ ] User registration endpoint works
- [ ] Multi-user login with username/password
- [ ] Session stores user ID
- [ ] List/delete users endpoints
- [ ] Auto-migration from single-user
- [ ] All tests pass
- [ ] No regressions

## Dependencies

**This is the foundational task.** All other Phase 13c tasks depend on this being merged first.

## Completion

1. Run all tests
2. Commit changes
3. Push: `git push -u origin feature/multiuser-backend-auth`
4. Report completion

**Do not merge.** Must be reviewed first.

## References

- `hscontrol/headplane_auth.go` - Current auth
- `hscontrol/headplane_settings.go` - Settings storage
- `ROADMAP.md` Phase 13c
