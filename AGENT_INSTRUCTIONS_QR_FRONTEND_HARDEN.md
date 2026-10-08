# Agent Instructions: QR Frontend Hardening

## Task Overview

**Task:** Frontend QR UX Hardening  
**Roadmap Item:** Phase 17 — Non-OIDC Interactive Registration  
**Branch:** `feature/qr-frontend-harden`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-qr-frontend-harden`  
**Base Commit:** `f808528f` (main)

## Objective

Harden the frontend QR scan UI, error handling, and user selection for non-OIDC interactive registration. Ensure clear feedback for edge cases and comprehensive test coverage for:
- Expired QR codes
- Invalid QR payloads
- User selection errors
- Scan errors and camera issues
- Approval flow edge cases

## Scope

### In Scope
- Review and harden QR scan UI (`headplane/app/routes/machines/scan-qr.tsx`)
- Improve error handling and user feedback
- Add frontend tests for edge cases
- Validate scanner works on supported browsers
- Improve accessibility and UX consistency

### Out of Scope
- Backend changes (handled by separate task)
- Documentation updates (handled by separate task)
- Integration/E2E tests (handled by separate task)
- Changes to core approval flow outside QR scan

## Files/Modules Expected to Change

- `headplane/app/routes/machines/scan-qr.tsx`
- `headplane/app/components/*` (if error/feedback components need updates)
- `headplane/tests/*` (component/route tests)
- Possibly: `headplane/app/lib/qr-validation.ts` (if validation logic needs extraction)

## Dependencies

None. This task is independent and can proceed immediately.

## Testing Requirements

1. **Component Tests:**
   - Test QR scan component with expired QR
   - Test QR scan component with invalid payload
   - Test user selection validation
   - Test error state rendering

2. **Error Handling Tests:**
   - Test camera permission denied
   - Test scan timeout
   - Test malformed QR data

3. **Validation:**
   - Run: `cd headplane && pnpm test:component`
   - Run: `cd headplane && pnpm run typecheck` (may have pre-existing errors)
   - Run: `cd headplane && pnpm exec oxlint app/routes/machines/scan-qr.tsx`
   - Run: `cd headplane && pnpm run build`

## Acceptance Criteria

- [ ] All edge cases for QR scanning are handled with clear feedback
- [ ] Error messages are user-friendly and actionable
- [ ] User selection validation is robust
- [ ] All new tests pass
- [ ] No regression in existing tests
- [ ] Code is linted and builds successfully
- [ ] Changes are committed to `feature/qr-frontend-harden` branch

## Constraints

- Preserve existing QR scan functionality
- Do not modify backend code
- Follow Headplane design system and component library
- Use existing error/toast patterns
- Maintain accessibility standards

## Relevant Documentation

- `ROADMAP.md` — Phase 17 description
- `session.md` — Phase 17 completion notes
- `docs/usage/registration.md` — Non-OIDC registration flow
- `headplane/docs/` — Frontend architecture and patterns
- `AGENTS.md` — Project interaction rules

## Implementation Notes

1. Start by reading the current QR scan implementation:
   - `headplane/app/routes/machines/scan-qr.tsx`
   - Related components in `headplane/app/components/`

2. Identify edge cases not currently covered by tests

3. Add tests first (TDD approach preferred)

4. Implement hardening and better error feedback

5. Verify all tests pass

6. Commit incrementally with clear messages

## Worktree Setup

```bash
cd /home/denny/Project/headscale-project/headscale-qr-frontend-harden
git status  # Verify on feature/qr-frontend-harden branch
git log --oneline -n 1  # Should show f808528f
cd headplane
pnpm install  # Ensure dependencies are installed
```

## Known Issues

- Headplane may have pre-existing typecheck errors unrelated to QR scan
- Node version should be 24.2+ (workspace may be on 22.22.1)
- Some dependencies (`html5-qrcode`, `better-sqlite3`) may need local installation

These are known project issues and should not block this task.

## Final Verification

Before considering this task complete:
- [ ] All tests pass
- [ ] Code is linted
- [ ] Build succeeds
- [ ] Changes are committed
- [ ] Branch pushed to remote (if applicable)

## Contact

If you discover scope creep or need architectural decisions outside this task, stop and request clarification.
