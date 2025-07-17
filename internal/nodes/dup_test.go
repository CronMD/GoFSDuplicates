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

func TestTwoSameDirs(t *testing.T) {
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

func TestUniqueFiles(t *testing.T) {
	finder := newTestFinder()

	dir1 := &Node[string]{
		Payload: "d1",
		Parent:  nil,
	}

	dir2 := &Node[string]{
		Payload: "d2",
		Parent:  nil,
	}

	dir3 := &Node[string]{
		Payload: "d3",
		Parent:  nil,
	}

	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			d1
				f1
			d2
				f2
			d3
				f3
			`),
		),
		[][]*Node[string]{
			{
				&Node[string]{
					Payload: "f1",
					Parent:  dir1,
				},
			},

			{
				&Node[string]{
					Payload: "f2",
					Parent:  dir2,
				},
			},

			{
				&Node[string]{
					Payload: "f3",
					Parent:  dir3,
				},
			},
		},
	)
}

func TestEqualFiles(t *testing.T) {
	finder := newTestFinder()

	dir1 := &Node[string]{
		Payload: "d1",
		Parent:  nil,
	}

	dir2 := &Node[string]{
		Payload: "d2",
		Parent:  nil,
	}

	dir3 := &Node[string]{
		Payload: "d3",
		Parent:  nil,
	}

	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			d1
				f1
			d2
				f2
				f1
			d3
				f3
			`),
		),
		[][]*Node[string]{
			{
				&Node[string]{
					Payload: "f1",
					Parent:  dir1,
				},
				&Node[string]{
					Payload: "f1",
					Parent:  dir2,
				},
			},

			{
				&Node[string]{
					Payload: "f2",
					Parent:  dir2,
				},
			},

			{
				&Node[string]{
					Payload: "f3",
					Parent:  dir3,
				},
			},
		},
	)
}

func TestExcessFiles(t *testing.T) {
	finder := newTestFinder()

	dir1 := &Node[string]{
		Payload: "d1",
		Parent:  nil,
	}

	dir2 := &Node[string]{
		Payload: "d2",
		Parent:  nil,
	}

	dir3 := &Node[string]{
		Payload: "d3",
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
				f2
				f1
				f3
			d3
				f3
			`),
		),
		[][]*Node[string]{
			{
				&Node[string]{
					Payload: "f1",
					Parent:  dir1,
				},
				&Node[string]{
					Payload: "f1",
					Parent:  dir2,
				},
			},

			{
				&Node[string]{
					Payload: "f2",
					Parent:  dir1,
				},
				&Node[string]{
					Payload: "f2",
					Parent:  dir2,
				},
			},

			{
				&Node[string]{
					Payload: "f3",
					Parent:  dir2,
				},
				&Node[string]{
					Payload: "f3",
					Parent:  dir3,
				},
			},
		},
	)
}

func TestFilesCopies(t *testing.T) {
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
				f1
				f2
			d2
				f1
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

func TestNestedEqualDirs(t *testing.T) {
	finder := newTestFinder()

	topDir1 := &Node[string]{
		Payload: "td1",
		Parent:  nil,
	}

	topDir2 := &Node[string]{
		Payload: "td2",
		Parent:  nil,
	}

	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			td1
				d1
					f1
					f2
			td2
				d2
					f1
					f2
			`),
		),
		[][]*Node[string]{
			{
				topDir1,
				topDir2,
			},
		},
	)
}

func TestSeveralNestedDirs(t *testing.T) {
	finder := newTestFinder()

	topDir1 := &Node[string]{
		Payload: "td1",
		Parent:  nil,
	}

	topDir2 := &Node[string]{
		Payload: "td2",
		Parent:  nil,
	}

	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			td1
				d1
					f1
				d2
					f2
			
			td2
				d3
					f1
				d4
					f2
			`),
		),
		[][]*Node[string]{
			{
				topDir1,
				topDir2,
			},
		},
	)
}

func TestNestedNotEqualDirs2(t *testing.T) {
	finder := newTestFinder()

	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			td1
				d1
					f1
				d2
					f2
			
			td2
				d4
					f2
			`),
		),
		[][]*Node[string]{
			{
				&Node[string]{
					Payload: "d2",
					Parent: &Node[string]{
						Payload: "td1",
						Parent:  nil,
					},
				},
				&Node[string]{
					Payload: "td2",
					Parent:  nil,
				},
			},

			{
				&Node[string]{
					Payload: "f1",
					Parent: &Node[string]{
						Payload: "d1",
						Parent: &Node[string]{
							Payload: "td1",
							Parent:  nil,
						},
					},
				},
			},
		},
	)
}

func TestMltipleRootsDuplicatedDirs(t *testing.T) {
	finder := newTestFinder()

	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			td1
				d1
					f1
				f2
			td2
				d2
					f1
				f2
			`),
		),
		[][]*Node[string]{
			{
				&Node[string]{
					Payload: "td1",
					Parent:  nil,
				},
				&Node[string]{
					Payload: "td2",
					Parent:  nil,
				},
			},
		},
	)
}

func TestNestedNotEqualDirs1(t *testing.T) {
	finder := newTestFinder()

	assertNodes(
		t,
		finder.FindFromSources(
			NewTxtSource(`
			td1
				d1
					f1
				d2
					f2
			
			td2
				d3
					f1
			`),
		),
		[][]*Node[string]{
			{
				&Node[string]{
					Payload: "d1",
					Parent: &Node[string]{
						Payload: "td1",
						Parent:  nil,
					},
				},
				&Node[string]{
					Payload: "td2",
					Parent:  nil,
				},
			},

			{
				&Node[string]{
					Payload: "f2",
					Parent: &Node[string]{
						Payload: "d2",
						Parent: &Node[string]{
							Payload: "td1",
							Parent:  nil,
						},
					},
				},
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
