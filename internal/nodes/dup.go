package nodes

import (
	"log"
	"slices"
)

type Indexer[T any] interface {
	Index(node *Node[T]) (interface{}, error)
}

type DupFinder[T any] struct {
	indexers []Indexer[T]
}

func NewDupFinder[T any](indexers []Indexer[T]) *DupFinder[T] {
	return &DupFinder[T]{
		indexers: indexers,
	}
}

func (f *DupFinder[T]) FindFromSources(srcs ...Source[T]) ([][]*Node[T], error) {
	leafs := make([]*Node[T], 0)
	for _, src := range srcs {
		for leaf := range src.Leafs() {
			leafs = append(leafs, leaf)
		}
	}

	return f.FindFromLeafs(leafs...)
}

func (f *DupFinder[T]) FindFromLeafs(leafs ...*Node[T]) ([][]*Node[T], error) {
	log.Println("got", len(leafs))

	indexedNodes, err := f.groupByIndexes(leafs...)
	if err != nil {
		return nil, err
	}

	return f.mergeParents(indexedNodes), nil
}

func (f *DupFinder[T]) groupByIndexes(leafs ...*Node[T]) ([][]*Node[T], error) {
	result := [][]*Node[T]{leafs}
	for _, indexer := range f.indexers {
		newResult := make([][]*Node[T], 0)
		for _, nodes := range result {
			groups := make(map[interface{}][]*Node[T])
			for _, node := range nodes {
				index, err := indexer.Index(node)
				if err != nil {
					return nil, err
				}

				if _, ok := groups[index]; !ok {
					groups[index] = make([]*Node[T], 0)
				}
				groups[index] = append(groups[index], node)
			}

			for _, groupNodes := range groups {
				newResult = append(newResult, groupNodes)
			}
		}

		result = newResult
	}

	return result, nil
}

func (f *DupFinder[T]) mergeParents(dups [][]*Node[T]) [][]*Node[T] {
	type dupsPath struct {
		path     []int
		pathSum  int
		node     *Node[T]
		children map[*dupsPath]bool
	}
	pathsMap := make(map[*Node[T]]*dupsPath)
	paths := make([]*dupsPath, 0)
	topmost := make(map[*dupsPath]bool)
	result := make([][]*Node[T], 0)

	dupsCount := 1
	for _, nodes := range dups {
		if len(nodes) < 2 {
			result = append(result, nodes)
		}

		fringe := make([]*dupsPath, len(nodes))
		for i, node := range nodes {
			curPath := &dupsPath{
				path:     []int{dupsCount},
				children: map[*dupsPath]bool{},
				node:     node,
				pathSum:  dupsCount,
			}
			pathsMap[node] = curPath
			paths = append(paths, curPath)
			topmost[curPath] = true

			fringe[i] = curPath
		}

		for len(fringe) > 0 {
			prevPath := fringe[0]
			fringe = fringe[1:]

			if prevPath.node.Parent == nil {
				continue
			}

			curPath, ok := pathsMap[prevPath.node.Parent]
			if !ok {
				curPath = &dupsPath{
					path:     []int{-1},
					pathSum:  -1,
					node:     prevPath.node.Parent,
					children: map[*dupsPath]bool{},
				}

				pathsMap[curPath.node] = curPath
				paths = append(paths, curPath)
			}

			curPath.path = append(curPath.path, dupsCount)
			curPath.pathSum += dupsCount
			curPath.children[prevPath] = true

			delete(topmost, prevPath)
			topmost[curPath] = true

			fringe = append(fringe, curPath)
		}

		dupsCount++
	}

	cmpDupPaths := func(p1, p2 *dupsPath) int {
		sumDiff := p1.pathSum - p2.pathSum
		if sumDiff != 0 {
			return sumDiff
		}

		lenDiff := len(p1.path) - len(p2.path)
		if lenDiff != 0 {
			return lenDiff
		}

		for i, dupNum := range p1.path {
			dupDiff := dupNum - p2.path[i]
			if dupDiff != 0 {
				return dupDiff
			}
		}

		return 0
	}

	slices.SortFunc(paths, cmpDupPaths)

	groupedPaths := make(map[*dupsPath][]*dupsPath)
	curGrop := []*dupsPath{}
	for _, curPath := range paths {
		if len(curGrop) == 0 {
			curGrop = []*dupsPath{curPath}
		} else {
			if cmpDupPaths(curPath, curGrop[0]) == 0 {
				curGrop = append(curGrop, curPath)
			} else {
				curGrop = []*dupsPath{curPath}
			}
		}

		for _, path := range curGrop {
			groupedPaths[path] = curGrop
		}
	}

	fringe := make([]*dupsPath, 0)
	for path := range topmost {
		fringe = append(fringe, path)
	}

	for len(fringe) > 0 {
		path := fringe[0]
		fringe = fringe[1:]

		group, ok := groupedPaths[path]
		if !ok {
			continue
		}

		if len(group) > 1 {
			parents := make(map[*Node[T]]bool)
			for _, groupPath := range group {
				parents[groupPath.node] = true
			}

			groupResult := []*Node[T]{}
			for _, groupPath := range group {
				if _, ok := parents[groupPath.node.Parent]; !ok {
					groupResult = append(groupResult, groupPath.node)
				}
			}

			result = append(result, groupResult)

			for _, groupPath := range group {
				delete(groupedPaths, groupPath)
			}
			continue
		}

		for child := range path.children {
			fringe = append(fringe, child)
		}
	}

	return result
}
