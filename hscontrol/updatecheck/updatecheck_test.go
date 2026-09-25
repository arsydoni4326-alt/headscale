package updatecheck

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_WithoutCheckParam(t *testing.T) {
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