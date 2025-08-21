package cache

type NOPCache[K comparable, V any] struct {
	defaultValue V
	size         int
}

func NewNOPCache[K comparable, V any](defaultValue V) *NOPCache[K, V] {
	return &NOPCache[K, V]{
		defaultValue: defaultValue,
		size:         0,
	}
}

func (c *NOPCache[K, V]) Size() int {
	return c.size
}

func (c *NOPCache[K, V]) Get(key K) (V, error) {
	return c.defaultValue, ErrNotFound
}

func (c *NOPCache[K, V]) Set(key K, value V) error {
	return ErrSetCache
}
