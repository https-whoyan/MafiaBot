package converter

// Same as /internal/converter/maps.go
// To avoid import /internal/*

func SliceToSet[S ~[]E, E comparable](sl S) map[E]bool {
	mp := make(map[E]bool)
	for _, val := range sl {
		mp[val] = true
	}
	return mp
}
