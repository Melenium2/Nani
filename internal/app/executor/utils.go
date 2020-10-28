package executor

import (
	"Nani/internal/app/inhuman"
	"sort"
)

func SortKeywords(keywords inhuman.Keywords) []string {
	switchValue := make(map[int]string, len(keywords)-1)
	values := make([]int, 0)
	for k, v := range keywords {
		switchValue[v] = k
		values = append(values, v)
	}
	sort.Ints(values)

	keys := make([]string, len(values))
	for i, v := range values {
		keys[i] = switchValue[v]
	}

	return keys
}
