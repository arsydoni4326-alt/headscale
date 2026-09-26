# Specification

This document describes the functional and non-functional requirements for the
`arsydoni4326-alt` fork of Headscale and its frontend, Headplane. It is the
reference for what the project must do, its constraints, and the acceptance
criteria used to evaluate changes.

> For how the system is structured internally, see
> [Architecture](./ARCHITECTURE.md). For the high-level direction, see
> [Roadmap](./ROADMAP.md). For past changes, see [Changelog](./CHANGELOG.md).
> Headplane's own specification and architecture live in
> [`headplane/docs/`](./headplane/docs/).

## Product Overview

Headscale is an open source, self-hosted implementation of the Tailscale
control server. It manages node registration, IP allocation, policy
enforcement, and DERP routing for self-hosted tailnets. Headplane is a
self-hosted web UI for Headscale.

This fork tracks upstream Headscale closely while adding a small set of
fork-specific features, most notably the update checker that spans both the Go
backend (`hscontrol/updatecheck/`) and the Headplane frontend
(`headplane/app/update-check/`).

## Feature Preservation

All existing features and customizations in this repository are load-bearing
and must be preserved. This includes, but is not limited to:

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

Every change must be **additive and non-destructive**. No change may remove,
rename, or break an existing feature without explicit review and a documented
migration path.

## Functional Requirements

### FR-1: Update check

- **FR-1.1**: `GET /api/v1/update-check` returns the running binary's version
  information (version, commit, build time, dirty flag).
- **FR-1.2**: With `?check=true`, the endpoint fetches the latest release tag
  from the remote GitHub repository and compares it with the running binary's
  version using semver precedence. For dev builds (dirty, pseudo-version, or
  untagged), it falls back to comparing short commit hashes.
- **FR-1.3**: The endpoint is intentionally unauthenticated.
- **FR-1.4**: Results are cached server-side for 15 minutes to stay within
  GitHub's unauthenticated rate limits (60 req/hour).
- **FR-1.5**: The remote repository is configurable via the
  `HEADSCALE_UPDATE_CHECK_REPO` environment variable (format: `owner/repo`),
  defaulting to `arsydoni4326-alt/headscale`.
- **FR-1.6**: The Headplane frontend checks for updates on page load and shows
  a modal when a new version is available. Results are cached in
  `sessionStorage` for 15 minutes. The modal includes "Dismiss for this
  session" and "Remind me later (24h)" options.
- **FR-1.7**: Prometheus metrics are exposed for the endpoint:
  `headscale_updatecheck_requests_total`,
  `headscale_updatecheck_remote_failures_total`,
  `headscale_updatecheck_cache_hits_total`,
  `headscale_updatecheck_cache_misses_total`.

### FR-2: Versioning

- **FR-2.1**: All version numbers end with the suffix `-arsydoni4326-alt`.

## Non-Functional Requirements

### NFR-1: Compatibility

- The fork tracks upstream Headscale closely; upstream features must keep
  working.
- No breaking changes to configuration, behavior, or API without an explicit
  migration path and prior discussion.

### NFR-2: Testing

- Backend changes require unit tests and, when practical, integration tests.
- Frontend changes require unit tests and a clean typecheck.
- Fork-specific features are protected by a CI check that fails if they are
  removed.

### NFR-3: Documentation

- Features are only complete when their documentation is complete.
- The user-facing docs site (`docs/`) must reflect fork-specific features.

## Constraints

- **Stack**: Go for the backend; React Router 7 + Vite + TypeScript for the
  Headplane frontend.
- **Versioning**: semantic versioning with the `-arsydoni4326-alt` suffix.
- **Database**: SQLite for local development, PostgreSQL for
  integration-heavy tests.

## Acceptance Criteria

A change is considered complete when:

1. It implements the smallest correct change for the requirement and does not
   alter unrelated behavior.
2. Tests exist and pass.
3. The documentation reflects the new behavior.
4. No breaking change to configuration or API exists without a documented
   migration path.