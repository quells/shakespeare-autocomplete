package main_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/quells/shakespeare-autocomplete/pkg/trie"
	"github.com/quells/shakespeare-autocomplete/pkg/word"
)

func BenchmarkCountIn(b *testing.B) {
	data, err := ioutil.ReadFile("../../shakespeare.txt")
	if err != nil {
		b.Fatal(err)
	}
	buf := bytes.NewBuffer(data)
	var counts map[string]int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counts, _ = word.CountIn(buf)
		if len(counts) != 23318 {
			b.Fatalf("%d", len(counts))
		}
	}
}

func BenchmarkTrieNew(b *testing.B) {
	data, err := ioutil.ReadFile("../../shakespeare.txt")
	if err != nil {
		b.Fatal(err)
	}
	buf := bytes.NewBuffer(data)
	counts, _ := word.CountIn(buf)
	var t *trie.Trie

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t = trie.New(counts)
		if t.Len() != 23318 {
			b.Fatalf("%d", len(counts))
		}
	}
}

func BenchmarkTrieFindWithPrefix(b *testing.B) {
	data, err := ioutil.ReadFile("../../shakespeare.txt")
	if err != nil {
		b.Fatal(err)
	}
	buf := bytes.NewBuffer(data)
	counts, _ := word.CountIn(buf)
	t := trie.New(counts)

	b.Run("a", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			results := t.FindWithPrefix("a")
			if len(results) != 1352 {
				b.Fatalf("%d", len(results))
			}
		}
	})

	b.Run("th", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			results := t.FindWithPrefix("th")
			if len(results) != 220 {
				b.Fatalf("%d", len(results))
			}
		}
	})

	b.Run("wu", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			results := t.FindWithPrefix("wul")
			if len(results) != 1 {
				b.Fatalf("%d", len(results))
			}
		}
	})

	b.Run("aardvark", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			results := t.FindWithPrefix("aardvark")
			if len(results) != 0 {
				b.Fatalf("%d", len(results))
			}
		}
	})
}
