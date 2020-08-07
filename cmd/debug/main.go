package main

import (
	"fmt"
	"os"
	"time"

	"github.com/quells/shakespeare-autocomplete/pkg/trie"
	"github.com/quells/shakespeare-autocomplete/pkg/word"
)

func main() {
	f, err := os.Open("shakespeare.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := time.Now()
	var counts map[string]int
	counts, err = word.CountIn(f)
	d := time.Since(s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("counted in %s\n", d)

	s = time.Now()
	t := trie.New(counts)
	d = time.Since(s)
	fmt.Printf("build trie in %s\n", d)

	s = time.Now()
	matches := t.FindWithPrefix("t")
	d = time.Since(s)
	fmt.Printf("found %d words in %s\n", len(matches), d)

	for i, f := range word.SortedByFreq(counts) {
		for _, c := range f.Word {
			if c < 'a' || 'z' < c {
				fmt.Printf("%v %d\n", []byte(f.Word), f.Count)
				break
			}
		}
		if len(f.Word) == 0 {
			fmt.Println(i, f)
		}
	}
}
