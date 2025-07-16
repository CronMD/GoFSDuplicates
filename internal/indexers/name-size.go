package indexers

import (
	"df/internal/nodes"
	"df/internal/sources"
)

type NameSizeFsIndexer struct{}

func NewNameSizeFsIndexer() *NameSizeFsIndexer {
	return &NameSizeFsIndexer{}
}

func (ix *NameSizeFsIndexer) Index(node *nodes.Node[sources.FsData]) interface{} {
	return struct {
		name string
		size int64
	}{
		node.Payload.Name,
		node.Payload.Size,
	}
}
