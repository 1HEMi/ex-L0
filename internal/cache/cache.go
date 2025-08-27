package cache

import (
	"Ex-L0/internal/domain"
	"container/list"
	"sync"
)

type entry struct {
	key   string
	value *domain.Order
}

type Cache struct {
	mx    sync.RWMutex
	cap   int
	list  *list.List
	index map[string]*list.Element
}

func New(cap int) *Cache {
	if cap <= 0 {
		cap = 512
	}
	return &Cache{
		cap:   cap,
		list:  list.New(),
		index: make(map[string]*list.Element, cap),
	}
}

func (c *Cache) Get(key string) (*domain.Order, bool) {
	c.mx.RLock()
	e, ok := c.index[key]
	if !ok {
		c.mx.RUnlock()
		return nil, false
	}
	val := e.Value.(*entry).value
	c.mx.RUnlock()
	c.mx.Lock()
	c.list.MoveToFront(e)
	c.mx.Unlock()
	return val, true
}

func (c *Cache) Set(key string, val *domain.Order) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if e, ok := c.index[key]; ok {
		e.Value = entry{key: key, value: val}
		c.list.MoveToFront(e)
		return
	}
	el := c.list.PushFront(entry{key: key, value: val})
	c.index[key] = el
	if c.list.Len() > c.cap {
		back := c.list.Back()
		if back != nil {
			ent := back.Value.(entry)
			delete(c.index, ent.key)
			c.list.Remove(back)
		}
	}
}
