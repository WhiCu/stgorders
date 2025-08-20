package cache

import (
	"log/slog"

	lru "github.com/hashicorp/golang-lru/v2"
)

type LRUCache[K comparable, V any] struct {
	lru  *lru.TwoQueueCache[K, V]
	Size int
	log  *slog.Logger
}

func NewLRUCache[K comparable, V any](size int, log *slog.Logger) *LRUCache[K, V] {
	if size <= 3 {
		panic("invalid size")
	}
	lru, err := lru.New2Q[K, V](size)
	if err != nil {
		panic(err)
	}
	return &LRUCache[K, V]{
		lru:  lru,
		Size: size,
		log:  log,
	}
}

func (c *LRUCache[K, V]) Get(key K) (V, error) {
	value, ok := c.lru.Get(key)
	if !ok {
		return value, ErrNotFound
	}
	c.log.Debug("cache get", slog.Any("key", key))
	return value, nil
}

func (c *LRUCache[K, V]) Set(key K, value V) error {
	c.lru.Add(key, value)
	c.log.Debug("cache set", slog.Any("key", key))
	return nil
}

func (c *LRUCache[K, V]) Map() map[K]V {
	m := make(map[K]V, c.lru.Len())
	for _, k := range c.lru.Keys() {
		v, ok := c.lru.Get(k)
		if !ok {
			continue
		}
		m[k] = v
	}
	return m
}
