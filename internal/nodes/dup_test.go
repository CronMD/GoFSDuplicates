package nodes

import (
	"slices"
	"strings"
	"testing"
)

func TestEmpty(t *testing.T) {
	finder := newTestFinder()

	result, err := finder.FindFromLeafs()
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(result) != 0 {
		t.Fatal("0 size exppected, got", len(result))
	}
}

func TestSameFilesAtRoot(t *testing.T) {
	assertTxt(
		t,
		`
		f1
		f2
		f1
		f2
		`,
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
	dir1 := &Node[string]{
		Payload: "d1",
		Parent:  nil,
	}

	dir2 := &Node[string]{
		Payload: "d2",
		Parent:  nil,
	}

	assertTxt(
		t,
		`
		d1
			f1
			f2
		d2
			f1
			f2
		`,
		[][]*Node[string]{
			{
				dir1,
				dir2,
			},
		},
	)
}

func TestUniqueFiles(t *testing.T) {
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

	assertTxt(
		t,
		`
		d1
			f1
		d2
			f2
		d3
			f3
		`,
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

	assertTxt(
		t,
		`
		d1
			f1
		d2
			f2
			f1
		d3
			f3
		`,
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

	assertTxt(
		t,
		`
		d1
			f1
			f2
		d2
			f2
			f1
			f3
		d3
			f3
		`,
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
	dir1 := &Node[string]{
		Payload: "d1",
		Parent:  nil,
	}

	dir2 := &Node[string]{
		Payload: "d2",
		Parent:  nil,
	}

	assertTxt(
		t,
		`
		d1
			f1
			f1
			f2
		d2
			f1
			f1
			f2
		`,
		[][]*Node[string]{
			{
				dir1,
				dir2,
			},
		},
	)
}

func TestNestedEqualDirs(t *testing.T) {
	topDir1 := &Node[string]{
		Payload: "td1",
		Parent:  nil,
	}

	topDir2 := &Node[string]{
		Payload: "td2",
		Parent:  nil,
	}

	assertTxt(
		t,
		`
		td1
			d1
				f1
				f2
		td2
			d2
				f1
				f2
		`,
		[][]*Node[string]{
			{
				topDir1,
				topDir2,
			},
		},
	)
}

func TestSeveralNestedDirs(t *testing.T) {
	topDir1 := &Node[string]{
		Payload: "td1",
		Parent:  nil,
	}

	topDir2 := &Node[string]{
		Payload: "td2",
		Parent:  nil,
	}

	assertTxt(
		t,
		`
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
			`,
		[][]*Node[string]{
			{
				topDir1,
				topDir2,
			},
		},
	)
}

func TestNestedNotEqualDirs1(t *testing.T) {
	assertTxt(
		t,
		`
		td1
			d1
				f1
			d2
				f2
		
		td2
			d3
				f1
		`,
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

func TestNestedNotEqualDirs2(t *testing.T) {
	assertTxt(
		t,
		`
			td1
				d1
					f1
				d2
					f2
			
			td2
				d4
					f2
			`,
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
	assertTxt(
		t,
		`
			td1
				d1
					f1
				f2
			td2
				d2
					f1
				f2
			`,
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

func TestTwoEqualDirsOfFour(t *testing.T) {
	assertTxt(
		t,
		`
		d1
			f1
			f2
		d2
			f1
		d3
			f1
			f2
		d4
			f1
		`,
		[][]*Node[string]{
			{
				&Node[string]{
					Payload: "d1",
					Parent:  nil,
				},
				&Node[string]{
					Payload: "d3",
					Parent:  nil,
				},
			},

			{
				&Node[string]{
					Payload: "d2",
					Parent:  nil,
				},
				&Node[string]{
					Payload: "d4",
					Parent:  nil,
				},
			},
		},
	)
}

func TestDeepNestedEqual(t *testing.T) {
	assertTxt(
		t,
		`
			td1
				p1_3
					p1_2
						p1_1
							d1
								f1
				f2
			td2
				p2_3
					p2_2
						p2_1
							d2
								f1
				f2
			`,
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

type equalPayloadIndexer struct{}

func (i *equalPayloadIndexer) Index(node *Node[string]) (interface{}, error) {
	return node.Payload, nil
}

func newTestFinder() *DupFinder[string] {
	return NewDupFinder([]Indexer[string]{
		&equalPayloadIndexer{},
	})
}

func assertTxt(t *testing.T, txt string, expected [][]*Node[string]) {
	finder := newTestFinder()

	result, err := finder.FindFromSources(
		NewTxtSource(txt),
	)
	if err != nil {
		t.Fatal(err)
		return
	}

	assertNodes(t, result, expected)
}

func assertNodes(t *testing.T, actual [][]*Node[string], expected [][]*Node[string]) {
	fail := func() {
		t.Fatal("actual", actual, "does not match", expected)
	}

	cmpNodes := func(node1, node2 *Node[string]) int {
		for node1 != nil && node2 != nil {
			diff := strings.Compare(node1.Payload, node2.Payload)
			if diff != 0 {
				return diff
			}

			node1 = node1.Parent
			node2 = node2.Parent
		}

		return 0
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

	if len(actual) != len(expected) {
		fail()
		return
	}

	slices.SortFunc(actual, cmpDups)
	slices.SortFunc(expected, cmpDups)

	if len(actual) == 1 {
		slices.SortFunc(actual[0], cmpNodes)
		slices.SortFunc(expected[0], cmpNodes)
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
