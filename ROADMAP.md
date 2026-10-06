# Roadmap

This document tracks the high-level direction of the `arsydoni4326-alt` fork of
Headscale and its frontend, Headplane. It is maintained alongside the
implementation and kept aligned with the actual project state.

> For past changes, see the [CHANGELOG](./CHANGELOG.md). For the frontend's own
> roadmap, specification, and architecture, see
> [`headplane/docs/`](./headplane/docs/).

## Preamble: Feature Preservation

All existing features and customizations in this repository are load-bearing and
must be preserved. This includes, but is not limited to:

- The `-arsydoni4326-alt` version suffix convention (tags, changelog entries,
  package versions).
- The update-check feature: `hscontrol/updatecheck/` (backend) and
  `headplane/app/update-check/` (frontend). Both are marked "DO NOT REMOVE" and
  must survive upstream merges.
- The `GET /api/v1/update-check` endpoint and its no-auth behavior.
- The Dockerfile build args (`APP_VERSION`, `APP_COMMIT`, `BUILD_DATE`) and the
  `__COMMIT_HASH__` Vite global.
- The fork repository targets (`github.com/arsydoni4326-alt/headscale` and
  `github.com/arsydoni4326-alt/headplane`).

Every roadmap item below is **additive and non-destructive**. No item may remove,
rename, or break an existing feature without explicit review and a documented
migration path.

## Overview

The fork tracks upstream Headscale closely while adding a small set of
fork-specific features, most notably the update checker that spans both the Go
backend (`hscontrol/updatecheck/`) and the Headplane frontend
(`headplane/app/update-check/`). The roadmap is organized into phases that
harden what exists today before expanding into new territory.

### Status legend

- **[Implemented]** — shipped and verified in the codebase.
- **[Partially implemented]** — shipped with known gaps or pending follow-ups.
- **[Planned]** — tracked in this roadmap or the Headplane roadmap, not yet
  implemented.
- **[Proposed]** — recommended by this analysis; not an existing requirement.
  Proposed items are additive and non-destructive, and require explicit
  approval before implementation.

### Current state

- Phases 1-12 are complete. Phase 13 and Phase 13c are complete; Headplane now
  supports a single local administrator with password stored as a bcrypt hash in
  configuration. Multi-user local authentication has been retired. Latest fork
  releases: Headscale `v0.34.0-arsydoni4326-alt`, Headplane `v0.8.3-arsydoni4326-alt`.
- Phase 11 below is **proposed** and must not be implemented without
  explicit approval.
- Phase 12 is **partially implemented**: the backend approval API is done; the
  Headplane frontend work is planned.

## Phase 1 — Foundation and Documentation (Completed)

Released in `v0.29.9-arsydoni4326-alt`. See [CHANGELOG](./CHANGELOG.md).

- [x] Create root-level `SPECIFICATION.md` and `ARCHITECTURE.md` describing the
      fork as a whole (backend + frontend), including the fork-specific features
      and the feature-preservation rule above.
- [x] Add a CI check that verifies `hscontrol/updatecheck/` and
      `headplane/app/update-check/` still exist after every merge, so an
      upstream merge cannot silently delete them.
- [x] Document the fork-specific features in the user-facing docs site
      (`docs/`), including the update-check endpoint and the version suffix
      convention.
- [x] Add the `/api/v1/update-check` endpoint to the OpenAPI specification so it
      is discoverable alongside the rest of the v1 API.

## Phase 2 — Update Checker Hardening (Completed)

See [CHANGELOG](./CHANGELOG.md) for the list of changes in this phase.

## Phase 3 — UI/UX Polish

### Part 1: Audit and Visual Documentation

- [x] UI Inventory & Audit: read all template files, catalog UI elements and
      shared components, identify inconsistencies and polish opportunities.
- [x] Visual Documentation Base: - Created `docs/usage/ui-guide.md` (placeholder screenshots, Mermaid flow
      diagrams, contributor guidance, maintenance checklist). - Created `docs/assets/screenshots/` and `docs/assets/diagrams/` directories. - Updated `mkdocs.yml` to include UI Guide in nav.
- [ ] Populate screenshot placeholders with actual images.

### Part 2: Visual Consistency & Minor Style Improvements

- [x] status-circle.tsx: Added `role="img"` for better SVG accessibility
- [x] table-list.tsx: Increased item padding from `p-2` to `p-3` for better touch targets and visual breathing room
- [x] token-list.tsx: Simplified add button icon sizing (removed oversized `size={30}` with `p-1`, replaced with standard `size-4`)
- [x] Verified all components use consistent focus-ring patterns, color palette, and spacing

### Part 3: Interaction & Feedback Polish

- [x] button.tsx: Added `active:scale-[0.98]` press-state feedback for tactile response
- [x] dialog.tsx: Added entrance animation (opacity + scale) to Dialog panel for polished feel
- [x] Verified all interactive elements have visible focus states, ARIA labels, and keyboard navigation
- [x] Toast/notification system already in place (toast-provider, CodeBlock copy feedback, Attribute copy feedback)

### Part 4: Responsive & Device Testing

- [x] Verified responsive layout patterns across all pages:
  - Container utility with breakpoint-aware padding
  - Grid layouts (grid-cols-2 sm:grid-cols-4 on home page)
  - Overflow handling for tables (overflow-x-auto on machine table)
  - Mobile navigation with horizontal scroll (block md:hidden)
  - Viewport-constrained dialogs (max-h-[90dvh])
  - Responsive width classes on attributes (w-1/3 sm:w-1/4 lg:w-1/3)
  - Settings page uses sm:w-2/3 for content width
- [x] header.tsx: Added transition-opacity to mobile nav for smoother appearance

### Part 5: Documentation & Testing

