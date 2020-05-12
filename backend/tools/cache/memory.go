package cache

import (
	"errors"
	"sync"
)

var ErrorNotFound error = errors.New("not found")

type memoryCache struct {
	c map[string]interface{}
	l sync.RWMutex
}

func (t *memoryCache) Get(key string) (interface{}, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	v, ok := t.c[key]
	if !ok {
		return nil, ErrorNotFound
	}
	return v, nil
}
func (t *memoryCache) Set(key string, value interface{}) error {
	t.l.Lock()
	defer t.l.Unlock()
	t.c[key] = value
	return nil
}
func (t *memoryCache) Remove(key string) {
	t.l.Lock()
	defer t.l.Unlock()
	delete(t.c, key)
}
func (t *memoryCache) Clear() {
	t.l.Lock()
	defer t.l.Unlock()
	t.c = make(map[string]interface{})
}
