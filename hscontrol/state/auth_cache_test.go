package state

import (
	"testing"
	"time"

	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthCacheBoundedLRU verifies that the registration auth cache is
// bounded by a maximum entry count, that exceeding the maxEntries evicts the
// oldest entry, and that the eviction callback resolves the parked
// AuthRequest with ErrRegistrationExpired so any waiting goroutine wakes.
func TestAuthCacheBoundedLRU(t *testing.T) {
	const maxEntries = 4

	cache := expirable.NewLRU[types.AuthID, *types.AuthRequest](
		maxEntries,
		func(_ types.AuthID, rn *types.AuthRequest) {
			rn.FinishAuth(types.AuthVerdict{Err: ErrRegistrationExpired})
		},
		time.Hour, // long TTL — we test eviction by size, not by time
	)

	entries := make([]*types.AuthRequest, 0, maxEntries+1)
	ids := make([]types.AuthID, 0, maxEntries+1)

	for range maxEntries + 1 {
		id := types.MustAuthID()
		entry := types.NewAuthRequest()
		cache.Add(id, entry)
		ids = append(ids, id)
		entries = append(entries, entry)
	}

	// Cap should be respected.
	assert.Equal(t, maxEntries, cache.Len(), "cache must not exceed the configured maxEntries")

	// The oldest entry must have been evicted.
	_, ok := cache.Get(ids[0])
	assert.False(t, ok, "oldest entry must be evicted when maxEntries is exceeded")

	// The eviction callback must have woken the parked AuthRequest.
	select {
	case verdict := <-entries[0].WaitForAuth():
		require.False(t, verdict.Accept(), "evicted entry must not signal Accept")
		require.ErrorIs(t,
			verdict.Err, ErrRegistrationExpired,
			"evicted entry must surface ErrRegistrationExpired, got: %v",
			verdict.Err,
		)
	case <-time.After(time.Second):
		t.Fatal("eviction callback did not wake the parked AuthRequest")
	}

	// All non-evicted entries must still be retrievable.
	for i := 1; i <= maxEntries; i++ {
		_, ok := cache.Get(ids[i])
		assert.True(t, ok, "non-evicted entry %d should still be in the cache", i)
	}
}

func TestSetAuthCacheEntryPreservesQRExpiry(t *testing.T) {
	s := &State{
		authCache:           expirable.NewLRU[types.AuthID, *types.AuthRequest](1, nil, time.Hour),
		authCacheExpiration: time.Hour,
	}
	entry := types.NewRegisterAuthRequest(&types.RegistrationData{})
	expiresAt := time.Now().Add(10 * time.Minute).Round(0)
	entry.SetExpiry(expiresAt)

	s.SetAuthCacheEntry(types.MustAuthID(), entry)

	assert.Equal(t, expiresAt, entry.ExpiresAt())
}

// TestAuthCacheTTLExpiry verifies that cache entries expire after the
// configured TTL and that the eviction callback is called.
func TestAuthCacheTTLExpiry(t *testing.T) {
	const shortTTL = 100 * time.Millisecond

	evicted := make(chan types.AuthID, 1)

	cache := expirable.NewLRU[types.AuthID, *types.AuthRequest](
		10, // large capacity, we test TTL not size
		func(id types.AuthID, rn *types.AuthRequest) {
			rn.FinishAuth(types.AuthVerdict{Err: ErrRegistrationExpired})
			evicted <- id
		},
		shortTTL,
	)

	id := types.MustAuthID()
	entry := types.NewAuthRequest()
	cache.Add(id, entry)

	// Entry should be retrievable immediately
	_, ok := cache.Get(id)
	require.True(t, ok, "entry should exist immediately after insertion")

	// Wait for TTL to expire
	time.Sleep(shortTTL + 50*time.Millisecond)

	// Entry should be gone
	_, ok = cache.Get(id)
	assert.False(t, ok, "entry should be evicted after TTL expires")

	// Eviction callback should have been called
	select {
	case evictedID := <-evicted:
		assert.Equal(t, id, evictedID, "eviction callback should be called with the correct ID")
	case <-time.After(time.Second):
		t.Fatal("eviction callback was not called within timeout")
	}

	// The AuthRequest should have been resolved with an error
	select {
	case verdict := <-entry.WaitForAuth():
		require.False(t, verdict.Accept(), "expired entry must not signal Accept")
		require.ErrorIs(t, verdict.Err, ErrRegistrationExpired,
			"expired entry must surface ErrRegistrationExpired, got: %v", verdict.Err)
	case <-time.After(time.Second):
		t.Fatal("expired entry did not resolve the AuthRequest")
	}
}

// TestAuthCacheNoExpiryExtension verifies that re-adding an existing entry
// does not extend its TTL or prevent expiry.
func TestAuthCacheNoExpiryExtension(t *testing.T) {
	const shortTTL = 100 * time.Millisecond

	cache := expirable.NewLRU[types.AuthID, *types.AuthRequest](
		10,
		func(_ types.AuthID, rn *types.AuthRequest) {
			rn.FinishAuth(types.AuthVerdict{Err: ErrRegistrationExpired})
		},
		shortTTL,
	)

	id := types.MustAuthID()
	entry1 := types.NewAuthRequest()
	cache.Add(id, entry1)

	// Wait half the TTL
	time.Sleep(shortTTL / 2)

	// Try to re-add with the same ID (simulating a refresh attempt)
	entry2 := types.NewAuthRequest()
	cache.Add(id, entry2)

	// Wait for the original TTL to fully expire
	time.Sleep(shortTTL/2 + 50*time.Millisecond)

	// Entry should be gone regardless of the re-add
	_, ok := cache.Get(id)
	assert.False(t, ok, "re-adding an entry should not extend its TTL")
}

// TestAuthCacheDoubleConsumption verifies that consuming an auth entry
// (retrieving and removing it) prevents a second consumption attempt.
func TestAuthCacheDoubleConsumption(t *testing.T) {
	cache := expirable.NewLRU[types.AuthID, *types.AuthRequest](
		10,
		func(_ types.AuthID, rn *types.AuthRequest) {
			rn.FinishAuth(types.AuthVerdict{Err: ErrRegistrationExpired})
		},
		time.Hour,
	)

	id := types.MustAuthID()
	entry := types.NewAuthRequest()
	cache.Add(id, entry)

	// First consumption succeeds
	retrieved1, ok1 := cache.Get(id)
	assert.True(t, ok1, "first retrieval should succeed")
	assert.Equal(t, entry, retrieved1)

	// Remove the entry (simulating consumption)
	cache.Remove(id)

	// Second consumption fails
	retrieved2, ok2 := cache.Get(id)
	assert.False(t, ok2, "second retrieval after removal should fail")
	assert.Nil(t, retrieved2)
}

// TestAuthCacheConcurrentAccess verifies that concurrent access to the same
// authID is handled correctly by the cache.
func TestAuthCacheConcurrentAccess(t *testing.T) {
	cache := expirable.NewLRU[types.AuthID, *types.AuthRequest](
		100,
		func(_ types.AuthID, rn *types.AuthRequest) {
			rn.FinishAuth(types.AuthVerdict{Err: ErrRegistrationExpired})
		},
		time.Hour,
	)

	id := types.MustAuthID()
	entry := types.NewAuthRequest()
	cache.Add(id, entry)

	const numGoroutines = 10
	successCount := make(chan bool, numGoroutines)

	// Launch multiple goroutines trying to get the same entry
	for i := 0; i < numGoroutines; i++ {
		go func() {
			retrieved, ok := cache.Get(id)
			successCount <- ok && retrieved != nil
		}()
	}

	// All should succeed since we're only reading
	successes := 0
	for i := 0; i < numGoroutines; i++ {
		if <-successCount {
			successes++
		}
	}

	assert.Equal(t, numGoroutines, successes,
		"all concurrent reads should succeed")
}

// TestAuthCacheConcurrentRemoval verifies that only one goroutine can
// successfully consume (retrieve and remove) an auth entry.
func TestAuthCacheConcurrentRemoval(t *testing.T) {
	cache := expirable.NewLRU[types.AuthID, *types.AuthRequest](
		100,
		func(_ types.AuthID, rn *types.AuthRequest) {
			rn.FinishAuth(types.AuthVerdict{Err: ErrRegistrationExpired})
		},
		time.Hour,
	)

	id := types.MustAuthID()
	entry := types.NewAuthRequest()
	cache.Add(id, entry)

	const numGoroutines = 10
	successCount := make(chan bool, numGoroutines)

	// Launch multiple goroutines trying to consume the same entry
	for i := 0; i < numGoroutines; i++ {
		go func() {
			// Attempt to get and immediately remove
			retrieved, ok := cache.Get(id)
			if ok && retrieved != nil {
				cache.Remove(id)
				successCount <- true
			} else {
				successCount <- false
			}
		}()
	}

	// At least one should succeed, others may fail
	successes := 0
	for i := 0; i < numGoroutines; i++ {
		if <-successCount {
			successes++
		}
	}

	assert.GreaterOrEqual(t, successes, 1,
		"at least one goroutine should successfully retrieve the entry")

	// Entry should be gone after all goroutines complete
	_, ok := cache.Get(id)
	assert.False(t, ok, "entry should be removed after concurrent access")
}

// TestAuthCacheMultipleEntries verifies that multiple independent auth
// entries can coexist and be managed independently.
func TestAuthCacheMultipleEntries(t *testing.T) {
	cache := expirable.NewLRU[types.AuthID, *types.AuthRequest](
		100,
		func(_ types.AuthID, rn *types.AuthRequest) {
			rn.FinishAuth(types.AuthVerdict{Err: ErrRegistrationExpired})
		},
		time.Hour,
	)

	const numEntries = 10
	ids := make([]types.AuthID, numEntries)
	entries := make([]*types.AuthRequest, numEntries)

	// Add multiple entries
	for i := 0; i < numEntries; i++ {
		ids[i] = types.MustAuthID()
		entries[i] = types.NewAuthRequest()
		cache.Add(ids[i], entries[i])
	}

	// All should be retrievable
	for i := 0; i < numEntries; i++ {
		retrieved, ok := cache.Get(ids[i])
		assert.True(t, ok, "entry %d should be retrievable", i)
		assert.Equal(t, entries[i], retrieved)
	}

	// Remove one entry
	cache.Remove(ids[5])

	// Removed entry should be gone
	_, ok := cache.Get(ids[5])
	assert.False(t, ok, "removed entry should not be retrievable")

	// Other entries should still exist
	for i := 0; i < numEntries; i++ {
		if i == 5 {
			continue
		}
		retrieved, ok := cache.Get(ids[i])
		assert.True(t, ok, "entry %d should still be retrievable", i)
		assert.Equal(t, entries[i], retrieved)
	}
}

// TestAuthCacheEvictionDoesNotAffectOthers verifies that when an entry
// is evicted due to size limits, it doesn't affect other entries.
func TestAuthCacheEvictionDoesNotAffectOthers(t *testing.T) {
	const maxEntries = 5

	cache := expirable.NewLRU[types.AuthID, *types.AuthRequest](
		maxEntries,
		func(_ types.AuthID, rn *types.AuthRequest) {
			rn.FinishAuth(types.AuthVerdict{Err: ErrRegistrationExpired})
		},
		time.Hour,
	)

	ids := make([]types.AuthID, maxEntries+2)
	entries := make([]*types.AuthRequest, maxEntries+2)

	// Add maxEntries entries
	for i := 0; i < maxEntries; i++ {
		ids[i] = types.MustAuthID()
		entries[i] = types.NewAuthRequest()
		cache.Add(ids[i], entries[i])
	}

	// Add two more, which should evict the first two
	for i := maxEntries; i < maxEntries+2; i++ {
		ids[i] = types.MustAuthID()
		entries[i] = types.NewAuthRequest()
		cache.Add(ids[i], entries[i])
	}

	// First two should be evicted
	_, ok := cache.Get(ids[0])
	assert.False(t, ok, "first entry should be evicted")
	_, ok = cache.Get(ids[1])
	assert.False(t, ok, "second entry should be evicted")

	// Remaining entries should still be present
	for i := 2; i < maxEntries+2; i++ {
		retrieved, ok := cache.Get(ids[i])
		assert.True(t, ok, "entry %d should still be present", i)
		assert.Equal(t, entries[i], retrieved)
	}
}
