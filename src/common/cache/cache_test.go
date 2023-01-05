package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpirationMap_New(t *testing.T) {
	expMap := New[string, string](3, int(10*time.Minute), true)
	require.NotNil(t, expMap.m)
	assert.Len(t, expMap.m, 0)
	assert.NotNil(t, expMap.m)
	assert.Equal(t, int(10*time.Minute), expMap.maxTTL)
	assert.True(t, expMap.slidingExpiration)
}

func TestExpirationMap_Put(t *testing.T) {
	now := time.Now().Unix()
	expMap := New[string, string](3, int(10*time.Minute), true)
	expMap.Put("key", "value")

	cachedItem := expMap.m["key"]
	assert.Equal(t, "value", cachedItem.value)
	assert.GreaterOrEqual(t, cachedItem.insertedAt, now)
	assert.GreaterOrEqual(t, cachedItem.lastAccess, now)
	assert.Equal(t, cachedItem.insertedAt, cachedItem.lastAccess)
}

func TestExpirationMap_Get(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expMap := New[string, string](3, int(10*time.Minute), true)
	expMap.Put("key", "value")

	delay := 1 * time.Second
	time.Sleep(delay)

	value := expMap.Get("key")

	assert.Equal(t, "value", value)

	cachedItem := expMap.m["key"]
	assert.Equal(t, "value", cachedItem.value)
	assert.GreaterOrEqual(t, cachedItem.insertedAt, now.Unix())
	assert.GreaterOrEqual(t, cachedItem.lastAccess, now.Add(delay).Unix())
}

func TestExpirationMap_Delete(t *testing.T) {
	expMap := New[string, *struct{}](3, int(10*time.Minute), true)
	expMap.Put("key", &struct{}{})

	expMap.Delete("key")

	assert.Nil(t, expMap.Get("key"))
	assert.Nil(t, expMap.m["key"])
}

func TestExpirationMap_ClearsExpiredWithoutSliding(t *testing.T) {
	t.Parallel()

	expMap := New[string, string](100, 1, false)

	for i := 0; i < 100; i++ {
		k, v := fmt.Sprint("key", i), fmt.Sprint("value", i)
		expMap.Put(k, v)
	}

	require.Equal(t, 100, expMap.Len())
	time.Sleep(3 * time.Second)
	assert.Equal(t, 0, expMap.Len())
}

func TestExpirationMap_ClearsExpiredWithSliding(t *testing.T) {
	t.Parallel()

	expMap := New[string, string](100, 1, true)

	for i := 0; i < 5; i++ {
		k, v := fmt.Sprint("key", i), fmt.Sprint("value", i)
		expMap.Put(k, v)
	}

	require.Equal(t, 5, expMap.Len())

	time.Sleep(950 * time.Millisecond)

	for i := 0; i < 5; i++ {
		time.Sleep(250 * time.Millisecond)
		expMap.Get(fmt.Sprint("key", i))
	}

	assert.Greater(t, expMap.Len(), 0)

	time.Sleep(2 * time.Second)
	assert.Equal(t, 0, expMap.Len())
}
