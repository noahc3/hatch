package utils

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func KeysContains[K any](arr map[string]K, str string) bool {
	for a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
