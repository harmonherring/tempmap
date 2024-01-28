package tempmap

import (
	"sync"
	"time"
)

type Map[K comparable, V any] struct {
	values map[K]tempValue[V]
	sync.RWMutex
}

type tempValue[V any] struct {
	cancel chan struct{}
	value  V
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		values: map[K]tempValue[V]{},
	}
}

// Get retrieves an item from the temporary map and returns existence status as a bool
func (m *Map[K, V]) Get(k K) (V, bool) {
	m.RLock()
	defer m.RUnlock()
	if value, exists := m.values[k]; exists {
		return value.value, true
	}

	var nothing V
	return nothing, false
}

// Put inserts a value into the temporary map
func (m *Map[K, V]) Put(k K, v V, lifetime time.Duration) {
	if value, exists := m.values[k]; exists {
		value.cancel <- struct{}{}
	}

	m.Lock()
	defer m.Unlock()

	m.values[k] = tempValue[V]{
		cancel: make(chan struct{}),
		value:  v,
	}

	go func() {
		select {
		case <-time.After(lifetime):
			m.Lock()
			defer m.Unlock()
			delete(m.values, k)
		case <-m.values[k].cancel:
			return
		}
	}()
}

// Close closes all pending deletion goroutines
func (m *Map[K, V]) Close() {
	for _, v := range m.values {
		v.cancel <- struct{}{}
	}
}
