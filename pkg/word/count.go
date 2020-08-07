package word

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// CountIn counts the words in a file or other io.Reader.
// This is overcomplicated for a single text file, but handles much larger
// files without spilling out into swap or getting OOM killed.
func CountIn(r io.Reader) (counts map[string]int, err error) {
	pipeline := make(chan string, 1)
	normalized := normalize(pipeline)
	countResult := count(normalized)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		for _, word := range words {
			if len(word) > 0 {
				pipeline <- word
			}
		}
	}
	err = scanner.Err()

	close(pipeline)
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

func normalize(in <-chan string) (out chan string) {
	out = make(chan string, 1)

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

		close(out)
	}()

	return
}

func count(words <-chan string) (result chan map[string]int) {
	result = make(chan map[string]int, 1)

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
