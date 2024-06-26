package converter

func GetElsOnlyIncludeFunc[S ~[]M, M any, K ~[]E, E comparable](
	operatedSet S, keys K, getKeyByEl func(m M) E) (S, int, error) {
	mpKeys := make(map[E]bool)
	for _, key := range keys {
		mpKeys[key] = true
	}

	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	var ans S
	for _, member := range operatedSet {
		memberKey := getKeyByEl(member)
		if mpKeys[memberKey] {
			ans = append(ans, member)
		}
	}

	return ans, len(ans), err
}

func SetDiff[S ~[]E, E comparable](operatedSet S, needNotIncludeS S) (S, int) {
	mpUsed := make(map[E]bool)
	for _, member := range needNotIncludeS {
		mpUsed[member] = true
	}
	var ans S
	for _, member := range operatedSet {
		if !mpUsed[member] {
			ans = append(ans, member)
		}
	}
	return ans, len(ans)
}

func SliceToSet[S ~[]E, E comparable](sl S) map[E]bool {
	mp := make(map[E]bool)
	for _, val := range sl {
		mp[val] = true
	}
	return mp
}
