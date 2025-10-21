package cache

import (
	"container/list"
	"sync"

	"l0/cons/internal/models"
)

type Cache struct {
	capacity int
	data     map[string]*list.Element
	mu       sync.RWMutex
	l        *list.List
}

type entity struct {
	key        string
	orderValue *models.Order
}

func NewCache(cap int) *Cache {
	return &Cache{
		capacity: cap,
		data:     make(map[string]*list.Element, cap),
		l:        list.New(),
	}
}

func (c *Cache) Get(orderID string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[orderID]
	if !ok {
		return nil, false
	}

	c.l.MoveToFront(val)

	return val.Value.(*entity).orderValue, true
}

func (c *Cache) Put(orderID string, order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[orderID]
	if ok {
		val.Value.(*entity).orderValue = order
		c.l.MoveToFront(val)
		return
	}

	if len(c.data) >= c.capacity {
		last := c.l.Back()
		if last != nil {
			c.l.Remove(last)
			delete(c.data, last.Value.(*entity).key)
		}
	}

	el := c.l.PushFront(&entity{key: orderID, orderValue: order})
	c.data[orderID] = el
}
