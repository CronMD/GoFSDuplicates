package sources

import (
	"df/internal/nodes"
	"iter"
	"math"
)

type txtParseState int

const tabSpaces = 4

const (
	stateEntry txtParseState = iota
	stateSearch
)

type entry struct {
	spaces  int
	payload string
}

type TextNoodesSource struct {
	txt []rune
}

func NewTextNoodesSource(txt string) *TextNoodesSource {
	return &TextNoodesSource{
		txt: []rune(txt),
	}
}

func (s *TextNoodesSource) Leafs() iter.Seq[*nodes.Node[string]] {
	return func(yield func(*nodes.Node[string]) bool) {
		prevSpaces := 0
		var prevNode *nodes.Node[string] = nil
		siblings := []*nodes.Node[string]{}

		for entry := range s.entries() {
			if prevNode == nil {
				prevSpaces = entry.spaces
			}

			node := &nodes.Node[string]{
				Payload: entry.payload,
				Parent:  nil,
			}

			spacesDiff := int(math.Ceil(float64(entry.spaces-prevSpaces) / float64(tabSpaces)))

			// fmt.Println(spacesDiff)
			if spacesDiff == 0 {
				if prevNode != nil {
					node.Parent = prevNode.Parent
				}

				siblings = append(siblings, node)
			} else if spacesDiff > 0 {
				node.Parent = prevNode
				siblings = []*nodes.Node[string]{node}
			} else if spacesDiff < 0 {
				parent := prevNode.Parent
				for i := 0; i < spacesDiff; i++ {
					parent = parent.Parent
				}

				node.Parent = parent

				for _, sibling := range siblings {
					if !yield(sibling) {
						return
					}
				}

				clear(siblings)
			}

			prevSpaces = entry.spaces
			prevNode = node
		}

		for _, sibling := range siblings {
			if !yield(sibling) {
				return
			}
		}
	}
}

func (s *TextNoodesSource) entries() iter.Seq[entry] {
	return func(yield func(entry) bool) {
		curSpaces := 0
		curPayload := []rune{}

		state := stateSearch
		pos := 0
		for {
			var cur rune = 0
			if pos < len(s.txt) {
				cur = s.txt[pos]
			}

			switch state {
			case stateSearch:
				switch cur {
				case '\n':
					curSpaces = 0
					pos++
				case ' ':
					curSpaces++
					pos++
				case '\t':
					curSpaces += 4
					pos++
				case 0:
					return
				default:
					state = stateEntry
				}
			case stateEntry:
				switch cur {
				case '\n', 0:
					newEntry := entry{
						spaces:  curSpaces,
						payload: string(curPayload),
					}
					if !yield(newEntry) {
						return
					}
					curSpaces = 0
					curPayload = []rune{}
					pos++
					state = stateSearch
				default:
					curPayload = append(curPayload, cur)
					pos++
				}
			}
		}
	}
}

type TxtSource struct {
	txt string
}

func NewTxtSource(txt string) *TxtSource {
	return &TxtSource{
		txt: txt,
	}
}

func (s *TxtSource) Leafs() iter.Seq[*nodes.Node[string]] {

	return func(yield func(*nodes.Node[string]) bool) {
		(&txtFsa{
			txt:        []rune(s.txt),
			curPayload: make([]rune, 0),
			yield:      yield,
			siblings:   make([]*nodes.Node[string], 0),
		}).start()
	}
}

type state func() state

type txtFsa struct {
	txt         []rune
	spacesCount int
	prevSpaces  int
	pos         int
	curPayload  []rune
	yield       func(*nodes.Node[string]) bool
	siblings    []*nodes.Node[string]
}

func (f *txtFsa) start() {
	curState := f.search
	for curState != nil {
		curState = curState()
	}
}

func (f *txtFsa) search() state {
	switch f.cur() {
	case ' ':
		f.spacesCount++
		f.pos++
		return f.search
	case '\n':
		f.spacesCount = 0
		f.pos++
		return f.search
	case '\t':
		f.spacesCount += 4
		f.pos++
		return f.search
	case 0:
		for _, node := range f.siblings {
			if !f.yield(node) {
				return nil
			}
		}

		return nil
	default:
		return f.entry
	}
}

func (f *txtFsa) entry() state {
	cur := f.cur()
	switch cur {
	case '\n', 0:
		payload := string(f.curPayload)
		node := &nodes.Node[string]{
			Payload: payload,
		}

		spacesDiff := int(math.Ceil(float64(f.spacesCount-f.prevSpaces) / float64(tabSpaces)))

		// if f.prevSpaces == f.spacesCount {
		if spacesDiff == 0 {
			var parent *nodes.Node[string] = nil
			if len(f.siblings) > 0 {
				parent = f.siblings[len(f.siblings)-1].Parent

			}
			node.Parent = parent
			f.siblings = append(f.siblings, node)
			// } else if f.prevSpaces < f.spacesCount {
		} else if spacesDiff > 0 {
			var parent *nodes.Node[string] = nil
			if len(f.siblings) > 0 {
				for _, node := range f.siblings[0 : len(f.siblings)-1] {
					if !f.yield(node) {
						return nil
					}
				}
				parent = f.siblings[len(f.siblings)-1]

				for i := spacesDiff; i > 1 && parent != nil; spacesDiff-- {
					parent = parent.Parent
				}
			}

			node.Parent = parent
			f.siblings = []*nodes.Node[string]{node}
			// } else if f.prevSpaces > f.spacesCount {
		} else if spacesDiff < 0 {
			for _, node := range f.siblings {
				if !f.yield(node) {
					return nil
				}
			}

			parent := f.siblings[len(f.siblings)-1].Parent
			if parent != nil {
				parent = parent.Parent
			}
			node.Parent = parent

			f.siblings = []*nodes.Node[string]{node}
		}

		f.curPayload = []rune{}
		f.prevSpaces = f.spacesCount
		f.spacesCount = 0
		f.pos++
		return f.search
	default:
		f.curPayload = append(f.curPayload, cur)
		f.pos++
		return f.entry
	}
}

func (f *txtFsa) cur() rune {
	if f.pos >= len(f.txt) {
		return 0
	}

	return f.txt[f.pos]
}
