package executor

import (
	"Nani/internal/app/inhuman"
	"sort"
)

type Pair struct {
	key string
	value int
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].value < p[j].value
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func SortKeywords(keywords inhuman.Keywords) []string {
	pl := make(PairList, len(keywords))
	i := 0
	for k, v := range keywords {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	keys := make([]string, len(keywords))
	for i, v := range pl {
		keys[i] = v.key
	}

	return keys
}
