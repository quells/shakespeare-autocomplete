package trie

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

// FindWithPrefix all matching words and their counts.
func (t *Trie) FindWithPrefix(p string) (counts map[string]int) {
	counts = make(map[string]int)
	charIdx := p[0] - 'a'
	t.children[charIdx].findWithPrefix(p, &counts)
	return
}

type node struct {
	depth    int
	word     string
	count    int
	children [26]*node
}

func newNode(depth int) *node {
	word := ""
	count := 0
	children := [26]*node{}
	return &node{depth, word, count, children}
}

func (n *node) insert(word string, count int, remaining string) {
	if remaining == "" {
		n.word = word
		n.count = count
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

	if len(p) <= n.depth {
		if n.count > 0 {
			(*m)[n.word] = n.count
		}
		for _, c := range n.children {
			c.findWithPrefix(p, m)
		}
		return
	}

	idx := p[n.depth] - 'a'
	n.children[idx].findWithPrefix(p, m)
}
