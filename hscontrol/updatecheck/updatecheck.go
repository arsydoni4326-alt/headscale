// Package updatecheck provides an indestructible, self-contained update check
// endpoint for the headscale fork at github.com/arsydoni4326-alt/headscale.
//
// DO NOT REMOVE: This is a fork-specific feature that must survive upstream
// merges. All update-check logic lives in this single package so it can be
// restored as a unit if removed by an upstream merge.
//
// Endpoint: GET /api/v1/update-check
//
// Without ?check=true: returns current version info only (fast, no external call).
// With    ?check=true: additionally fetches the latest commit hash from the remote
//                      fork's main branch via the GitHub API and compares it with
//                      the running binary's embedded commit hash.
package updatecheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/juanfont/headscale/hscontrol/types"
	"github.com/rs/zerolog/log"
)

const (
	// RemoteRepo is the fork repository URL for update checking.
	RemoteRepo = "https://github.com/arsydoni4326-alt/headscale.git"

	// RemoteRepoOwner is the GitHub owner of the fork.
	RemoteRepoOwner = "arsydoni4326-alt"

	// RemoteRepoName is the GitHub repository name of the fork.
	RemoteRepoName = "headscale"

	// githubAPIURL is the GitHub API endpoint for the latest commit on main.
	githubAPIURL = "https://api.github.com/repos/" + RemoteRepoOwner + "/" + RemoteRepoName + "/commits/main"

	// httpTimeout is the timeout for the GitHub API request.
	httpTimeout = 10 * time.Second
)

// CurrentVersionResponse is the version information for the running binary.
type CurrentVersionResponse struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
	Dirty     bool   `json:"dirty"`
}

// RemoteVersionResponse is the version information from the remote repository.
type RemoteVersionResponse struct {
	Commit string `json:"commit"`
	URL    string `json:"url"`
}

// UpdateCheckResponse is the full response from the update check endpoint.
type UpdateCheckResponse struct {
	Current         CurrentVersionResponse  `json:"current"`
	UpdateAvailable *bool                   `json:"updateAvailable,omitempty"`
	Remote          *RemoteVersionResponse  `json:"remote,omitempty"`
	Error           string                  `json:"error,omitempty"`
}

// Handler returns an http.HandlerFunc that serves the update check endpoint.
//
// It is mounted at /api/v1/update-check in the main router. The handler:
//  1. Always returns the running binary's version info.
//  2. When ?check=true is set, fetches the remote latest commit hash
//     and compares it with the embedded commit.
func Handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		versionInfo := types.GetVersionInfo()

		resp := UpdateCheckResponse{
			Current: CurrentVersionResponse{
				Version:   versionInfo.Version,
				Commit:    versionInfo.Commit,
				BuildTime: versionInfo.BuildTime,
				Dirty:     versionInfo.Dirty,
			},
		}

		// Only perform the remote check if ?check=true is in the query
		if strings.EqualFold(req.URL.Query().Get("check"), "true") {
			remoteCommit, err := fetchRemoteShortCommit()
			if err != nil {
				log.Warn().
					Err(err).
					Str("url", githubAPIURL).
					Msg("update-check: failed to fetch remote commit")

				resp.Error = fmt.Sprintf("failed to check remote: %s", err.Error())
			} else {
				updateAvailable := remoteCommit != "" && remoteCommit != shortCommit(versionInfo.Commit)
				resp.UpdateAvailable = &updateAvailable
				resp.Remote = &RemoteVersionResponse{
					Commit: remoteCommit,
					URL:    fmt.Sprintf("https://github.com/%s/%s", RemoteRepoOwner, RemoteRepoName),
				}
			}
		}

		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(resp)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Msg("update-check: failed to write response")
		}
	}
}

// fetchRemoteShortCommit fetches the short commit hash of the latest commit
// on the remote fork's main branch via the GitHub API.
func fetchRemoteShortCommit() (string, error) {
	client := &http.Client{
		Timeout: httpTimeout,
	}

	resp, err := client.Get(githubAPIURL) //nolint:noctx // simple GET with timeout is sufficient
	if err != nil {
		return "", fmt.Errorf("github api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var data struct {
		SHA string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode github api response: %w", err)
	}

	if data.SHA == "" {
		return "", fmt.Errorf("github api returned empty sha")
	}

	return shortCommit(data.SHA), nil
}

// shortCommit returns the first 7 characters of a commit hash, or the
// full string if it is already shorter than 7 characters.
func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

// BuildInfo returns the vcs.revision and vcs.time from the build info,
// falling back to the VersionInfo's commit if build info is unavailable.
// This is used for embedding additional build metadata.
func BuildInfo() (revision, time string) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown", "unknown"
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.time":
			time = setting.Value
		}
	}

	if revision == "" {
		revision = "unknown"
	}
	if time == "" {
		time = "unknown"
	}

	return revision, time
}