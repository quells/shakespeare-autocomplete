package trie

import (
	"fmt"
	"strings"
)

// A Trie is a prefix tree.
type Trie struct {
	children [26]*node
}

func New(freqs map[string]int) *Trie {
	children := [26]*node{}

	for word, count := range freqs {
		idx := word[0] - 'a'
		child := children[idx]
		if child == nil {
			child = newNode(1)
			children[idx] = child
		}
		child.insert(word, count, word[1:])
	}

	return &Trie{children}
}

// String serializes a Trie to a String - beware, this gets large fast.
func (t *Trie) String() string {
	children := make([]string, 26)
	for i, c := range t.children {
		children[i] = c.String()
	}
	return "[" + strings.Join(children, " ") + "]"
}

// FindWithPrefix all matching words and their counts.
func (t *Trie) FindWithPrefix(p string) (counts map[string]int) {
	counts = make(map[string]int)
	charIdx := p[0] - 'a'
	t.children[charIdx].findWithPrefix(p, &counts)
	return
}

type node struct {
	depth    int
	words    map[string]int // word : count
	children [26]*node
}

func newNode(depth int) *node {
	words := make(map[string]int)
	children := [26]*node{}
	return &node{depth, words, children}
}

func (n *node) String() string {
	if n == nil {
		return "[]"
	}

	children := make([]string, 26)
	for i, c := range n.children {
		children[i] = c.String()
	}
	return "[ " + fmt.Sprintf("%v", n.words) + strings.Join(children, " ") + "]"
}

func (n *node) insert(word string, count int, remaining string) {
	if remaining == "" {
		n.words[word] = count
		return
	}

	idx := remaining[0] - 'a'
	c := n.children[idx]
	if c == nil {
		c = newNode(n.depth + 1)
		n.children[idx] = c
	}
	c.insert(word, count, remaining[1:])
}

func (n *node) findWithPrefix(p string, m *map[string]int) {
	if n == nil {
		return
	}

	for w, c := range n.words {
		if strings.HasPrefix(w, p) {
			(*m)[w] = c
		}
	}

	if len(p) <= n.depth {
		for _, c := range n.children {
			c.findWithPrefix(p, m)
		}
		return
	}

	idx := p[n.depth] - 'a'
	n.children[idx].findWithPrefix(p, m)
}
