package word

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"sync"
)

const (
	countChunkSize = 32 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// CountIn counts the words in a file or other io.Reader.
// This is overcomplicated for a single text file, but handles much larger
// files without spilling out into swap or getting OOM killed.
func CountIn(r io.Reader) (counts map[string]int, err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	pipeline := make(chan string)
	normalized := normalize(pipeline, &wg)
	countResult := count(normalized)

	// TODO: Use two sets of buffers and empty one while filling the other.

	// TODO: Double check that leftovers are handled correctly (maybe on text
	// where the words aren't so unusual)

	var leftover, combined []byte
	var lines, words [][]byte
	buf := make([]byte, countChunkSize)
	var n, m int
	for {
		n, err = r.Read(buf)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		if buf[n-1] == ' ' {
			leftover = []byte{}
		} else {
			for m = n - 1; buf[m] != ' ' && m > 0; m-- {
				leftover = buf[m:n]
				n = m - 1
			}
		}
		combined = append(leftover, buf[:n]...)
		lines = bytes.Split(combined, newline)
		for _, line := range lines {
			words = bytes.Split(line, space)
			for _, word := range words {
				if len(word) > 0 {
					pipeline <- string(word)
				}
			}
		}
	}

	close(pipeline)
	wg.Wait()
	close(normalized)
	counts = <-countResult

	return
}

var (
	passAlpha    = regexp.MustCompile(`^([a-z]*).*$`)
	passAlphaSub = []byte("$1")
)

// Normalize a word by converting to lowercase and removing trailing non-letters.
// This throws away words which start with a non-letter - for example, ordinals (1st, 2nd, etc) -
// which may not be intended.
func Normalize(word string) string {
	lower := strings.ToLower(word)
	norm := passAlpha.ReplaceAll([]byte(lower), passAlphaSub)
	// TODO: strip out roman numerals
	return string(norm)
}

func normalize(in <-chan string, wg *sync.WaitGroup) (out chan string) {
	out = make(chan string)

	go func() {
		for {
			word, ok := <-in
			if !ok {
				break
			}

			// PERF: inline Normalize if speed is critical
			norm := Normalize(word)
			if len(norm) > 0 {
				out <- string(norm)
			}
		}

		wg.Done()
	}()

	return
}

func count(words <-chan string) (result chan map[string]int) {
	result = make(chan map[string]int)

	go func() {
		counts := make(map[string]int)

		var count int
		for {
			word, ok := <-words
			if !ok {
				break
			}
			count = counts[word] + 1
			counts[word] = count
		}

		result <- counts
	}()

	return
}
