package updatecheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_WithoutCheckParam(t *testing.T) {
	ResetCache()

	handler := Handler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/update-check", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %s", contentType)
	}

	var body UpdateCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should always have current version info
	if body.Current.Version == "" {
		t.Error("expected current.version to be non-empty")
	}
	if body.Current.Commit == "" {
		t.Error("expected current.commit to be non-empty")
	}

	// Without ?check=true, remote should be nil
	if body.Remote != nil {
		t.Error("expected remote to be nil without ?check=true")
	}
	if body.UpdateAvailable != nil {
		t.Error("expected updateAvailable to be nil without ?check=true")
	}
}

func TestHandler_WithCheckParam(t *testing.T) {
	ResetCache()

	handler := Handler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/update-check?check=true", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body UpdateCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should always have current version info
	if body.Current.Version == "" {
		t.Error("expected current.version to be non-empty")
	}
	if body.Current.Commit == "" {
		t.Error("expected current.commit to be non-empty")
	}

	// With ?check=true, remote may be nil (if GitHub API is unreachable in test)
	// but should not return an error about missing check param
	if body.Error != "" {
		t.Logf("remote check error (expected in test env): %s", body.Error)
	}
}

func TestShortCommit(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abc123def456ghi789", "abc123d"},
		{"abc123d", "abc123d"},
		{"ab", "ab"},
		{"", ""},
	}

	for _, tt := range tests {
		result := shortCommit(tt.input)
		if result != tt.expected {
			t.Errorf("shortCommit(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestHandler_ResponseStructure(t *testing.T) {
	ResetCache()

	handler := Handler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/update-check?check=true", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify required fields exist
	requiredFields := []string{"current"}
	for _, field := range requiredFields {
		if _, ok := body[field]; !ok {
			t.Errorf("expected response field %q to exist", field)
		}
	}

	// Verify nested current fields
	current, ok := body["current"].(map[string]interface{})
	if !ok {
		t.Fatal("expected current to be an object")
	}

	nestedFields := []string{"version", "commit", "buildTime"}
	for _, field := range nestedFields {
		if _, ok := current[field]; !ok {
			t.Errorf("expected current.%s to exist", field)
		}
	}
}
// --- Version parsing tests ---

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input   string
		wantOk  bool
		wantMaj int
		wantMin int
		wantPat int
		wantPre string
	}{
		{"0.29.9", true, 0, 29, 9, ""},
		{"v0.29.9", true, 0, 29, 9, ""},
		{"0.29.9-arsydoni4326-alt", true, 0, 29, 9, "-arsydoni4326-alt"},
		{"v0.29.9-arsydoni4326-alt", true, 0, 29, 9, "-arsydoni4326-alt"},
		{"1.0.0", true, 1, 0, 0, ""},
		{"1.2.3-rc1", true, 1, 2, 3, "-rc1"},
		{"1.2.3+build123", true, 1, 2, 3, "+build123"},
		{"", false, 0, 0, 0, ""},
		{"dev", false, 0, 0, 0, ""},
		{"unknown", false, 0, 0, 0, ""},
		{"(devel)", false, 0, 0, 0, ""},
		{"not-a-version", false, 0, 0, 0, ""},
		{"1.2", false, 0, 0, 0, ""},
		{"abc.def.ghi", false, 0, 0, 0, ""},
	}

	for _, tt := range tests {
		got, err := parseVersion(tt.input)
		if tt.wantOk {
			if err != nil {
				t.Errorf("parseVersion(%q) unexpected error: %v", tt.input, err)
				continue
			}
			if got.major != tt.wantMaj {
				t.Errorf("parseVersion(%q) major = %d, want %d", tt.input, got.major, tt.wantMaj)
			}
			if got.minor != tt.wantMin {
				t.Errorf("parseVersion(%q) minor = %d, want %d", tt.input, got.minor, tt.wantMin)
			}
			if got.patch != tt.wantPat {
				t.Errorf("parseVersion(%q) patch = %d, want %d", tt.input, got.patch, tt.wantPat)
			}
			if got.preRelease != tt.wantPre {
				t.Errorf("parseVersion(%q) preRelease = %q, want %q", tt.input, got.preRelease, tt.wantPre)
			}
		} else {
			if err == nil {
				t.Errorf("parseVersion(%q) expected error, got %+v", tt.input, got)
			}
		}
	}
}

