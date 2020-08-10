package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
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

	var before, after runtime.MemStats

	runtime.ReadMemStats(&before)
	s := time.Now()
	var counts map[string]int
	counts, err = word.CountIn(f)
	d := time.Since(s)
	runtime.ReadMemStats(&after)
	if err != nil {
		panic(err)
	}
	fmt.Printf("counted in %s\n", d)
	fmt.Printf("heap diff %d\n", after.TotalAlloc-before.TotalAlloc)

	s = time.Now()
	t := trie.New(counts)
	d = time.Since(s)
	fmt.Printf("build trie in %s\n", d)

	runtime.ReadMemStats(&before)
	s = time.Now()
	matches := t.FindWithPrefix("t")
	d = time.Since(s)
	runtime.ReadMemStats(&after)
	fmt.Printf("found %d words in %s\n", len(matches), d)
	fmt.Printf("heap diff %d\n", after.TotalAlloc-before.TotalAlloc)

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

	runtime.ReadMemStats(&before)
	s = time.Now()
	data, err := trie.Marshal(t)
	d = time.Since(s)
	runtime.ReadMemStats(&after)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile("shakespeare.bin", data, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("serialized %d words in %s\n", t.Len(), d)
	fmt.Printf("heap diff %d\n", after.TotalAlloc-before.TotalAlloc)

	runtime.ReadMemStats(&before)
	s = time.Now()
	t2, err := trie.Unmarshal(data)
	d = time.Since(s)
	runtime.ReadMemStats(&after)
	if err != nil {
		panic(err)
	}
	fmt.Printf("deserialized %d words in %s\n", t2.Len(), d)
	fmt.Printf("heap diff %d\n", after.TotalAlloc-before.TotalAlloc)
}
