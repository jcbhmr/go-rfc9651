package orderedmap

import (
	"iter"
	"slices"
)

type OrderedMap[K comparable, V any] struct {
	orderedKeys  []K
	unorderedMap map[K]V
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{
		orderedKeys:  []K{},
		unorderedMap: map[K]V{},
	}
}

func (om *OrderedMap[K, V]) Get(key K) (value V, ok bool) {
	value, ok = om.unorderedMap[key]
	return
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	if _, ok := om.unorderedMap[key]; !ok {
		om.orderedKeys = append(om.orderedKeys, key)
	}
	om.unorderedMap[key] = value
}

func (om *OrderedMap[K, V]) Has(key K) bool {
	_, ok := om.unorderedMap[key]
	return ok
}

func (om *OrderedMap[K, V]) Delete(key K) {
	if _, ok := om.unorderedMap[key]; ok {
		delete(om.unorderedMap, key)
		for i, k := range om.orderedKeys {
			if k == key {
				om.orderedKeys = append(om.orderedKeys[:i], om.orderedKeys[i+1:]...)
				break
			}
		}
	}
}

func (om *OrderedMap[K, V]) Len() int {
	return len(om.unorderedMap)
}

func (om *OrderedMap[K, V]) Keys() iter.Seq[K] {
	return slices.Values(om.orderedKeys)
}

func (om *OrderedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, key := range om.orderedKeys {
			value := om.unorderedMap[key]
			if !yield(value) {
				return
			}
		}
	}
}

func (om *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, key := range om.orderedKeys {
			value := om.unorderedMap[key]
			if !yield(key, value) {
				return
			}
		}
	}
}