- [x] Updated ROADMAP.md to reflect Parts 2-4 completion
- [x] Updated session.md with Phase 3 progress
- [x] No UI component tests exist in the headplane frontend — adding a testing framework (vitest, testing-library) is out of scope for this polish phase; noted as a known limitation

### Part 6: Review & Acceptance

- [ ] Pending manual review and stakeholder sign-off

## Phase 4 — Feature Expansion (Completed)

See [CHANGELOG](./CHANGELOG.md) for the list of changes in this phase.

## Phase 5 — Testing, Performance, and CI/CD (Completed)

See [CHANGELOG](./CHANGELOG.md) for the list of changes in this phase.

## Phase 6 — Community and Ecosystem (Completed)

See [CHANGELOG](./CHANGELOG.md) for the list of changes in this phase.

## Phase 7 — Testing and Technical Hardening [Implemented]

**Objective:** Close the remaining testing, documentation, and dead-code gaps
that make future work riskier and the codebase harder to maintain.

**Problems addressed:**

- Headplane has no UI component tests — only service-level unit tests,
  integration tests, and e2e/a11y tests (`headplane/tests/unit/` covers
  services, config, and utilities; the shared component library in
  `headplane/app/components/` is untested). This was noted as a known
  limitation in Phase 3.
- `GET /api/v1/derp` is missing from the OpenAPI specification
  (`openapi/v1/headscale.yaml` documents `/api/v1/update-check` but not the
  DERP endpoint), so the fork's own API surface is incompletely documented.
- Documentation inconsistency: `headplane/docs/configuration/index.md` still
  says environment overrides require `HEADPLANE_LOAD_ENV_OVERRIDES=true`,
  while `headplane/docs/SPECIFICATION.md` (NFR-2) and
  `headplane/docs/ARCHITECTURE.md` state the variable is deprecated and env
  overrides are always loaded.
- The legacy `useUpdateCheck.ts` hook is reported as unwired in `session.md`
  but is still referenced by `headplane/app/root.tsx`,
  `headplane/app/layout/header.tsx`, and `headplane/app/update-check/index.ts`;
  it needs verification and either wiring or removal.
- The upstream `DisableUpdateCheck` config option
  (`hscontrol/types/config.go`) is not wired to the fork's update-check
  endpoint; the fork endpoint always responds.
- The UI guide (`docs/usage/ui-guide.md`) still contains screenshot
  placeholders.

**Features/improvements:**

- [x] Add a UI component testing framework (Vitest + React Testing Library) to
      Headplane and cover the shared component library (button, chip,
      status-circle, table-list, token-list, tabs, switch). Wire into CI.
- [x] Add the `/api/v1/derp` handler to the OpenAPI spec (code-first Huma, so
      this is a handler annotation change plus regenerated spec). Add a test
      that asserts the endpoint and its schema appear in the emitted spec.
- [x] Fix the `HEADPLANE_LOAD_ENV_OVERRIDES` documentation inconsistency.
- [x] Verify and remove the legacy `useUpdateCheck.ts` hook (dead code; only
      `useUpdateCheckContext` from `UpdateCheckProvider` is used).
- [x] Wire `DisableUpdateCheck` into the fork's update-check endpoint: when
      `disable_check_updates` is true, the endpoint skips the remote comparison
      and returns version info only.
- [x] Convert the UI guide screenshot placeholders into a tracked maintenance
      task.

**Technical work:** Headplane test tooling + component tests; backend OpenAPI
annotation; docs corrections; dead-code cleanup.

**UI/UX work:** None beyond test coverage of existing components.

**Dependencies:** None.

**Expected outcome:** Higher confidence in UI changes, complete API
documentation, consistent docs, and a cleaner codebase.

**Priority:** High (foundational — everything below builds on it).

## Phase 8 — UI/UX and Accessibility Refinement [Implemented]

**Objective:** Polish the user experience and make accessibility a
first-class, verifiable property.

**Problems addressed:**

- Empty, loading, and error states are now standardized across all routes.
- Accessibility documentation is comprehensive; axe-core runs in CI.
- Contextual help is provided through tooltips on form fields and attributes.
- Confirmation dialogs use a standardized pattern for destructive actions.

**Features/improvements:**

- [x] Document accessibility features and guidance in the UI guide (keyboard
      navigation, screen reader support, visual accessibility, forms, testing,
      known limitations, reporting process).
- [x] Create standardized UI state components (`EmptyState`, `LoadingSpinner`)
      with full accessibility support and comprehensive component tests.
- [x] Refactor all data-displaying routes to use standardized components:
      machines, users, audit, derp, topology, auth-keys all use `EmptyState`
      with appropriate variants (`default`, `filtered`, `error`).
- [x] Extend e2e accessibility test coverage (empty/loading states, icon-only
      buttons).
- [x] Confirmation dialogs standardized: 25 dialogs use the shared `Dialog` +
      `DialogPanel` pattern with `destructive`, `normal`, and `unactionable`
      variants.
- [x] Contextual help provided through tooltip component on 12+ UI elements
      (form fields, attribute labels, chips, status indicators).
- [ ] Extend the a11y suite beyond axe-core: keyboard navigation, focus
  management, color contrast, and screen-reader flows; publish a conformance
  statement (deferred to future phase).
- [ ] Add first-run onboarding flow (deferred to future phase).
- [ ] Keep the UI guide screenshots (from Phase 7) current (tracked separately).

**Technical work:** Shared state/empty-state components; a11y test expansion;
help/onboarding content.

**UI/UX work:** This phase is primarily UI/UX.

**Dependencies:** Phase 7 (component tests make UI changes safer).

**Expected outcome:** A more professional, accessible, and user-friendly UI.

**Priority:** High.

## Phase 9 — Supporting Features and Observability [Implemented]

