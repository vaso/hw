package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// CacheItem stores item structure to keep back-link with cache.items map.
type CacheItem struct {
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
	elementValue := CacheItem{key: key, value: value}
	el, ok := cache.items[key]
	if ok {
		el.Value = elementValue
		cache.queue.MoveToFront(el)
		return true
	}

	el = cache.queue.PushFront(elementValue)
	cache.items[key] = el
	if cache.queue.Len() <= cache.capacity {
		return false
	}

	lastEl := cache.queue.Back()
	if lastEl != nil {
		cache.queue.Remove(lastEl)
		lastElValue := lastEl.Value.(CacheItem)
		delete(cache.items, lastElValue.key)
	}
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	el, ok := cache.items[key]
	if !ok {
		return nil, false
	}
	if el == nil {
		return nil, false
	}

	cache.queue.MoveToFront(el)
	elValue := el.Value.(CacheItem)

	return elValue.value, true
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}
