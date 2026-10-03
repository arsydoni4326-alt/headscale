# Phase 13c Task 3: Backend User Management API

**Worktree:** `/home/denny/Project/headscale-multiuser-user-mgmt-api`  
**Branch:** `feature/multiuser-user-mgmt-api`  
**Base Commit:** `880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a`

## Objective

Add comprehensive user management API endpoints:
- List all users (admin-only)
- Get user by ID (admin-only)
- Update user (admin-only, change username/role)
- Delete user (admin-only)
- Optional: Avatar upload endpoint

## Dependencies

**Blocked until Task 1 (Backend User Model & Auth) is merged.**

## Context

Task 1 provides basic RegisterUser/ListUsers/DeleteUser functions. This task adds:
- Full CRUD API endpoints
- Admin authorization checks
- Optional avatar support

## API Endpoints

1. **GET /api/v1/headplane/users** - List all users (admin-only)
2. **GET /api/v1/headplane/users/:id** - Get user by ID (admin-only)
3. **PUT /api/v1/headplane/users/:id** - Update user (admin-only)
4. **DELETE /api/v1/headplane/users/:id** - Delete user (admin-only)
5. **POST /api/v1/headplane/users/:id/avatar** - Upload avatar (optional)

## Authorization

All endpoints require:
- Valid session token
- Authenticated user has `role = 'admin'`

Return 403 Forbidden if not admin.

## Files to Create

- `hscontrol/api/v1/headplane_users.go` - User management handlers

## Files to Modify

- `hscontrol/app.go` - Register endpoints
- `hscontrol/headplane_users.go` - Add UpdateUser, GetUser functions
- Tests

## Implementation Steps

1. Implement admin authorization middleware
2. Add GetUser(id) function
3. Add UpdateUser(id, username, role) function
4. Create API handlers for GET/PUT/DELETE
5. Add comprehensive tests
6. Optional: Add avatar upload support

## Testing

```bash
cd /home/denny/Project/headscale-multiuser-user-mgmt-api
go test ./hscontrol -v -run TestHeadplaneUserManagement
```

## Acceptance Criteria

- [ ] List/get/update/delete user endpoints
- [ ] Admin-only authorization enforced
- [ ] Cannot delete last admin user
- [ ] All tests pass

## Completion

Push: `git push -u origin feature/multiuser-user-mgmt-api`

**Do not merge** until Task 1 is merged.
