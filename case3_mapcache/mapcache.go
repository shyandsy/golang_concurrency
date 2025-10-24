package cache

import (
	"log"
	"sync"
	"time"
)

type mapValue struct {
	value    string
	deadline time.Time
}

type mapCache struct {
	items map[string]mapValue
	mtx   *sync.RWMutex
	quit  chan struct{}
}

func NewMapCahe() Cache {
	c := &mapCache{
		items: make(map[string]mapValue, 100),
		mtx:   &sync.RWMutex{},
		quit:  make(chan struct{}),
	}
	go c.autoExpiration()
	return c
}

func (c *mapCache) Set(key string, value string, ttl time.Duration) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.items[key] = mapValue{
		value:    value,
		deadline: time.Now().Add(ttl),
	}
}

func (c *mapCache) Get(key string) (string, bool) {
	c.mtx.RLock()
	val, ok := c.items[key]
	c.mtx.RUnlock()

	if !ok {
		return "", false
	}

	if time.Now().After(val.deadline) {
		c.mtx.Lock()
		delete(c.items, key)
		c.mtx.Unlock()
		return "", false
	}

	return val.value, true
}

func (c *mapCache) Close() {
	c.quit <- struct{}{}
}

func (c *mapCache) autoExpiration() {
	tick := time.NewTicker(300 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			itemsToRemove := []string{}

			// check
			c.mtx.RLock()
			for k, v := range c.items {
				if time.Now().After(v.deadline) {
					itemsToRemove = append(itemsToRemove, k)
				}
			}
			c.mtx.RUnlock()

			// delete
			if len(itemsToRemove) > 0 {
				c.mtx.Lock()
				for _, k := range itemsToRemove {
					delete(c.items, k)
				}
				c.mtx.Unlock()
			}
		case <-c.quit:
			log.Println("autoExpiration quit")
			return
		}
	}
}
