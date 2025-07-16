package nodes

import (
	"slices"
	"strings"
	"testing"
)

func TestEmpty(t *testing.T) {
	finder := newTestFinder()

	result := finder.FindFromLeafs()

	if len(result) != 0 {
		t.Fatal("0 size exppected, got", len(result))
	}
}

func TestSameFilesAtRoot(t *testing.T) {
	finder := newTestFinder()
	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			f1
			f2
			f1
			f2
			`),
		),
		[][]*Node[string]{
			{
				&Node[string]{
					Payload: "f1",
					Parent:  nil,
				},
				&Node[string]{
					Payload: "f1",
					Parent:  nil,
				},
			},

			{
				&Node[string]{
					Payload: "f2",
					Parent:  nil,
				},
				&Node[string]{
					Payload: "f2",
					Parent:  nil,
				},
			},
		},
	)
}

func TestSameDirs(t *testing.T) {
	finder := newTestFinder()

	dir1 := &Node[string]{
		Payload: "d1",
		Parent:  nil,
	}

	dir2 := &Node[string]{
		Payload: "d2",
		Parent:  nil,
	}

	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			d1
				f1
				f2
			d2
				f1
				f2
			`),
		),
		[][]*Node[string]{
			{
				dir1,
				dir2,
			},
		},
	)
}

type equalPayloadIndexer struct{}

func (i *equalPayloadIndexer) Index(node *Node[string]) interface{} {
	return node.Payload
}

func newTestFinder() *DupFinder[string] {
	return NewDupFinder([]Indexer[string]{
		&equalPayloadIndexer{},
	})
}

func assertNodes(t *testing.T, actual [][]*Node[string], expected [][]*Node[string]) {
	fail := func() {
		t.Fatal("actual", actual, "does not match", expected)
	}

	cmpNodes := func(node1, node2 *Node[string]) int {
		diff := strings.Compare(node1.Payload, node2.Payload)
		for diff == 0 {
			node1 = node1.Parent
			node2 = node2.Parent

			if node1 == nil || node2 == nil {
				break
			}

			diff = strings.Compare(node1.Payload, node2.Payload)
		}

		return diff
	}

	cmpDups := func(nodes1, nodes2 []*Node[string]) int {
		lenDiff := len(nodes1) - len(nodes2)
		if lenDiff != 0 {
			return lenDiff
		}

		slices.SortFunc(nodes1, cmpNodes)
		slices.SortFunc(nodes2, cmpNodes)

		for i, node := range nodes1 {
			nodesDiff := cmpNodes(node, nodes2[i])
			if nodesDiff != 0 {
				return nodesDiff
			}
		}

		return 0
	}

	slices.SortFunc(actual, cmpDups)
	slices.SortFunc(expected, cmpDups)

	if len(actual) != len(expected) {
		fail()
		return
	}

	for i, actualNodes := range actual {
		for j, actualNode := range actualNodes {
			if cmpNodes(actualNode, expected[i][j]) != 0 {
				fail()
				return
			}
		}
	}
}
