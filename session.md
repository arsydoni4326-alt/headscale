# Session

## Objective

Add update check API endpoint to headscale backend.

## Progress

### Headplane (frontend)
- [x] Created `app/update-check/` domain (README, types, provider, modal, hook)
- [x] Modified build pipeline: `__COMMIT_HASH__` global, `HEADPLANE_COMMIT` env
- [x] Integrated: root.tsx provider, header menu item, modal on auto-check
- [x] Updated docs: ROADMAP.md, SPECIFICATION.md (FR-10)

### Headscale (backend)
- [x] Created `hscontrol/updatecheck/` package with:
  - `Handler()` - HTTP handler returning version info + remote check
  - `GET /api/v1/update-check` - returns current version
  - `GET /api/v1/update-check?check=true` - additionally fetches remote commit
  - Built-in GitHub API fetch against `github.com/arsydoni4326-alt/headscale`
- [x] Wired route in `app.go`
- [x] Added test file

## Decisions and Assumptions

- Headscale API endpoint requires no auth (like `/version` and `/health`)
- Uses the fork repo: `https://github.com/arsydoni4326-alt/headscale.git`
- The remote check optionally performs a GitHub API call when `?check=true`
- Both packages (frontend `app/update-check/` and backend `hscontrol/updatecheck/`) are self-contained to survive upstream merges

---

# Previous Session

## Objective

Git Flow release `v0.29.6-arsydoni4326-alt` — completed and merged to `main`.

## Progress

- [x] Committed staged `.github/workflows/deploy.yaml` change on `dev`
      (`5fb5e622 ci: use Dockerfile for GHCR deploy build`) — the GHCR deploy
      build now uses the repo `Dockerfile` instead of `Dockerfile.tailscale-HEAD`
- [x] Created release branch `release/v0.29.6-arsydoni4326-alt` from `dev`
- [x] Updated `CHANGELOG.md`: added the `0.29.6-arsydoni4326-alt (2026-09-25)`
      fork release section with the `Dockerfile` change
      (`52ae3658 changelog: add 0.29.6-arsydoni4326-alt fork release entry`)
- [x] Finished release with `git flow release finish`:
  - Merged release branch into `main` (merge commit `28273878`)
  - Tagged `main` with `v0.29.6-arsydoni4326-alt`
  - Merged release tag back into `dev` (merge commit `b059e497`)
  - Deleted `release/v0.29.6-arsydoni4326-alt` branch
- [ ] Push `main`, `dev`, and tag `v0.29.6-arsydoni4326-alt` to `origin` — **blocked: no GitHub credentials available in this environment**

## Decisions and Assumptions

- This repo uses `dev` as the integration branch (there is no `develop` branch).
- Version format convention: all version numbers in this fork must end with
  `-arsydoni4326-alt` (documented in `CONTRIBUTING.md`).
- Version number auto-increments from the latest tag: latest was
  `v0.29.5-arsydoni4326-alt`, so this release is `v0.29.6-arsydoni4326-alt`.
- All staged files are committed and included in the current release.
- `git flow release start` refuses to run with staged changes, so the staged
  file is committed on `dev` first, then the release branch is created from
  `dev` (same pattern as the v0.29.4 release).
- The v0.29.5 release did not add a CHANGELOG entry; this release follows the
  fork convention (suffix must appear in changelog entries) and adds one.

## Known Issues and Limitations

- Pushing to `origin` requires the user to authenticate (e.g., `gh auth login`,
  a PAT, or an authorized SSH key).

## Pending Work

- Push `main`, `dev`, and tag `v0.29.6-arsydoni4326-alt` to `origin` once
  credentials are available.
