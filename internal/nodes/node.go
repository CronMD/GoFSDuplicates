package nodes

import (
	"fmt"
	"iter"
)

type Node[T any] struct {
	Payload T
	Parent  *Node[T]
}

func (n *Node[T]) String() string {
	parentPayload := ""
	if n.Parent != nil {
		parentPayload = fmt.Sprint(n.Parent.Payload)
	}

	return fmt.Sprintf("{%v - %s}", n.Payload, parentPayload)
}

type Source[T any] interface {
	Leafs() iter.Seq[*Node[T]]
}
