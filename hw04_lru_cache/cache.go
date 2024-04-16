package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type item struct {
	key   Key
	value interface{}
}

// lruCache кэш, который хранит в себе последние capacity элементов.
type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

// Set добавить элемент в кэш. Если элемент уже в кэше, то обновить значение.
// Возвращает флаг: присутствовал ли элемент в кэше.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	i, ok := c.items[key]
	if ok {
		v := item{key: key, value: value}
		i.Value = v
		c.queue.MoveToFront(i)
		return true
	}
	if c.capacity == c.queue.Len() {
		back := c.queue.Back()
		v, _ := back.Value.(item)
		delete(c.items, v.key)
		c.queue.Remove(back)
	}
	i = c.queue.PushFront(item{key: key, value: value})
	c.items[key] = i
	return false
}

// Get получить элемент из кэша. Возвращает значение и флаг: присутствовал ли элемент в кэше.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	i, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(i)
	v, _ := i.Value.(item)
	return v.value, true
}

// Clear очистить кэш.
func (c *lruCache) Clear() {
	c.mu.Lock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
	c.mu.Unlock()
}

// NewCache создать кэш емкостью capacity.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
