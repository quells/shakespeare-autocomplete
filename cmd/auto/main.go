package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/quells/shakespeare-autocomplete/pkg/trie"
	"github.com/quells/shakespeare-autocomplete/pkg/word"
)

/*
TODO:
- read port from environment variable
- move routes into own file (if there were more)
- serialize trie to disk to amortize build time (only use if the input filepath matches)
*/

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: auto <input>")
		os.Exit(1)
	}

	log.Println("starting to read file")
	auto, err := buildAutoHandler(os.Args[1])
	if err != nil {
		fmt.Printf("could not build trie: %v\n", err)
		os.Exit(1)
	}

	http.Handle("/autocomplete", auto)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Println("listening on :5000")
	if err = http.ListenAndServe(":5000", nil); err != nil {
		fmt.Println(err)
	}
}

type autoHandler struct {
	t *trie.Trie
}

func buildAutoHandler(filename string) (h autoHandler, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	var counts map[string]int
	counts, err = word.CountIn(f)
	f.Close()
	if err != nil {
		return
	}

	t := trie.New(counts)
	h = autoHandler{t}

	return
}

func (h autoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	term := word.Normalize(q.Get("term"))
	if term == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing or invalid `term` query parameter"))
		return
	}

	results := word.SortedByFreq(h.t.FindWithPrefix(term))
	if len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("no results found"))
		return
	}

	if len(results) > 25 {
		results = results[:25]
	}

	for _, r := range results {
		w.Write([]byte(r.Word))
		w.Write([]byte{'\n'})
	}
}
