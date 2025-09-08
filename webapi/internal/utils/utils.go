package utils

import "strings"

func ContainsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

func HasSubstringKeys[V any](m map[string]V) bool {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	for i := 0; i < len(keys); i++ {
		for j := 0; j < len(keys); j++ {
			if i == j {
				continue
			}
			if strings.Contains(strings.ToLower(keys[i]), strings.ToLower(keys[j])) {
				return true
			}
		}
	}
	return false
}
