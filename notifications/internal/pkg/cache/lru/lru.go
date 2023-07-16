package lru

import (
	"container/list"
	"sync"
)

// Cache item
type Item[Key comparable, Val any] struct {
	key   Key
	value Val
}

// Custom LRU Cacche implementation
type LRUCache[Key comparable, Val any] struct {
	items    map[Key]*list.Element
	queue    *list.List
	capacity int
	dataMu   sync.RWMutex
}

// Initialize new lru cache
func NewLRUCache[Key comparable, Val any](capacity int) *LRUCache[Key, Val] {
	return &LRUCache[Key, Val]{
		items:    make(map[Key]*list.Element),
		queue:    list.New(),
		capacity: capacity,
	}
}

// Get element from cache by key
func (c *LRUCache[Key, Val]) Get(key Key) (Val, bool) {
	c.dataMu.RLock()
	defer c.dataMu.RUnlock()

	if elem, exists := c.items[key]; exists {
		c.queue.MoveToFront(elem)
		return elem.Value.(*Item[Key, Val]).value, true
	}

	return Item[Key, Val]{}.value, false
}

// Save elemnt in cache
func (c *LRUCache[Key, Val]) Set(key Key, val Val) error {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	// Проверяем, существует ли элемент с заданным ключем
	if elem, exists := c.items[key]; exists {
		elem.Value.(*Item[Key, Val]).value = val
		c.queue.MoveToFront(elem)
		return nil
	}

	// Если элемента нет, создаем новый и добавляем в кэш
	newItem := &Item[Key, Val]{key: key, value: val}
	elem := c.queue.PushFront(newItem)
	c.items[key] = elem

	// Если количество превышает вместимость кэша, удаляем самый старый
	if c.queue.Len() > c.capacity {
		oldestElem := c.queue.Back()
		if oldestElem != nil {
			oldestElem := c.queue.Remove(oldestElem).(*Item[Key, Val])
			delete(c.items, oldestElem.key)
		}
	}

	return nil
}

// Delete element from cache by key
func (c *LRUCache[Key, Val]) Delete(key Key) {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	if elem, exists := c.items[key]; exists {
		delete(c.items, key)
		c.queue.Remove(elem)
	}
}

// Return count of elements in cache
func (c *LRUCache[Key, Val]) Count() int {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	return c.queue.Len()
}

// Clear cache
func (c *LRUCache[Key, Val]) Clear() {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	c.items = make(map[Key]*list.Element)
	c.queue.Init()
}
