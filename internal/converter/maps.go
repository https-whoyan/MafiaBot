package converter

import (
	"cmp"
)

// Почему еще не добавили эти функции в стандартную библиотеку....

func GetMapKeys[K cmp.Ordered, E any](mp map[K]E) []K {
	var keys []K
	for key := range mp {
		keys = append(keys, key)
	}
	return keys
}

func GetMapValues[K cmp.Ordered, E any](mp map[K]E) []E {
	var values []E
	for _, val := range mp {
		values = append(values, val)
	}
	return values
}

func SliceToSet[S ~[]E, E comparable](sl S) map[E]bool {
	mp := make(map[E]bool)
	for _, val := range sl {
		mp[val] = true
	}
	return mp
}
