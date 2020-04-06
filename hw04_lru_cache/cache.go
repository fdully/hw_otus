package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	// Place your code here
	Set(Key, interface{}) bool
	Get(Key) (interface{}, bool)
	Clear()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*listItem),
	}
}

type lruCache struct {
	// Place your code here:
	capacity int               // - capacity
	queue    List              // - queue
	items    map[Key]*listItem // - items

}

func (l *lruCache) Set(k Key, v interface{}) bool {
	if item, ok := l.items[k]; ok {
		l.queue.MoveToFront(item)
		item.Value.(*cacheItem).value = v
		return true
	}

	itemValue := &cacheItem{
		key:   k,
		value: v,
	}
	item := l.queue.PushFront(itemValue)
	l.items[k] = item

	if l.queue.Len() > l.capacity {
		item := l.queue.Back()
		if item != nil {
			delete(l.items, item.Value.(*cacheItem).key)
			l.queue.Remove(item)
		}
	}
	return false
}

func (l lruCache) Get(k Key) (interface{}, bool) {
	if item, ok := l.items[k]; ok {
		l.queue.MoveToFront(item)
		if item.Value.(*cacheItem) == nil {
			return nil, false
		}
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	for k, v := range l.items {
		delete(l.items, k)
		l.queue.Remove(v)
	}
}

type cacheItem struct {
	// Place your code here
	key   Key
	value interface{}
}
