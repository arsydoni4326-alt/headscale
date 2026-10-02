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

- Phases 1-10 are complete. Latest fork releases: Headscale
  `v0.34.0-arsydoni4326-alt`, Headplane `v0.8.3-arsydoni4326-alt`.
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
