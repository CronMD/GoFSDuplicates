package nodes

import (
	"fmt"
	"iter"
	"testing"
)

func TestParseEmpty(t *testing.T) {
	txt := `
	`
	leafsSlice := make([]*Node[string], 0)
	for leaf := range NewTxtSource(txt).Leafs() {
		leafsSlice = append(leafsSlice, leaf)
	}

	if len(leafsSlice) != 0 {
		t.Fatal("expected", 0, "leafs, got", len(leafsSlice))
	}
}

func TestParseToplevelFiles(t *testing.T) {
	txt := `
		f1
		f2
	`
	assertLeafs(
		t,
		NewTxtSource(txt).Leafs(),
		&Node[string]{
			Payload: "f1",
		},
		&Node[string]{
			Payload: "f2",
		},
	)
}

func TestParseDir(t *testing.T) {
	txt := `
	 	d1
			f1
			f2
	`

	d1 := &Node[string]{
		Payload: "d1",
	}

	assertLeafs(
		t,
		NewTxtSource(txt).Leafs(),
		&Node[string]{
			Payload: "f1",
			Parent:  d1,
		},
		&Node[string]{
			Payload: "f2",
			Parent:  d1,
		},
	)
}

func TestSameLevelDir(t *testing.T) {
	txt := `
	 	d1
			f1
			f2
		d2
			f3
			f4
	`

	d1 := &Node[string]{
		Payload: "d1",
	}

	d2 := &Node[string]{
		Payload: "d2",
	}

	assertLeafs(
		t,
		NewTxtSource(txt).Leafs(),
		&Node[string]{
			Payload: "f1",
			Parent:  d1,
		},
		&Node[string]{
			Payload: "f2",
			Parent:  d1,
		},
		&Node[string]{
			Payload: "f3",
			Parent:  d2,
		},
		&Node[string]{
			Payload: "f4",
			Parent:  d2,
		},
	)
}

func TestMultiLevel(t *testing.T) {
	txt := `
	 	d1
			f1
			f2
		d2
			d3
				f3
				f4
			d4
				f1
				f2
	`

	d1 := &Node[string]{
		Payload: "d1",
	}

	d2 := &Node[string]{
		Payload: "d2",
	}

	d3 := &Node[string]{
		Payload: "d3",
		Parent:  d2,
	}

	d4 := &Node[string]{
		Payload: "d4",
		Parent:  d2,
	}

	assertLeafs(
		t,
		NewTxtSource(txt).Leafs(),
		&Node[string]{
			Payload: "f1",
			Parent:  d1,
		},
		&Node[string]{
			Payload: "f2",
			Parent:  d1,
		},
		&Node[string]{
			Payload: "f3",
			Parent:  d3,
		},
		&Node[string]{
			Payload: "f4",
			Parent:  d3,
		},
		&Node[string]{
			Payload: "f1",
			Parent:  d4,
		},
		&Node[string]{
			Payload: "f2",
			Parent:  d4,
		},
	)
}

func TestMultiLevel2(t *testing.T) {
	txt := `
	 	d1
			d2
				f3
				d3
					f4
			d4
				f1
				f2
			f1
			f2
	`

	d1 := &Node[string]{
		Payload: "d1",
	}

	d2 := &Node[string]{
		Payload: "d2",
		Parent:  d1,
	}

	d3 := &Node[string]{
		Payload: "d3",
		Parent:  d2,
	}

	d4 := &Node[string]{
		Payload: "d4",
		Parent:  d1,
	}

	assertLeafs(
		t,
		NewTxtSource(txt).Leafs(),
		&Node[string]{
			Payload: "f3",
			Parent:  d2,
		},
		&Node[string]{
			Payload: "f4",
			Parent:  d3,
		},
		&Node[string]{
			Payload: "f1",
			Parent:  d4,
		},
		&Node[string]{
			Payload: "f2",
			Parent:  d4,
		},
		&Node[string]{
			Payload: "f1",
			Parent:  d1,
		},
		&Node[string]{
			Payload: "f2",
			Parent:  d1,
		},
	)
}

func assertLeafs(
	t *testing.T,
	leafsSeq iter.Seq[*Node[string]],
	expected ...*Node[string]) {
	actual := []*Node[string]{}
	for leaf := range leafsSeq {
		actual = append(actual, leaf)
	}

	if len(actual) != len(expected) {
		fail(t, actual, expected)
		return
	}

	for i, actualNode := range actual {
		if actualNode.Payload != expected[i].Payload {
			fail(t, actual, expected)
			return
		}

		if actualNode.Parent == nil && expected[i].Parent != nil {
			fail(t, actual, expected)
			return
		}

		if actualNode.Parent != nil && expected[i].Parent != nil {
			if actualNode.Parent.Payload != expected[i].Parent.Payload {
				fail(t, actual, expected)
				return
			}
		}
	}
}

func fail(t *testing.T, actual []*Node[string], expected []*Node[string]) {
	t.Fatal(
		"actual",
		nodes2txt(actual),
		"not equal to expected",
		nodes2txt(expected),
	)
}

func nodes2txt(pNodes []*Node[string]) string {
	sNodes := []Node[string]{}
	for _, pNode := range pNodes {
		sNodes = append(sNodes, *pNode)
	}

	return fmt.Sprint(sNodes)
}
