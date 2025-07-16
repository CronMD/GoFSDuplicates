package sources

import (
	"df/internal/nodes"
	"fmt"
	"io/fs"
	"iter"
	"log"
	"path/filepath"
	"sync"
)

type FsData struct {
	Path   string
	Name   string
	Size   int64
	IsFile bool
}

type MultipleDirsFsDataSource struct {
	dirsPaths []string
}

func NewMultipleDirsFsDataSource(dirsPaths ...string) *MultipleDirsFsDataSource {
	return &MultipleDirsFsDataSource{
		dirsPaths: dirsPaths,
	}
}

func (s *MultipleDirsFsDataSource) Leafs() iter.Seq[*nodes.Node[FsData]] {
	dirsNumber := len(s.dirsPaths)
	var dirsWg sync.WaitGroup

	dirsWg.Add(dirsNumber)
	dirsLeafs := make([][]*nodes.Node[FsData], dirsNumber)
	for i := 0; i < dirsNumber; i++ {
		go func(num int) {
			dirsLeafs[num] = s.readDir(s.dirsPaths[num])
			dirsWg.Done()
		}(i)
	}
	dirsWg.Wait()

	return func(yield func(*nodes.Node[FsData]) bool) {
		for _, leafs := range dirsLeafs {
			for _, leaf := range leafs {
				if !yield(leaf) {
					return
				}
			}
		}
	}
}

func (s *MultipleDirsFsDataSource) readDir(dirPath string) []*nodes.Node[FsData] {
	if dirPath[len(dirPath)-1] != filepath.Separator {
		dirPath = fmt.Sprintf("%s%c", dirPath, filepath.Separator)
	}

	parents := map[string]*nodes.Node[FsData]{}
	leafs := make([]*nodes.Node[FsData], 0)
	filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		info, err := d.Info()
		if err != nil {
			log.Println(err)
			return nil
		}

		newNode := &nodes.Node[FsData]{
			Payload: FsData{
				Path:   path,
				Name:   info.Name(),
				Size:   info.Size(),
				IsFile: !info.IsDir(),
			},
		}

		parentPath := filepath.Dir(path)
		if parent, ok := parents[parentPath]; ok {
			newNode.Parent = parent
		}

		parents[path] = newNode

		if newNode.Payload.IsFile {
			leafs = append(leafs, newNode)
		}

		return nil
	})

	return leafs
}