**Objective:** Add supporting functionality and operational visibility.

**Status:** Completed in v0.34.0-arsydoni4326-alt.

**Completed work:**

- [x] Added `/ready` endpoint for Kubernetes-style readiness probes (returns 200 OK
      when ready, 503 Service Unavailable when database is unreachable).
- [x] Created comprehensive observability documentation (`docs/usage/observability.md`):
  - Documented all operational endpoints (`/health`, `/ready`, `/version`, `/api/v1/health`)
  - Cataloged all Prometheus metrics (HTTP, MapResponse, NodeStore, Mapper, Update Check, HA Health Probe)
  - Documented structured logging with zerolog (log levels, formats, best practices)
  - Documented debug endpoints (`/debug/overview`, `/debug/config`, `/debug/policy`, `/debug/pprof/`, `/debug/statsviz`)
  - Added Prometheus scrape configuration examples and recommended alerting rules
  - Added troubleshooting guide for common observability issues
- [x] Added unit tests for health endpoints (`TestHealthHandler`, `TestReadyHandler`, `TestVersionHandler`)
- [x] All tests pass successfully

**Deferred items (require separate approval):**

- In-app notifications for key events (machine expiry, pending approvals, update available)
- Admin-only, opt-in usage analytics and reporting
- Audit log export (CSV/JSON) and retention policy
- Request tracing and log correlation between Headplane and Headscale
- Search/filtering enhancements (fuzzy search, saved filters)

**Technical work:** Operational endpoints, comprehensive metrics documentation, test coverage.

**Dependencies:** Phase 7, 8.

**Expected outcome:** Better operational visibility and monitoring capabilities.

**Priority:** High.

## Phase 10 — Core Feature Expansion [Implemented]

**Objective:** Close feature gaps and expand product capabilities.

**Problems addressed:**

- OIDC groups cannot be used in ACLs (documented limitation in
  `docs/about/features.md`).
- No user self-service; all user management is admin-driven (deffered - need further consideration).
- Machine management parity gaps (route management, key rotation, device
  posture) — already planned in the Headplane roadmap.
- DNS management improvements (`extra_records_path` provisioning) — already
  planned in the Headplane roadmap.
- Upstream Tailscale features not implemented: Funnel, Serve, network flow
  logs.

**Features/improvements:**

- [x] OIDC group support in ACLs (backend) [Implemented].
- [ ] OIDC group support in ACLs (frontend) [Deferred — follow-up within Phase 10].
- User self-service: registration, password reset, profile management (deffered - need further consideration) [Proposed].
- [x] Machine management parity — key rotation (frontend) [Implemented].
- [x] Machine management parity — route overview page (frontend) [Implemented].
- [ ] Machine management parity — device posture [Deferred — requires significant policy engine changes].
- [x] DNS management improvements — `extra_records_path` visibility, CNAME
  records, and IP validation (frontend) [Implemented].
- [x] Extensibility foundation — documented extension points
  (`docs/ref/extending.md`), added Grafana dashboard example, validated
  Terraform/K8s operator support [Implemented].
- Evaluate Funnel / Serve / network flow logs; implement only if stable
  upstream and justified by demand [Proposed].

**Technical work:** Backend policy/API changes for OIDC groups; self-service
flows; parity features.

**UI/UX work:** Self-service pages, parity UI, DNS improvements.

**Dependencies:** Phase 7-9.

**Expected outcome:** Parity with more Tailscale features and improved user
management.

**Priority:** Medium.

## Phase 12 — Web-based Machine Approval [Implemented]

**Objective:** Enable administrators to approve pending machines through the Headplane
web interface, in addition to existing CLI and API methods.

**Status:** Implemented in v0.30.0-arsydoni4326-alt.

**Problems addressed:**

- Manual CLI-only approval workflow is cumbersome for administrators managing multiple machines.
- No visibility into pending machines without running CLI commands.
- Bulk operations require scripting around the CLI or API.

**Features/improvements:**

- [x] Backend API endpoints for machine approval [Implemented]:
  - `POST /api/v1/machines/{id}/approve` - Single machine approval
  - `POST /api/v1/machines/approve` - Bulk machine approval
- [x] Frontend UI for machine approval workflow [Implemented]:
  - View pending machines in Headplane
  - Approve individual machines with confirmation
  - Bulk selection and approval
  - Real-time status updates
  - Error handling and user feedback
- [x] Documentation [Implemented]:
  - API reference updated with endpoint schemas and examples
  - Registration workflow updated to mention web-based approval
  - UI guide with approval workflow and troubleshooting
  - Updated ROADMAP and CHANGELOG

**Technical work:** REST API endpoints, machine state management, permission checks,
frontend integration with existing machine management UI.

**UI/UX work:** Pending machines view, approval confirmation dialogs, bulk selection
interface, status badges, error messaging.

**Dependencies:** Phase 7-10 (API infrastructure, Headplane UI foundation).

**Expected outcome:** Streamlined machine approval workflow accessible through the
web interface, reducing administrator friction.

**Priority:** Medium.

## Phase 11 — Advanced Features and Integrations [Proposed / Future]

**Objective:** Prepare for future extensibility and large-scale use. These
items are deliberately deferred and should not be implemented until the
current tasks require them.

**Problems addressed:**

- No plugin/extension system for third-party UI components.
- No multi-Headscale-instance management from a single dashboard.
- Limited integrations with monitoring/alerting tooling.

**Features/improvements:**

- Plugin/extension system for Headplane [Proposed — deferred].
- Multi-Headscale-instance dashboard [Proposed — deferred].
- Monitoring/alerting integrations (Prometheus/Grafana webhooks) [Proposed].
- Terraform provider / Kubernetes operator support — the v2 API OAuth
  client-credentials flow already enables this; document and validate
  [Proposed].

