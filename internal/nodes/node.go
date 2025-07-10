package nodes

import (
	"iter"
)

type Node[T any] struct {
	Payload T
	Parent  *Node[T]
}

// func (n *Node[T]) String() string {
// 	parentPayload := ""
// 	if n.Parent != nil {
// 		parentPayload = fmt.Sprint(n.Parent.Payload)
// 	}

// 	return fmt.Sprintf("%s - %s", string(n.Payload), parentPayload)
// }

type Source[T any] interface {
	Leafs() iter.Seq[*Node[T]]
}
