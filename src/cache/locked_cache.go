package cache

import "sync"

type mapType map[string][]byte

type LockedCache struct {
	*sync.Mutex
	data mapType
}

// NewLockedCache creates and returns a *LockedCache with prefill values.
func NewLockedCache(data mapType) *LockedCache {
	if data == nil {
		data = make(mapType)
	}
	return &LockedCache{
		Mutex: &sync.Mutex{},
		data:  data,
	}
}

// Add add a key / value pair into the cache. If a value
// already exist, replace it. It is threadsafe.
func (c *LockedCache) Add(key string, value []byte) {
	c.Lock()
	c.data[key] = value
	c.Unlock()
}

// GetWithStatus return the value into the cache and a boolean that
// indicates if a value was found. It is threadsafe. If `c`
// is nil, act as an empty cache.
func (c *LockedCache) GetWithStatus(key string) ([]byte, bool) {
	c.Lock()
	if c == nil {
		return nil, false
	}
	v, ok := c.data[key]
	c.Unlock()
	return v, ok
}

// Get return the value into the cache. It is threadsafe.
// If `c` is nil, act as an empty cache.
func (c *LockedCache) Get(key string) []byte {
	if c == nil {
		return nil
	}
	c.Lock()
	v, _ := c.data[key]
	c.Unlock()
	return v
}