**Technical work:** Depends on the chosen item.

**UI/UX work:** Depends on the chosen item.

**Dependencies:** Phase 7-10.

**Expected outcome:** Future-proofing and extensibility.

**Priority:** Low / Deferred.

## Phase 12 — Web-based Machine Approval [Partially implemented]

**Objective:** Let administrators approve pending machines from the Headplane
web UI, without the CLI, via single and bulk REST endpoints with RBAC and audit
logging.

**Problems addressed:**

- Approving a pending machine required shell access to the Headscale CLI
  (`headscale auth register` / route approval), which is impractical for
  web-only operators.

**Features/improvements:**

- [x] Backend REST endpoints `POST /api/v1/machines/{id}/approve` (single) and
      `POST /api/v1/machines/approve` (bulk) [Implemented].
- [x] RBAC: admin API key (all-access) or OAuth access token with the
      `devices:core` scope; read-only tokens are rejected with `403`
      [Implemented].
- [x] Audit logging for approvals (structured log entries; no backend audit
      table exists) [Implemented].
- [x] OpenAPI documentation, unit tests, and `docs/ref/api.md` coverage
      [Implemented].
- [ ] Headplane frontend: pending-machine list and one-click / bulk approve
      actions [Planned].

**Technical work:** v1 API handlers, v1 auth middleware RBAC extension,
`hscontrol/db/schema.sql` webhooks-table fix (pre-existing blocker).

**UI/UX work:** Headplane pending-approval page (pending).

**Dependencies:** Phase 4 (audit log in Headplane), Phase 6 (bulk operations).

**Expected outcome:** Pending machines can be approved end-to-end from the web
UI.

**Priority:** High.

## Phase 13 — Headplane Local Authentication and Settings

**Status:** Partially implemented. Phases 13a and 13b shipped under a
superseded Headscale-backed local-account design. Phase 13c replaces that
design with a single local Headplane administrator and a safe migration path.

**Objective:** Provide a local authentication system for Headplane that is
entirely independent from Headscale API authentication. After successful login,
users access a Settings menu where they can configure their Headscale API key
and customize their Headplane experience.

**Important note:** This phase was originally titled "Simple Password Login for
Headplane" and was partially implemented (commits 8e2131bf-966978f9) with a
simplified single-password approach. A later multi-user implementation moved
local accounts into the Headscale database. That model is now superseded: this
fork supports one local Headplane administrator, while Headscale API keys remain
the secondary authentication method.

---

### Phase 13a — Simple Password Authentication [Implemented]

**Status:** ✅ Complete (merged in v0.35.9-arsydoni4326-alt)

**Implemented features:**

- [x] Single password authentication for Headplane (no multi-user).
- [x] Password configured via `headplane.password` in config.yaml or
      `HEADSCALE_HEADPLANE_PASSWORD` environment variable.
- [x] Login endpoint: `POST /api/v1/headplane/login`.
- [x] Session management with 24-hour expiry and secure cookies.
- [x] Rate limiting (5 login attempts per minute per IP).
- [x] Constant-time password comparison to prevent timing attacks.
- [x] Password/API key toggle in login UI (password is default).
- [x] Full backward compatibility with API key authentication.
- [x] Comprehensive documentation in `docs/usage/authentication.md`.

**Files changed:**

- Backend: `hscontrol/headplane_auth.go`, `hscontrol/headplane_auth_test.go`,
  `hscontrol/types/config.go`, `hscontrol/app.go`
- Frontend: `headplane/app/routes/auth/login/*`, `headplane/app/server/web/auth.ts`
- Docs: `docs/usage/authentication.md`, `config-example.yaml`

**Limitations:**

- Single shared password (no user accounts).
- No settings menu or user preferences.
- API key must still be entered at login (not stored).

---

### Phase 13b — Settings Menu (Single User) [Implemented]

**Status:** ✅ Implemented (merged 2026-10-03).

**Objective:** Add a Settings menu for the authenticated user to manage their
Headscale API key, change password, and customize preferences. This phase does
NOT add multi-user support — it enhances the single-user experience from 13a.

**Requirements:**

- [x] Settings page accessible after password authentication.
- [x] **API Key Management**: Store and update the Headscale API key used by
      Headplane for API calls (eliminating need to re-enter at each login).
- [x] **Change Password**: Update the Headplane login password with current
      password verification.
- [x] **Theme Selection**: Choose between dark/light/system themes, persisted across
      sessions.
- [x] **Profile Name**: Optional display name for the authenticated user.
- [x] **Session Info**: Display current session expiry and logout button.

**Technical work:**

**Backend:**

- [x] Add settings storage (SQLite table for API key, theme, profile name).
- [x] Implement `POST /api/v1/headplane/settings` — update settings.
- [x] Implement `GET /api/v1/headplane/settings` — retrieve settings.
- [x] Implement `POST /api/v1/headplane/change-password` — change password with
      current password verification.
- [x] Store Headscale API key encrypted at rest (AES-256-GCM with PBKDF2-derived key).
- [x] Update session validation to load stored API key automatically.

**Frontend:**

- [x] Create `/settings/profile` route with sections:
  - Account (change password, session info, logout)
  - Integration (Headscale API key input with save/update)
  - Preferences (theme selector: dark/light/system)
  - Profile (display name, optional)
- [x] Implement theme persistence and application.
- [x] Add "Settings" link to navigation bar after login.
- [x] Form validation and error handling for all settings operations.
- [x] Visual feedback for save/update operations.

**Implementation details:**

- Settings stored in single-row SQLite table (`headplane_settings`).
- API key encrypted with AES-256-GCM using PBKDF2-derived key (100,000 iterations).
- Per-key salt and nonce for encryption security.
- CSRF protection on all state-changing endpoints.
- Comprehensive unit and integration tests.
- Documentation: user guide, API reference, integration testing guide.

