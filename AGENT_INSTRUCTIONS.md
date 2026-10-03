# Phase 13c Task 2: Backend Per-User Settings

**Worktree:** `/home/denny/Project/headscale-multiuser-per-user-settings`  
**Branch:** `feature/multiuser-per-user-settings`  
**Base Commit:** `880a8a7d9abb0cfb77fcb20e767aa67cee4cc02a`

## Objective

Refactor settings storage from single-user to per-user:
- Modify `headplane_settings` table to link to user ID
- Update settings API endpoints to be user-aware
- Migrate existing single-user settings to default admin user
- Each user gets their own API key, theme, profile

## Dependencies

**Blocked until Task 1 (Backend User Model & Auth) is merged.**

## Context

- Phase 13b has single-user settings in `headplane_settings` table
- Current schema has single-row constraint (id=1)
- Need to change to multi-row, keyed by user_id

## Schema Changes

**Before:**
```sql
CREATE TABLE headplane_settings (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  api_key_encrypted TEXT,
  ...
);
```

**After:**
```sql
CREATE TABLE headplane_settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER UNIQUE NOT NULL,
  api_key_encrypted TEXT,
  api_key_salt TEXT,
  theme TEXT DEFAULT 'light',
  profile_name TEXT,
  updated_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES headplane_users(id) ON DELETE CASCADE
);
```

## Key Changes

1. Remove single-row constraint
2. Add `user_id` column with UNIQUE constraint
3. Add foreign key to `headplane_users`
4. Update GetSettings/UpdateSettings to filter by user_id
5. Migrate existing settings row to admin user (id=1)

## Files to Modify

- `hscontrol/headplane_settings.go` - Add user_id parameter
- `hscontrol/db/db.go` - Migration for schema change
- `hscontrol/headplane_settings_test.go` - Update tests
- `hscontrol/app.go` - Update endpoint handlers

## Implementation Steps

1. Add migration to modify headplane_settings table
2. Update HeadplaneSettings struct with UserID field
3. Modify GetSettings(userID) and UpdateSettings(userID, ...)
4. Update API handlers to extract user ID from session
5. Migrate existing settings row to user_id=1 (admin)
6. Update all tests for multi-user settings

## Testing

```bash
cd /home/denny/Project/headscale-multiuser-per-user-settings
go test ./hscontrol -v -run TestHeadplaneSettings
go test ./...
```

## Acceptance Criteria

- [ ] Settings table has user_id column
- [ ] Foreign key to headplane_users
- [ ] Each user has independent settings
- [ ] Settings API filtered by authenticated user
- [ ] Existing settings migrated to admin
- [ ] All tests pass

## Completion

1. Run tests
2. Commit and push: `git push -u origin feature/multiuser-per-user-settings`
3. Report completion

**Do not merge** until Task 1 is merged and this is reviewed.
