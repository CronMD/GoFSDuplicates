package sources

import (
	"df/internal/nodes"
	"iter"
	"math"
)

const tabSpaces = 4

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

		if spacesDiff == 0 {
			var parent *nodes.Node[string] = nil
			if len(f.siblings) > 0 {
				parent = f.siblings[len(f.siblings)-1].Parent

			}
			node.Parent = parent
			f.siblings = append(f.siblings, node)
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
		} else if spacesDiff < 0 {
			for _, node := range f.siblings {
				if !f.yield(node) {
					return nil
				}
			}

			parent := f.siblings[len(f.siblings)-1].Parent
			parentsLevel := -1 * spacesDiff
			for i := 0; i < parentsLevel; i++ {
				if parent != nil {
					parent = parent.Parent
				}
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
