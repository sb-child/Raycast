package utility

func IndexOf[T comparable](collection []T, el T) int {
	for i, x := range collection {
		if x == el {
			return i
		}
	}
	return -1
}

func IndexOf2[T comparable](collection []T, el T) (int, bool) {
	for i, x := range collection {
		if x == el {
			return i, true
		}
	}
	return -1, false
}
