package trie

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// A Trie is a prefix tree.
type Trie struct {
	len      int
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

	len := len(freqs)
	return &Trie{len, children}
}

// FindWithPrefix all matching words and their counts.
func (t *Trie) FindWithPrefix(p string) (counts map[string]int) {
	counts = make(map[string]int)
	charIdx := p[0] - 'a'
	t.children[charIdx].findWithPrefix(p, &counts)
	return
}

func (t *Trie) Len() int {
	if t == nil {
		return 0
	}
	return t.len
}

func Marshal(t *Trie) (data []byte, err error) {
	if t == nil {
		err = errors.New("trie is nil")
		return
	}

	var buf bytes.Buffer

	// TODO: add file header
	var childrenPopulated uint32
	for i, c := range t.children {
		if c != nil {
			childrenPopulated |= 0x1 << i
		}
	}
	childrenPopBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(childrenPopBytes, childrenPopulated)
	buf.Write(childrenPopBytes)

	for _, c := range t.children {
		var cData []byte
		cData, err = marshal(c)
		if err != nil {
			return
		}
		buf.Write(cData)
	}

	data = buf.Bytes()
	return
}

func Unmarshal(data []byte) (t *Trie, err error) {
	t = &Trie{}
	// TODO: read file header
	childrenPopBytes := data[:4]
	childrenPopulated := binary.LittleEndian.Uint32(childrenPopBytes)
	offset := 4
	for i := range t.children {
		if childrenPopulated&(0x1<<i) != 0 {
			t.children[i], offset, err = unmarshal(data, offset)
			if err != nil {
				return
			}
			t.len += t.children[i].len()
		}
	}
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

func (n *node) len() int {
	if n == nil {
		return 0
	}

	var l int
	if n.count != 0 {
		l = 1
	}
	for _, c := range n.children {
		l += c.len()
	}
	return l
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

func marshal(n *node) (data []byte, err error) {
	var buf bytes.Buffer

	if n != nil {
		if n.count != 0 {
			lBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(lBytes, uint16(len(n.word)))
			buf.Write(lBytes)

			buf.WriteString(n.word)

			cBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(cBytes, uint32(n.count))
			buf.Write(cBytes)
		} else {
			buf.Write([]byte{0, 0})
		}

		var childrenPopulated uint32
		for i, c := range n.children {
			if c != nil {
				childrenPopulated |= 0x1 << i
			}
		}
		childrenPopBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(childrenPopBytes, childrenPopulated)
		buf.Write(childrenPopBytes)

		for _, c := range n.children {
			if c != nil {
				var cBytes []byte
				cBytes, err = marshal(c)
				if err != nil {
					return
				}
				buf.Write(cBytes)
			}
		}
	}

	data = buf.Bytes()
	return
}

func unmarshal(data []byte, offset int) (n *node, newOffset int, err error) {
	n = &node{}

	newOffset = offset + 2
	lBytes := data[offset:newOffset]
	l := binary.LittleEndian.Uint16(lBytes)
	if l > 0 {
		offset = newOffset
		newOffset += int(l)
		wBytes := data[offset:newOffset]
		n.word = string(wBytes)

		offset = newOffset
		newOffset += 4
		cBytes := data[offset:newOffset]
		n.count = int(binary.LittleEndian.Uint32(cBytes))
	}

	offset = newOffset
	newOffset += 4
	childrenPopBytes := data[offset:newOffset]
	childrenPopulated := binary.LittleEndian.Uint32(childrenPopBytes)
	for i := range n.children {
		if childrenPopulated&(0x1<<i) != 0 {
			n.children[i], newOffset, err = unmarshal(data, newOffset)
			if err != nil {
				return
			}
		}
	}

	return
}
