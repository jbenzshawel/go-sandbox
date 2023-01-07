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
	maxTTL            int64
	slidingExpiration bool
	stopSignal        chan struct{}
}

func NewExpirationMap[K comparable, V any](maxDuration time.Duration, slidingExpiration bool) *ExpirationMap[K, V] {
	return &ExpirationMap[K, V]{
		m:                 map[K]*item[V]{},
		maxTTL:            maxDuration.Milliseconds(),
		slidingExpiration: slidingExpiration,
		stopSignal:        make(chan struct{}),
	}
}

func (m *ExpirationMap[K, V]) Len() int {
	return len(m.m)
}

func (m *ExpirationMap[K, V]) Get(k K) (v V) {
	m.l.RLock()
	defer m.l.RUnlock()

	if it, ok := m.m[k]; ok {
		v = it.value
		it.lastAccess = time.Now().UnixMilli()
	}
	return v
}

func (m *ExpirationMap[K, V]) Set(k K, v V) {
	m.l.Lock()
	defer m.l.Unlock()

	startCleanup := len(m.m) == 0

	it, ok := m.m[k]
	if !ok {
		it = &item[V]{
			value: v,
		}
		m.m[k] = it
	}
	it.insertedAt = time.Now().UnixMilli()
	it.lastAccess = it.insertedAt

	if startCleanup {
		go m.startCleanup()
	}
}

func (m *ExpirationMap[K, V]) startCleanup() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-m.stopSignal:
			ticker.Stop()
			break
		case t := <-ticker.C:
			m.DeleteExpired(t)
			if len(m.m) == 0 {
				ticker.Stop()
			}
		}
	}
}

func (m *ExpirationMap[K, V]) stopCleanup() {
	m.stopSignal <- struct{}{}
}

func (m *ExpirationMap[K, V]) Delete(k K) {
	m.l.Lock()
	defer m.l.Unlock()

	delete(m.m, k)

	if len(m.m) == 0 {
		m.stopCleanup()
	}
}

func (m *ExpirationMap[K, V]) DeleteExpired(now time.Time) {
	m.l.Lock()
	defer m.l.Unlock()
	for k, v := range m.m {
		var cachedTime int64
		if m.slidingExpiration {
			cachedTime = now.UnixMilli() - v.lastAccess
		} else {
			cachedTime = now.UnixMilli() - v.insertedAt
		}
		if cachedTime > m.maxTTL {
			delete(m.m, k)
		}
	}
}
