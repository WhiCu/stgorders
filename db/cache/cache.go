package cache

import (
	"log/slog"

	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	MinSize = 2
)

type Cache[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, value V) error
	Size() int
}

type LRUCache[K comparable, V any] struct {
	lru  *lru.TwoQueueCache[K, V]
	size int
	log  *slog.Logger
}

func NewLRUCache[K comparable, V any](size int, log *slog.Logger) (*LRUCache[K, V], error) {
	lru, err := lru.New2Q[K, V](size)
	if err != nil {
		return nil, err
	}
	return &LRUCache[K, V]{
		lru:  lru,
		size: size,
		log:  log,
	}, nil
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
		v, _ := c.lru.Get(k)
		m[k] = v
	}
	return m
}

func (c *LRUCache[K, V]) Size() int {
	return c.size
}
