package main

import (
	"fmt"
	"os"

	"github.com/quells/shakespeare-autocomplete/pkg/word"
)

func main() {
	f, err := os.Open("shakespeare.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var counts map[string]int
	counts, err = word.CountIn(f)
	if err != nil {
		panic(err)
	}

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
