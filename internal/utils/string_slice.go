package utils

// RemoveValueFromSlice with only one of needle in the haystack slice. Will not preserve order.
func RemoveValueFromSlice(haystack []string, needle string) []string {
	opSlice := haystack

	for idx, val := range haystack {
		if needle == val {
			opSlice[idx] = opSlice[len(opSlice)-1]
			opSlice = opSlice[:len(opSlice)-1]
			return opSlice
		}
	}

	return haystack
}