**Files changed:**

- Backend: `hscontrol/headplane_settings.go`, `hscontrol/headplane_settings_test.go`, `hscontrol/app.go`
- Frontend: `headplane/app/routes/settings/profile.tsx`, API client updates
- Docs: `docs/usage/settings.md`, `docs/ref/api/headplane-settings.md`

**Expected outcome:**

1. ✅ Users log in once with their password and configure their Headscale API key in
   Settings.
2. ✅ API key is stored securely and reused across sessions (no re-entry needed).
3. ✅ Users can change their password without editing config files.
4. ✅ Theme preference persists across sessions.
5. ✅ Enhanced single-user experience without the complexity of multi-user accounts.

---

### Phase 13c — Single Local Administrator Migration [Implemented]

**Status:** Complete (2026-10-06). This phase superseded the Headscale database-backed
multi-user local-auth implementation completed in October 2026.

**Objective:** Support exactly one local Headplane administrator with
username/password login as the primary method and Headscale API-key login as the
secondary method. Local administrator credentials must belong to Headplane, not
to Headscale's database or configuration.

**Implementation outcome:**

- ✅ Headplane supports one local administrator only; local account creation, role
  assignment, user deletion, and per-user preferences are not supported.
- ✅ Existing OIDC and proxy-auth code remains in the repository but is disabled by
  single-admin mode. No user-facing configuration to re-enable it is provided.
- ✅ Headscale users, nodes, policies, and Headscale API keys remain independent of
  the local Headplane administrator.
- ✅ Legacy `headplane_users` and `headplane_settings` database records are
  retained for rollback but are no longer used for local Headplane authentication.

**Configuration and authentication delivered:**

- ✅ Added required `user` section to Headplane's `config.yaml`:
  ```yaml
  user:
    username: admin
    password: "$2b$12$..." # bcrypt verification hash; never plaintext
  ```
- ✅ Validates non-empty administrator username and supported bcrypt hash at
  Headplane startup. Rejects plaintext passwords and malformed hashes.
- ✅ Uses bcrypt cost 12 to verify local passwords and hash new passwords.
- ✅ Password verification, rate limiting, session issuance, and session
  invalidation moved into Headplane; no calls to Headscale's local password endpoints.
- ✅ Headscale API-key login preserved as an independent secondary login method.
- ✅ Requires configured `headscale.api_key` for local-password sessions to
  access Headscale data, with clear configuration error when absent.
- ✅ OIDC and proxy authentication disabled in this mode without deleting their
  code or unrelated persistent identity data.

**Administration UI delivered:**

- ✅ `/admin/users` route preserved but replaced user CRUD with a
  single-administrator Administration page.
- ✅ UI clearly states that only one local administrator is supported.
- ✅ Password reset/change form requires current password, confirmation, atomic
  config update, bcrypt rehash, and invalidation of all existing password sessions.
- ✅ When configuration source is immutable or externally managed, UI reset flow
  is disabled with precise operator recovery instructions.
- ✅ API-key list, creation, rotation, expiry/revocation, and deletion
  controls use Headscale's existing API-key endpoints.
- ✅ API-key secrets masked by default; newly created key displayed only once.
- ✅ Direct revocation or deletion of the configured Headplane API key prevented.
  Rotation creates and verifies replacement, atomically updates credential, then
  offers revocation of prior key.
- ✅ API-key-authenticated browser session ends immediately when its own key
  is revoked.

**SQLite migration tool delivered:**

- ✅ Headplane command `headplane migrate-local-admin --config <headplane-config> --legacy-db <headscale-sqlite-db>`
  shipped for native, Docker, and Nix deployments.
- ✅ Supports `--dry-run` and `--username`; requires `--username` when more than
  one legacy `admin` row exists.
- ✅ Reads legacy SQLite `headplane_users` table without modifying it.
- ✅ Copies exactly one selected administrator's existing valid bcrypt hash into
  `user.password`; passwords never recoverable or rewritten as plaintext.
- ✅ Does not migrate secondary accounts, roles, themes, profile names, or legacy
  per-user API-key ciphertext. Operators configure `headscale.api_key` separately.
- ✅ Creates timestamped, mode-`0600` Headplane-config backup; writes and
  `fsync`s same-directory temporary file; atomically renames into place;
  preserves restrictive ownership and permissions.
- ✅ Idempotent: matching config is successful no-op; config with different
  local-admin credentials fails without overwrite; interrupted writes leave
  original config valid; repeated runs never alter legacy database.
- ✅ Host-level `hash-password --password-stdin` and
  `reset-local-admin-password --password-stdin` commands provided for fresh setup and
  lockout recovery without exposing plaintext passwords in process arguments.

**Rollback and operational safety delivered:**

- ✅ Documented pre-migration backup of Headplane config and complete
  Headscale SQLite set, including `-wal` and `-shm` files when present.
- ✅ Documented stopping Headplane before migration and using `--dry-run` before
  actual write.
- ✅ Migration failures before atomic rename allow operators to correct reported
  precondition and rerun; config remains unchanged.
- ✅ Completed migration can be abandoned: operators stop Headplane, restore
  timestamped config backup, redeploy prior Headplane/Headscale versions, clear
  browser session cookie, and restart.
- ✅ No destructive Headscale database migration or legacy Headplane table drops.
  Any later cleanup requires separately announced migration after rollback window.

**Security and documentation delivered:**

- ✅ Documented that config encryption is not substitute for secret management:
  decryption key would still need to be available at startup.
- ✅ Requires `0600` configuration permissions, trusted owner, and parent
  directory not writable by untrusted users; never commit config files to source control.
- ✅ Recommends `headscale.api_key_path` backed by systemd credentials, Docker
  secrets, Kubernetes Secrets, or equivalent managed secret file.