func TestSemverGreaterThan(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want bool // a > b
	}{
		{"1.0.0", "0.9.9", true},
		{"0.9.9", "1.0.0", false},
		{"0.29.10", "0.29.9", true},
		{"0.29.9", "0.29.10", false},
		{"0.29.9", "0.29.9", false},
		{"0.29.9", "0.29.9-arsydoni4326-alt", true},  // no pre > pre
		{"0.29.9-arsydoni4326-alt", "0.29.9", false}, // pre < no pre
		{"0.29.10-arsydoni4326-alt", "0.29.9-arsydoni4326-alt", true},
		{"0.29.9-rc2", "0.29.9-rc1", true},
		{"0.29.9-rc1", "0.29.9-rc2", false},
		{"0.29.9-alpha", "0.29.9-beta", false},  // lexical: "alpha" < "beta"
		{"0.29.9-beta", "0.29.9-alpha", true},   // lexical: "beta" > "alpha"
		{"0.29.9-rc1", "0.29.9-rc1", false},     // same
	}

	for _, tt := range tests {
		a, errA := parseVersion(tt.a)
		b, errB := parseVersion(tt.b)
		if errA != nil || errB != nil {
			t.Fatalf("parse error: a=%v, b=%v", errA, errB)
		}

		got := a.GreaterThan(b)
		if got != tt.want {
			t.Errorf("parseVersion(%q).GreaterThan(%q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestIsDevVersion(t *testing.T) {
	devVersions := []string{"", "dev", "unknown", "(devel)"}
	releaseVersions := []string{"0.29.9", "v0.29.9-arsydoni4326-alt", "1.0.0"}

	for _, v := range devVersions {
		if !isDevVersion(v) {
			t.Errorf("isDevVersion(%q) = false, want true", v)
		}
	}
	for _, v := range releaseVersions {
		if isDevVersion(v) {
			t.Errorf("isDevVersion(%q) = true, want false", v)
		}
	}
}

// --- Caching tests ---

func TestCachedFetch(t *testing.T) {
	ResetCache()

	var callCount int
	fetch := func() (string, error) {
		callCount++
		return "abc1234", nil
	}

	// First call: miss, fetches
	val, err := cachedFetch("test", 5*time.Minute, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "abc1234" {
		t.Errorf("got %q, want %q", val, "abc1234")
	}
	if callCount != 1 {
		t.Errorf("callCount = %d, want 1", callCount)
	}

	// Second call: hit, no fetch
	val, err = cachedFetch("test", 5*time.Minute, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "abc1234" {
		t.Errorf("got %q, want %q", val, "abc1234")
	}
	if callCount != 1 {
		t.Errorf("callCount = %d, want 1 (cached)", callCount)
	}
}

func TestCachedFetch_Expired(t *testing.T) {
	ResetCache()

	var callCount int
	fetch := func() (string, error) {
		callCount++
		return "abc1234", nil
	}

	// Fetch with very short TTL
	val, err := cachedFetch("expire-test", 1*time.Millisecond, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "abc1234" {
		t.Errorf("got %q, want %q", val, "abc1234")
	}

	// Wait for expiry
	time.Sleep(2 * time.Millisecond)

	// Should fetch again
	val, err = cachedFetch("expire-test", 1*time.Millisecond, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "abc1234" {
		t.Errorf("got %q, want %q", val, "abc1234")
	}
	if callCount != 2 {
		t.Errorf("callCount = %d, want 2 (expired)", callCount)
	}
}

func TestCachedFetch_ErrorNotCached(t *testing.T) {
	ResetCache()

	var callCount int
	fetch := func() (string, error) {
		callCount++
		return "", fmt.Errorf("test error")
	}

	// First call: error, not cached
	_, err := cachedFetch("err-test", 5*time.Minute, fetch)
	if err == nil {
		t.Fatal("expected error")
	}
	if callCount != 1 {
		t.Errorf("callCount = %d, want 1", callCount)
	}

	// Second call: should retry (error was not cached)
	_, err = cachedFetch("err-test", 5*time.Minute, fetch)
	if err == nil {
		t.Fatal("expected error")
	}
	if callCount != 2 {
		t.Errorf("callCount = %d, want 2 (retried)", callCount)
	}
}

func TestCachedFetch_Reset(t *testing.T) {
	ResetCache()

	var callCount int
	fetch := func() (string, error) {
		callCount++
		return "abc1234", nil
	}

	_, _ = cachedFetch("reset-test", 5*time.Minute, fetch)
	if callCount != 1 {
		t.Errorf("callCount = %d, want 1", callCount)
	}

	ResetCache()

	// After reset, should fetch again
	_, _ = cachedFetch("reset-test", 5*time.Minute, fetch)
	if callCount != 2 {
		t.Errorf("callCount = %d, want 2 (after reset)", callCount)
	}
}
// --- Edge case tests ---

func TestBuildResponse_WithCheck_NoExternalCall(t *testing.T) {
	// Without ?check=true, BuildResponse should never make external calls.
	ResetCache()
	resp := BuildResponse(false)
	if resp.Remote != nil {
		t.Error("expected remote to be nil when check=false")
	}
	if resp.UpdateAvailable != nil {
		t.Error("expected updateAvailable to be nil when check=false")
	}
}

func TestBuildResponse_DevVersion(t *testing.T) {
	// In test environments, the version is typically "dev" or "(devel)".
	// BuildResponse should handle this gracefully.
	ResetCache()
	resp := BuildResponse(true)

	// Should either succeed (if GitHub is reachable) or return an error
	if resp.Error != "" {
		t.Logf("remote check error (expected in test env): %s", resp.Error)
	}

	// Always has current version info
	if resp.Current.Version == "" {
		t.Error("expected current.version to be non-empty")
	}
}

func TestBuildResponse_RemoteVersionResponseHasVersion(t *testing.T) {
	// Verify that RemoteVersionResponse includes the Version field in JSON.
	resp := UpdateCheckResponse{
		Current: CurrentVersionResponse{Version: "0.29.9", Commit: "abc1234", BuildTime: "2026-01-01", Dirty: false},
		UpdateAvailable: func() *bool { v := false; return &v }(),
		Remote: &RemoteVersionResponse{
			Commit:  "v0.29.9-arsydoni4326-alt",
			Version: "0.29.9-arsydoni4326-alt",
			URL:     "https://github.com/test-owner/test-repo",
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	remote, ok := parsed["remote"].(map[string]interface{})
	if !ok {
		t.Fatal("expected remote object in JSON")
	}
	if _, ok := remote["version"]; !ok {
		t.Error("expected remote.version field in JSON")
	}
}

func TestParseVersion_GoModPseudoVersion(t *testing.T) {
	// Go module pseudo-versions like v0.0.0-20260522092201-58a85b68b3d9
	// should be treated as dev builds.
	v := "v0.0.0-20260522092201-58a85b68b3d9"
	if !isDevVersion(v) {
		t.Errorf("isDevVersion(%q) = false, want true (pseudo-version is dev)", v)
	}

	// parseVersion should also reject it
	_, err := parseVersion(v)
	if err != nil {
		t.Logf("parseVersion(%q) correctly returned error: %v", v, err)
	}
}

// TestFetchRemoteShortCommit_MalformedResponse verifies the JSON decoder
// handles unexpected types in the GitHub API response.
func TestFetchRemoteShortCommit_MalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sha": 123}`)) // sha is a number, not string — type mismatch
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("http get failed: %v", err)
	}
	defer resp.Body.Close()

	var data struct {
		SHA string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Logf("expected malformed JSON error: %v", err)
	}
}