package cache

import (
	"sync"
	"time"
)

type item[V any] struct {
	value      V
	insertedAt int64
	lastAccess int64
}

type ExpirationMap[K comparable, V any] struct {
	m                 map[K]*item[V]
	l                 sync.RWMutex
	maxTTL            int
	slidingExpiration bool
}

func New[K comparable, V any](ln int, maxTTL int, slidingExpiration bool) (m *ExpirationMap[K, V]) {
	m = &ExpirationMap[K, V]{
		m:                 make(map[K]*item[V], ln),
		maxTTL:            maxTTL,
		slidingExpiration: slidingExpiration,
	}
	// TODO: Update expiration logic to only run when cache map has values
	go func() {
		for now := range time.Tick(time.Second) {
			m.l.Lock()
			for k, v := range m.m {
				var cachedTime int64
				if m.slidingExpiration {
					cachedTime = now.Unix() - v.lastAccess
				} else {
					cachedTime = now.Unix() - v.insertedAt
				}
				if cachedTime > int64(maxTTL) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return
}

func (m *ExpirationMap[K, V]) Len() int {
	return len(m.m)
}

func (m *ExpirationMap[K, V]) Put(k K, v V) {
	m.l.Lock()
	defer m.l.Unlock()
	it, ok := m.m[k]
	if !ok {
		it = &item[V]{
			value: v,
		}
		m.m[k] = it
	}
	it.insertedAt = time.Now().Unix()
	it.lastAccess = it.insertedAt
}

func (m *ExpirationMap[K, V]) Get(k K) (v V) {
	m.l.RLock()
	defer m.l.RUnlock()

	if it, ok := m.m[k]; ok {
		v = it.value
		it.lastAccess = time.Now().Unix()
	}
	return v
}

func (m *ExpirationMap[K, V]) Delete(k K) {
	m.l.Lock()
	defer m.l.Unlock()

	delete(m.m, k)
}
