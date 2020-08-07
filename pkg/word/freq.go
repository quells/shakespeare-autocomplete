package word

import (
	"sort"
	"strings"
)

type Freq struct {
	Word  string
	Count int
}

func SortedByFreq(counts map[string]int) (freqs []Freq) {
	freqs = make([]Freq, len(counts))
	var idx int
	for w, c := range counts {
		freqs[idx] = Freq{Word: w, Count: c}
		idx++
	}
	sort.Sort(ByCount(freqs))
	return
}

type ByCount []Freq

func (c ByCount) Len() int      { return len(c) }
func (c ByCount) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ByCount) Less(i, j int) bool {
	// Most frequent first, alphabetical
	if c[i].Count == c[j].Count {
		return strings.Compare(c[i].Word, c[i].Word) < 0
	}
	return c[i].Count > c[j].Count
}
