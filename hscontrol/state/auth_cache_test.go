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
