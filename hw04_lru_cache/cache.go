package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type CacheItem struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, exists := c.items[key]; exists {
		item.Value = CacheItem{Key: key, Value: value}
		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() >= c.capacity {
		backItem := c.queue.Back()
		c.queue.Remove(backItem)
		delete(c.items, backItem.Value.(CacheItem).Key)
	}

	listItem := c.queue.PushFront(CacheItem{Key: key, Value: value})
	c.items[key] = listItem
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		return item.Value.(CacheItem).Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