- ✅ Documents HTTPS, secure cookies, API-key rotation, recovery, migration,
  rollback, and one-local-admin limitation in Headscale and Headplane docs.

**Verification completed:**

- ✅ Unit tests cover bcrypt config validation, local login, rate limiting,
  password reset, session invalidation, disabled OIDC/proxy paths, API-key
  lifecycle safeguards, and every migration success/failure/idempotency path.
- ✅ Browser tests cover local login, API-key login, single-admin notice,
  password reset, API-key rotation/revocation, and absence of local user
  management controls.
- ✅ Focused Headscale tests confirm legacy Headplane local-auth routes no
  longer mounted, API-key endpoints continue unchanged, and no destructive
  database migration introduced.
- ✅ Headplane unit, typecheck, lint, E2E, and docs-build checks plus
  focused Headscale test and formatting/lint checks completed before release.

See the [CHANGELOG](./CHANGELOG.md) for detailed implementation changes.

---

### Phase 13d — Restore `/admin/admin/users` Editable User Profile UI [Planned]

**Status:** Planned

**Objective:** Restore the `/admin/admin/users` route in Headplane with a full editable user profile UI that allows the local administrator to update their username, password, name, and avatar. All changes must persist to Headplane's `config.yaml`.

**Problem addressed:**

During the Phase 13c implementation, the `/admin/admin/users` route was accidentally removed. This route should provide a dedicated editable user profile interface for the single local administrator.

**Requirements:**

- [ ] Restore `/admin/admin/users` route in Headplane.
- [ ] Provide an editable UI with the following fields:
  - [ ] Username (editable text field)
  - [ ] Password (secure password input with current password verification)
  - [ ] Name/Display Name (editable text field)
  - [ ] Avatar Picture (optional file upload or URL input)
  - [ ] Save button to persist all changes
- [ ] All submitted field updates must persist to the corresponding values in Headplane's `config.yaml`:
  - `user.username`
  - `user.password` (stored as bcrypt hash)
  - `user.name`
  - `user.avatar` (if implemented)
- [ ] **Graceful handling of missing configuration fields:**
  - If any user fields are absent or undefined in `config.yaml`, the UI must display placeholder values and allow editing.
  - The UI must not raise errors or fail to load when fields are missing.
  - Empty or missing fields should render as empty input fields with appropriate placeholders (e.g., "Enter username", "No name set").
- [ ] Form validation and error handling:
  - Username uniqueness validation (if applicable)
  - Password strength requirements
  - Current password verification before allowing password changes
  - File size/type validation for avatar uploads
