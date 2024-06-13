package converter

import "cmp"

// Same as /internal/bot/converter/maps.go
// To avoid import /internal/bot/*

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
