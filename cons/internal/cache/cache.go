package cache

import (
	"l0/cons/internal/models"
	"sync"
)

type Cache struct {
	data map[string]*models.Order
	mu sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]*models.Order),
	}
}

func (c *Cache) Get(orderID string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[orderID]
	return val, ok
}

func (c *Cache) Set(orderID string, order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[orderID] = order
}

func (c *Cache) GetAll() []*models.Order {
	c.mu.Lock()
	defer c.mu.Unlock()

	var res []*models.Order

	for _, order := range c.data {
		res = append(res, order)
	}

	return res
}