- [ ] Visual feedback for successful save operations and error states.
- [ ] Atomic configuration updates with backup creation (following Phase 13c's config update pattern).

**Technical work:**

- Backend: Extend or create Headplane configuration update endpoints that support username, name, and avatar field updates alongside existing password change functionality.
- Frontend: Implement `/admin/admin/users` route with a comprehensive user profile edit form.
- Configuration: Ensure `config.yaml` schema supports optional `user.name` and `user.avatar` fields.
- Testing: Unit and integration tests covering all field update scenarios, including missing field handling.

**Expected outcome:**

1. The local administrator can navigate to `/admin/admin/users` and edit their profile information.
2. All changes persist to `config.yaml` securely and atomically.
3. The UI gracefully handles partial or missing configuration without errors.
4. Improved user experience for managing local administrator account details.

**Priority:** Medium

---

## Phase 13 Architecture Note

The separation of concerns across all sub-phases:

- **Headplane local administrator** (Phase 13c): One bcrypt-backed username and
  password configured in Headplane's own `config.yaml`; Headplane verifies the
  password and owns its local browser session lifecycle.
- **Headscale API key**: A distinct Headscale credential used for dashboard API
  calls and available as the secondary Headplane login method. It is configured
  through Headplane, preferably with `headscale.api_key_path` and managed by a
  deployment secret mechanism.
- **Legacy local-account data**: The prior Headscale `headplane_users` and
  `headplane_settings` tables are migration input and rollback data only. They
  are not the runtime source of authentication after Phase 13c.

## Phase 15 — Database Version Migration v0.36.3 → v1.0.0 [In Progress]

**Status:** In Progress  
**Priority:** Critical (blocks v1.0.0 adoption)  
**Impact:** Database version metadata update (no schema changes)

### Problem Statement

Headscale v1.0.0-arsydoni4326 cannot open databases created with v0.36.3-arsydoni4326-alt due to major version check (v0→v1), even though:

- **No schema changes** between versions
- **Phase 13c retirement** is code-only (runtime user creation removed)
- **Legacy tables preserved** for rollback (`headplane_users`, `headplane_settings`)
- **Database structure identical** between v0.36.3 and v1.0.0

### Root Cause

The `checkVersionUpgradePath()` function in `hscontrol/db/db.go` blocks major version jumps (v0→v1) without considering that this fork's v1.0.0 is feature-equivalent to v0.36.3 with Phase 13c code retirement.

### Solution: Manual Version Migration

Since no schema changes exist, the migration updates database version metadata only:

#### Step 1: Backup Everything

```bash
# Stop Headscale
docker compose down

# Backup database files
cp /var/lib/headscale/headscale.db /var/lib/headscale/headscale.db.v0.36.3.backup
cp /var/lib/headscale/headscale.db-wal /var/lib/headscale/headscale.db-wal.backup 2>/dev/null || true
cp /var/lib/headscale/headscale.db-shm /var/lib/headscale/headscale.db-shm.backup 2>/dev/null || true

# Backup entire directory
tar -czf /tmp/headscale-backup-$(date +%Y%m%d-%H%M%S).tar.gz /var/lib/headscale/
```

#### Step 2: Update Database Version Metadata

Create migration script `/tmp/migrate-v0-to-v1.sql`:

```sql
-- Phase 15: Update database version metadata for v1.0.0 migration
-- NO SCHEMA CHANGES - version metadata only

BEGIN TRANSACTION;

-- Update the last_seen_version if it exists
UPDATE kv 
SET value = '1.0.0-arsydoni4326' 
WHERE key = 'last_seen_version';

-- If no version record exists, insert it
INSERT OR IGNORE INTO kv (key, value) 
VALUES ('last_seen_version', '1.0.0-arsydoni4326');

-- Verify the migration
SELECT key, value FROM kv WHERE key = 'last_seen_version';

COMMIT;
```

Apply the migration:

```bash
# Apply SQL migration
sqlite3 /var/lib/headscale/headscale.db < /tmp/migrate-v0-to-v1.sql

# Verify version updated
sqlite3 /var/lib/headscale/headscale.db "SELECT key, value FROM kv WHERE key = 'last_seen_version';"
```

#### Step 3: Verify Phase 13c Configuration

Ensure Headplane is configured for Phase 13c (single local admin):

```yaml
# /etc/headplane/config.yaml or equivalent
user:
  username: admin
  password: "$2b$12$..."  # bcrypt hash, NOT plaintext

headscale:
  url: "http://headscale:8080"
  api_key: "your-admin-api-key"  # or api_key_path
```

If not configured, run Headplane migration:

```bash
headplane migrate-local-admin \
  --config /etc/headplane/config.yaml \
  --legacy-db /var/lib/headscale/headscale.db \
  --dry-run

# After reviewing, run actual migration
headplane migrate-local-admin \
  --config /etc/headplane/config.yaml \
  --legacy-db /var/lib/headscale/headscale.db
```

#### Step 4: Start v1.0.0

```bash
# Switch to v1.0.0
cd /home/denny/Project/headscale-project/headscale
git checkout v1.0.0-arsydoni4326

# Rebuild
make clean
make build

# Or for Docker
docker compose build --no-cache

# Start services
docker compose up -d

# Verify version
docker compose exec headscale headscale version
# Should show: v1.0.0-arsydoni4326

# Check logs for successful startup
docker compose logs headscale | head -50
```

### Rollback Procedure

If v1.0.0 fails:

```bash
# Stop services
docker compose down

# Restore backup
cp /var/lib/headscale/headscale.db.v0.36.3.backup /var/lib/headscale/headscale.db

# Revert to v0.36.3
git checkout v0.36.3-arsydoni4326-alt
make build
# Or: docker compose build --no-cache

# Restart
docker compose up -d
```

### What v1.0.0 Changes (Code Only)

- ✅ **Removed:** Automatic default admin user creation from `cfg.Headplane.Password`
- ✅ **Removed:** Tests for database-backed Headplane user creation
- ✅ **Preserved:** Legacy `headplane_users` and `headplane_settings` tables
- ✅ **Preserved:** All Headscale core functionality (nodes, users, routes, policies)
- ✅ **Required:** Headplane must use Phase 13c config-based authentication

### Verification Checklist

After migration to v1.0.0:

- [ ] Headscale starts without "version check" error
- [ ] `headscale version` reports v1.0.0-arsydoni4326
- [ ] `headscale nodes list` shows all existing nodes
- [ ] `headscale users list` shows all existing users
- [ ] `headscale routes list` shows all existing routes
- [ ] Headplane login works (config-based authentication)
- [ ] Headplane `/admin` page accessible
- [ ] API keys work for Headscale API access
- [ ] No data loss compared to v0.36.3 backup

### Future Permanent Fix (Phase 15b)

Update `hscontrol/db/db.go` to allow v0.36.3 → v1.0.0 migration:

```go
// In checkVersionUpgradePath()
// Allow fork's v0.36.3-arsydoni4326-alt → v1.0.0-arsydoni4326 migration
// This is safe because v1.0.0 has no schema changes from v0.36.3
if lastVersion == "0.36.3-arsydoni4326-alt" && 
   currentVersion == "1.0.0-arsydoni4326" {
    return nil  // Allow this specific migration
}
```

This permanent fix will be included in v1.0.1-arsydoni4326.

---

## Phase 14 — Module Path Rewrite to `github.com/arsydoni4326-alt/headscale` [Planned]

**Status:** Planned  
**Priority:** High  
**Impact:** Breaking change for internal imports and build tooling

### Objective

Rewrite the Go module path from `github.com/juanfont/headscale` to
`github.com/arsydoni4326-alt/headscale` to properly reflect the fork's
independent identity and fix version reporting issues in built binaries.

### Motivation

Currently, the Go module path still references the upstream repository, which
causes:

1. **Version reporting issues**: Build-time ldflags with
   `-X 'github.com/arsydoni4326-alt/headscale/hscontrol/types.Version=...'`
   fail to inject version information because the actual module path is
   `github.com/juanfont/headscale`.
2. **Import confusion**: Internal imports reference the upstream path, making
   the fork's independence less clear.
3. **Dependency management**: The module path doesn't match the repository URL,
   which can cause issues with Go tooling and downstream consumers.

### Implementation Plan

#### Part 1: Preparation and Analysis [Planned]

- [ ] Audit all Go source files for import statements referencing
      `github.com/juanfont/headscale`
- [ ] Identify all non-Go files that reference the module path:
  - [ ] `Makefile` and build scripts
  - [ ] `Dockerfile*` files (especially ldflags in build commands)
  - [ ] CI/CD configuration files (`.github/workflows/*.yml`, `nix/` files)
  - [ ] Documentation files (`docs/`, `README.md`, `CONTRIBUTING.md`)
  - [ ] Test fixtures and integration test configurations
  - [ ] `flake.nix` and Nix-related files
  - [ ] `tools/bump` and other tooling scripts
- [ ] Document the current module dependency graph to identify potential
      downstream impact
- [ ] Create a backup branch before starting the rewrite

#### Part 2: Core Module Path Update [Planned]

- [ ] Update `go.mod` module declaration to
      `github.com/arsydoni4326-alt/headscale`
- [ ] Run automated find-and-replace for all Go import statements:
  ```bash
  find . -name "*.go" -type f -exec sed -i \
    's|github.com/juanfont/headscale|github.com/arsydoni4326-alt/headscale|g' {} +
  ```
- [ ] Update `go.sum` by running `go mod tidy`
- [ ] Verify no references to the old path remain in Go files:
  ```bash
  grep -r "github.com/juanfont/headscale" --include="*.go"
  ```

#### Part 3: Build System and Tooling [Planned]

- [ ] Update `Makefile` ldflags to use the new module path
- [ ] Update all `Dockerfile*` files:
  - [ ] Fix ldflags in `Dockerfile`
  - [ ] Fix ldflags in `Dockerfile.debug`
  - [ ] Fix ldflags in `Dockerfile.tailscale-HEAD`
  - [ ] Fix ldflags in `Dockerfile.derper`
- [ ] Update `cmd/headscale/headscale.go` if it contains version information
- [ ] Update `tools/bump` and related tooling scripts
- [ ] Update `flakehashes.json` by running `go run ./cmd/vendorhash update`
- [ ] Update Nix files (`flake.nix`, `nix/*.nix`) if they reference the module
      path

#### Part 4: CI/CD and Testing [Planned]

- [ ] Update GitHub Actions workflows (`.github/workflows/*.yml`)
- [ ] Update integration test configurations (`integration/`, `cmd/hi/`)
- [ ] Update Nix flake checks and builders
- [ ] Run local build to verify:
  ```bash
  make clean
  make build
  ./headscale version  # should show correct version/commit
  ```
- [ ] Run unit tests: `make test`
- [ ] Run integration tests: `go run ./cmd/hi doctor` and sample tests
- [ ] Verify Nix build: `nix build`

#### Part 5: Documentation and Communication [Planned]

- [ ] Update `README.md` with the new module path for go get instructions
- [ ] Update `CONTRIBUTING.md` with the new import path conventions
- [ ] Update `docs/` documentation referencing the module path
- [ ] Update `ARCHITECTURE.md` if it references the module structure
- [ ] Update `SPECIFICATION.md` if it references the module path
- [ ] Add a migration note to `CHANGELOG.md` under Breaking Changes
- [ ] Update Headplane documentation if it references the backend module path

#### Part 6: Verification and Release [Planned]

- [ ] Build and test all binaries:
  - [ ] `headscale version` reports correct version and commit
  - [ ] Docker images build successfully
  - [ ] Nix builds complete without errors
- [ ] Run full integration test suite: `go run ./cmd/hi run --all`
- [ ] Verify the update-check endpoint still works
- [ ] Test Headplane integration with the updated backend
- [ ] Create a release tag following the `-arsydoni4326-alt` convention
- [ ] Communicate the change to any downstream consumers or collaborators

### Risk Assessment

- **High impact**: All Go imports change; automated tools required.
- **Breaking change**: Downstream projects importing this fork as a library
  will need to update their imports.
- **Build tooling**: Requires careful coordination with Dockerfiles, Makefiles,
  and CI/CD.
- **Testing burden**: Must verify all build paths (native, Docker, Nix) and
  all test suites.

### Success Criteria

1. ✅ `headscale version` command shows the correct version, commit, and build
   date (not "version=dev commit=unknown")
2. ✅ All unit tests pass
3. ✅ All integration tests pass
4. ✅ Docker builds complete successfully
5. ✅ Nix builds complete successfully
6. ✅ No grep matches for `github.com/juanfont/headscale` in Go files
7. ✅ Documentation accurately reflects the new module path
8. ✅ Update-check feature continues to work
9. ✅ Headplane integration remains functional

### Rollback Plan

If critical issues arise:

1. Revert to the backup branch created in Part 1
2. Document the specific failure mode
3. Re-plan with lessons learned

### Dependencies

- **Blocks**: None (this is an independent infrastructure change)
- **Blocked by**: None
- **Related**: This fixes the version reporting issue where ldflags fail to
  inject version information due to module path mismatch

### Estimated Timeline

- Part 1 (Preparation): 1-2 hours
- Part 2 (Core Update): 30 minutes
- Part 3 (Build System): 1-2 hours
- Part 4 (CI/CD and Testing): 2-3 hours
- Part 5 (Documentation): 1-2 hours
- Part 6 (Verification): 2-3 hours

**Total estimated effort**: 8-13 hours

---

## Tracking

- Day-to-day work is tracked via GitHub issues on the fork repositories.
- Long-horizon items live here. When an item moves from Planned to In Progress,
  update this document in the same change as the implementation.
- Completed items are removed from this document and recorded in the
  [CHANGELOG](./CHANGELOG.md) instead.

### Headplane roadmap sync

The Headplane frontend maintains its own roadmap at
[`headplane/docs/ROADMAP.md`](./headplane/docs/ROADMAP.md). Items that span
both backend and frontend (e.g. new API endpoints with UI pages) must be
tracked in **both** roadmaps. When an item moves between Planned and In
Progress:

1. Update this roadmap (Headscale) in the implementation commit.
2. Update `headplane/docs/ROADMAP.md` in the corresponding Headplane commit.
3. If the item spans a single release that includes both repos, the release
   notes should reference both roadmap updates.

The headplane roadmap follows the same lifecycle: Planned → In Progress →
removed (recorded in its CHANGELOG). Keep the status labels consistent across
both documents.
