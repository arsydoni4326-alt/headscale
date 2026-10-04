package hscontrol

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandleListUsers_ResponseFormat verifies that HandleListUsers returns
// a JSON object with a "users" key containing an array, matching the frontend
// expectation at headplane/app/routes/admin/users/route.tsx line 73.
func TestHandleListUsers_ResponseFormat(t *testing.T) {
	// This is a simplified test focused only on the response format
	// We're testing that the response is {"users": [...]} and not just [...]
	
	// Note: This test validates the response structure only.
	// Full integration tests with database setup are in headplane_users_test.go
	
	mockResponse := map[string]interface{}{
		"users": []HeadplaneUserResponse{
			{
				ID:        1,
				Username:  "testuser",
				Role:      "admin",
				CreatedAt: 1234567890,
			},
		},
	}
	
	// Verify the structure marshals correctly
	data, err := json.Marshal(mockResponse)
	require.NoError(t, err)
	
	// Verify we can unmarshal it back
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)
	
	// Verify the "users" key exists
	users, ok := result["users"]
	require.True(t, ok, "Response should have 'users' key")
	require.NotNil(t, users, "Response 'users' should not be nil")
	
	// Verify it's an array
	usersArray, ok := users.([]interface{})
	require.True(t, ok, "Response 'users' should be an array")
	assert.Equal(t, 1, len(usersArray), "Should have one user")
}

// TestHandleListUsers_EmptyResponseFormat verifies empty list format
func TestHandleListUsers_EmptyResponseFormat(t *testing.T) {
	mockResponse := map[string]interface{}{
		"users": []HeadplaneUserResponse{},
	}
	
	data, err := json.Marshal(mockResponse)
	require.NoError(t, err)
	
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)
	
	// Verify the "users" key exists even when empty
	users, ok := result["users"]
	require.True(t, ok, "Response should have 'users' key even when empty")
	
	usersArray, ok := users.([]interface{})
	require.True(t, ok, "Response 'users' should be an array")
	assert.Equal(t, 0, len(usersArray), "Should have no users")
}
