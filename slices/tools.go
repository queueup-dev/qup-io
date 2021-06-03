package slices

func HasString(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func CountedHasString(needle string, haystack []string, count int) (int, bool) {
	currentCount := 0
	for k, item := range haystack {
		if item == needle {
			if count >= currentCount {
				return k, true
			}
			currentCount++
		}
	}
	return 0, false
}
