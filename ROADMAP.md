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

- Accessibility audit against WCAG 2.1 AA: keyboard navigation, visible focus
  states, screen-reader labels, and color contrast across both themes.
- Responsive design pass for tablet and mobile widths (machine tables, ACL
  editor, dialogs).
- Consistent empty states, loading skeletons, and error states across all
  routes.
- Toast/notification system for action feedback (rename, expire, delete,
  ACL apply) instead of relying on inline banners alone.
- Rich machine detail view: per-machine page showing routes, tags, expiry,
  OS/version, and recent activity, instead of only the list row.
- Theme polish: verify both light and dark themes render every component
  correctly and consistently.

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