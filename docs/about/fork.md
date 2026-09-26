# The `arsydoni4326-alt` fork

This repository is a fork of [Headscale](https://headscale.net) maintained by
`arsydoni4326-alt`. It tracks upstream Headscale closely while adding a small
set of fork-specific features, most notably the update checker that spans both
the Go backend and the [Headplane](../ref/integration/web-ui.md) web UI.

## Version suffix convention

All version numbers in this fork end with the suffix `-arsydoni4326-alt`. For
example, a release version is written as `v0.29.8-arsydoni4326-alt`. This
suffix identifies releases of the fork and must always be present in version
numbers, tags, and changelog entries.

## Update check

The fork ships an update checker that compares the running binary's version
against the latest release on the fork's GitHub repository.

### Backend endpoint

`GET /api/v1/update-check` returns the running binary's version information:

- Without `?check=true`: returns current version info only (fast, no external
  call).
- With `?check=true`: fetches the latest release tag from the fork's GitHub
  repository and compares it with the running binary's version using semver
  precedence. For dev builds (dirty, pseudo-version, or untagged), it falls
  back to comparing short commit hashes. Results are cached for 15 minutes.

The endpoint is intentionally unauthenticated. The remote repository can be
overridden via the `HEADSCALE_UPDATE_CHECK_REPO` environment variable (format:
`owner/repo`). See the [API reference](../ref/api.md#update-check) for usage
examples.

### Headplane UI

When the [Headplane](../ref/integration/web-ui.md) web UI is deployed, it checks
for updates on page load and shows a modal when a new version is available. A
"Check for updates" button in the header menu triggers an immediate check.
Results are cached in `sessionStorage` for 15 minutes. The modal includes
"Dismiss for this session" and "Remind me later (24h)" options, and a direct
link to the release notes when a tag-based update is detected.

## Feature preservation

All fork-specific features are load-bearing and must survive upstream merges.
They are protected by a CI check that fails if any of the following are
removed:

- `hscontrol/updatecheck/` (backend update-check package)
- `headplane/app/update-check/` (frontend update-check domain)