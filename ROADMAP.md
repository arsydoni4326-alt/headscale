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

## Phase 1 — Foundation and Documentation

- Create root-level `SPECIFICATION.md` and `ARCHITECTURE.md` describing the fork
  as a whole (backend + frontend), including the fork-specific features and the
  feature-preservation rule above.
- Add a CI check that verifies `hscontrol/updatecheck/` and
  `headplane/app/update-check/` still exist after every merge, so an upstream
  merge cannot silently delete them.
- Document the fork-specific features in the user-facing docs site
  (`docs/`), including the update-check endpoint and the version suffix
  convention.
- Add the `/api/v1/update-check` endpoint to the OpenAPI specification so it is
  discoverable alongside the rest of the v1 API.

## Phase 2 — Update Checker Hardening (Completed)

See [CHANGELOG](./CHANGELOG.md) for the list of changes in this phase.

## Phase 3 — UI/UX Polish

### Part 1: Audit and Visual Documentation

- [x] UI Inventory & Audit: read all template files, catalog UI elements and
      shared components, identify inconsistencies and polish opportunities.
- [x] Visual Documentation Base:
      - Created `docs/usage/ui-guide.md` (placeholder screenshots, Mermaid flow
        diagrams, contributor guidance, maintenance checklist).
      - Created `docs/assets/screenshots/` and `docs/assets/diagrams/` directories.
      - Updated `mkdocs.yml` to include UI Guide in nav.
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

## Phase 4 — Feature Expansion

- Network topology visualization of the tailnet (nodes, routes, exit nodes).
- Audit log / activity feed showing who changed what (machine, ACL, DNS,
  settings).
- Bulk operations on machines: expire, delete, and tag multiple selected
  machines at once.
- DERP status and health page (regions, latency, relay usage).
- Search and filter improvements: filter machines by OS, user, tag, online
  status, and expiry; persist filters in the URL.
- Export/import of ACL policy and Headscale configuration for backup and
  migration.
- Headscale version compatibility tracking: proactively detect and surface
  unsupported features based on the server's reported version.

## Phase 5 — Testing, Performance, and CI/CD

- Frontend end-to-end tests (Playwright) for critical flows: login, machine
  management, ACL editing, DNS settings.
- Accessibility tests (axe-core) wired into CI.
- Lighthouse CI with performance budgets for the main routes.
- Bundle-size analysis and code splitting; lazy-load the WASM SSH payload so it
  is only fetched when the SSH page is opened.
- Backend tests for the update-check endpoint edge cases (see Phase 2).
- CI gate that fails on missing fork-specific features (see Phase 1).

## Phase 6 — Community and Ecosystem

- Issue templates for feature requests, bug reports, and UX feedback.
- Contribution guide updates covering the fork-specific features and the
  feature-preservation rule.
- Keep this roadmap and the headplane roadmap in sync when items move between
  Planned and In Progress.

## Tracking

- Day-to-day work is tracked via GitHub issues on the fork repositories.
- Long-horizon items live here. When an item moves from Planned to In Progress,
  update this document in the same change as the implementation.
- Completed items are removed from this document and recorded in the
  [CHANGELOG](./CHANGELOG.md) instead.