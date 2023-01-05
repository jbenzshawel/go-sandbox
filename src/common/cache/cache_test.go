package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpirationMap_New(t *testing.T) {
	expMap := NewExpirationMap[string, string](10*time.Minute, true)
	require.NotNil(t, expMap.m)
	assert.Len(t, expMap.m, 0)
	assert.NotNil(t, expMap.m)
	assert.Equal(t, (10 * time.Minute).Milliseconds(), expMap.maxTTL)
	assert.True(t, expMap.slidingExpiration)
}

func TestExpirationMap_Set(t *testing.T) {
	now := time.Now().UnixMilli()
	expMap := NewExpirationMap[string, string](10*time.Minute, true)
	expMap.Set("key", "value")

	cachedItem := expMap.m["key"]
	assert.Equal(t, "value", cachedItem.value)
	assert.GreaterOrEqual(t, cachedItem.insertedAt, now)
	assert.GreaterOrEqual(t, cachedItem.lastAccess, now)
	assert.Equal(t, cachedItem.insertedAt, cachedItem.lastAccess)
}

func TestExpirationMap_Get(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expMap := NewExpirationMap[string, string](10*time.Minute, true)
	expMap.Set("key", "value")

	delay := 1 * time.Second
	time.Sleep(delay)

	value := expMap.Get("key")

	assert.Equal(t, "value", value)

	cachedItem := expMap.m["key"]
	assert.Equal(t, "value", cachedItem.value)
	assert.GreaterOrEqual(t, cachedItem.insertedAt, now.UnixMilli())
	assert.GreaterOrEqual(t, cachedItem.lastAccess, now.Add(delay).UnixMilli())
}

func TestExpirationMap_Delete(t *testing.T) {
	expMap := NewExpirationMap[string, *struct{}](10*time.Minute, true)
	expMap.Set("key", &struct{}{})

	expMap.Delete("key")

	assert.Nil(t, expMap.Get("key"))
	assert.Nil(t, expMap.m["key"])
}

func TestExpirationMap_DeleteExpired(t *testing.T) {
	expMap := NewExpirationMap[string, string](10*time.Minute, false)

	expMap.Set("key1", "val1")
	expMap.Set("key2", "val2")
	cachedItem2 := expMap.m["key2"]
	require.NotNil(t, cachedItem2)
	cachedItem2.insertedAt = time.Now().Add(5 * time.Minute).UnixMilli()

	expMap.Set("key3", "val3")
	cachedItem3 := expMap.m["key3"]
	require.NotNil(t, cachedItem3)
	cachedItem3.insertedAt = time.Now().Add(10 * time.Minute).UnixMilli()

	require.Equal(t, 3, expMap.Len())

	expMap.DeleteExpired(time.Now())
	assert.Equal(t, 3, expMap.Len())

	expMap.DeleteExpired(time.Now().Add(11 * time.Minute))
	assert.Equal(t, 2, expMap.Len())

	expMap.DeleteExpired(time.Now().Add(16 * time.Minute))
	assert.Equal(t, 1, expMap.Len())

	expMap.DeleteExpired(time.Now().Add(21 * time.Minute))
	assert.Equal(t, 0, expMap.Len())
}

func TestExpirationMap_ClearsExpiredWithoutSliding(t *testing.T) {
	t.Parallel()

	expMap := NewExpirationMap[string, string](1*time.Second, false)

	for i := 0; i < 100; i++ {
		k, v := fmt.Sprint("key", i), fmt.Sprint("value", i)
		expMap.Set(k, v)
	}

	require.Equal(t, 100, expMap.Len())
	time.Sleep(3 * time.Second)
	assert.Equal(t, 0, expMap.Len())
}

func TestExpirationMap_ClearsExpiredWithSliding(t *testing.T) {
	t.Parallel()

	expMap := NewExpirationMap[string, string](1*time.Second, true)

	for i := 0; i < 5; i++ {
		k, v := fmt.Sprint("key", i), fmt.Sprint("value", i)
		expMap.Set(k, v)
	}

	require.Equal(t, 5, expMap.Len())

	time.Sleep(500 * time.Millisecond)

	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		expMap.Get(fmt.Sprint("key", i))
	}

	assert.Greater(t, expMap.Len(), 0)

	time.Sleep(2 * time.Second)
	assert.Equal(t, 0, expMap.Len())
}
