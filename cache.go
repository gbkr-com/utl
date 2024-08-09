package utl

import "time"

// Cache is an in-memory cache of items of type V each having a key K. There is
// a fixed time-to-live for any item. Items retrieved that exceed that duration
// are replaced from another source, such as a Redis database.
type Cache[K comparable, V any] struct {
	values   map[K]V
	inserted map[K]time.Time
	ttl      time.Duration
	replace  func(K) (V, bool)
}

// NewCache returns a new [*Cache] ready to use. The function argument is used
// to replace expired items: given a key value it should return an item V and
// a boolean set to true if the item was found.
func NewCache[K comparable, V any](ttl time.Duration, replace func(K) (V, bool)) *Cache[K, V] {
	return &Cache[K, V]{
		values:   map[K]V{},
		inserted: map[K]time.Time{},
		ttl:      ttl,
		replace:  replace,
	}
}

// Get the item having the given key. Return also true if the item was found,
// otherwise false.
func (x *Cache[K, V]) Get(key K) (item V, ok bool) {

	when := x.inserted[key]

	if time.Since(when) > x.ttl {
		item, ok = x.replace(key)
		if !ok {
			return
		}
		x.values[key] = item
		x.inserted[key] = time.Now()
		return
	}

	item, ok = x.values[key]
	return
}
