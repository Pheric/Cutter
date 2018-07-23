package server

import (
	"sync"
	"fmt"
)

var GlobalCache = Cache{items:make(map[string][]byte)}

type Cache struct {
	sync.RWMutex
	items map[string][]byte
}

func (c *Cache) Save(name string, data []byte) {
	c.Lock()
	c.items[name] = data
	c.Unlock()
}

func (c *Cache) Get(name string) ([]byte, error) {
	// Returning an error so checking ok is not forgotten later.
	c.RLock()
	defer c.RUnlock()
	d, ok := c.items[name]
	if !ok {
		return nil, fmt.Errorf("error getting file \"%s\" from cache: does not exist")
	}

	return d, nil
}