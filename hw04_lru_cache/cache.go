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

type hash struct {
	key   Key
	value interface{}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	v, isExisted := lc.items[key]

	if isExisted {
		v.Value.(*hash).value = value
		lc.queue.MoveToFront(v)

		return isExisted
	}

	h := &hash{key: key, value: value}
	newItem := lc.queue.PushFront(h)
	lc.items[key] = newItem

	if lc.queue.Len() > lc.capacity {
		back := lc.queue.Back()

		if back != nil {
			delete(lc.items, back.Value.(*hash).key)
			lc.queue.Remove(back)
		}
	}

	return isExisted
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	v, isExisted := lc.items[key]

	if isExisted {
		lc.queue.MoveToFront(v)

		return v.Value.(*hash).value, isExisted
	}

	return nil, isExisted
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
