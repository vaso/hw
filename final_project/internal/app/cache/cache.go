package cache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// Item stores item structure to keep back-link with cache.items map.
type Item struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	elementValue := Item{key: key, value: value}
	el, ok := cache.items[key]
	if ok {
		// update existing el
		el.Value = elementValue
		cache.queue.MoveToFront(el)
		return true
	}

	if cache.queue.Len() < cache.capacity {
		// add new el
		el = cache.queue.PushFront(elementValue)
		cache.items[key] = el
		return false
	}

	// remove last in queue if necessary
	lastEl := cache.queue.Back()
	if lastEl != nil {
		cache.queue.Remove(lastEl)
		lastElValue := lastEl.Value.(Item)
		delete(cache.items, lastElValue.key)
	}
	el = cache.queue.PushFront(elementValue)
	cache.items[key] = el

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	el, ok := cache.items[key]
	if !ok {
		return nil, false
	}
	if el == nil {
		return nil, false
	}

	cache.queue.MoveToFront(el)
	elValue := el.Value.(Item)

	return elValue.value, true
}

func (cache *lruCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}
