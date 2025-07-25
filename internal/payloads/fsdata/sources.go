package fsdata

import (
	"context"
	"df/internal/nodes"
	"fmt"
	"io/fs"
	"iter"
	"log"
	"path/filepath"
	"sync"
)

type MultipleDirsFsDataSource struct {
	dirsPaths   []string
	failOnError bool
}

type multipleDirsFsDataSourceOption func(src *MultipleDirsFsDataSource)

func MulDirsDataSrcWithDirs(dirsPaths ...string) multipleDirsFsDataSourceOption {
	return func(src *MultipleDirsFsDataSource) {
		src.dirsPaths = make([]string, len(dirsPaths))
		copy(src.dirsPaths, dirsPaths)
	}
}

func MulDirsDataSrcWithDFailOnError(val bool) multipleDirsFsDataSourceOption {
	return func(src *MultipleDirsFsDataSource) {
		src.failOnError = val
	}
}

func NewMultipleDirsFsDataSource(
	opts ...multipleDirsFsDataSourceOption,
) *MultipleDirsFsDataSource {
	src := &MultipleDirsFsDataSource{
		dirsPaths:   nil,
		failOnError: false,
	}

	for _, opt := range opts {
		opt(src)
	}

	return src
}

func (s *MultipleDirsFsDataSource) Leafs() iter.Seq2[*nodes.Node[FsData], error] {
	dirsNumber := len(s.dirsPaths)
	var dirsWg sync.WaitGroup
	dirsWg.Add(dirsNumber)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var dirsErr error = nil

	dirsLeafs := make([][]*nodes.Node[FsData], dirsNumber)
	for i := 0; i < dirsNumber; i++ {
		go func(num int) {
			dirsLeafs[num], dirsErr = s.readDir(ctx, s.dirsPaths[num])
			dirsWg.Done()
			if dirsErr != nil {
				cancel()
				log.Println(s.dirsPaths[i], "[fail]")
				return
			}

			log.Println(s.dirsPaths[i], "[done]")
		}(i)
	}
	dirsWg.Wait()

	return func(yield func(*nodes.Node[FsData], error) bool) {
		if dirsErr != nil {
			yield(nil, dirsErr)
			return
		}

		for _, leafs := range dirsLeafs {
			for _, leaf := range leafs {
				if !yield(leaf, nil) {
					return
				}
			}
		}
	}
}

func (s *MultipleDirsFsDataSource) readDir(ctx context.Context, dirPath string) ([]*nodes.Node[FsData], error) {
	if dirPath[len(dirPath)-1] != filepath.Separator {
		dirPath = fmt.Sprintf("%s%c", dirPath, filepath.Separator)
	}

	parents := map[string]*nodes.Node[FsData]{}
	leafs := make([]*nodes.Node[FsData], 0)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s [interrupted]", dirPath)
		default:
			break
		}

		if err != nil {
			log.Println(err)

			if s.failOnError {
				return err
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			log.Println(err)

			if s.failOnError {
				return err
			}
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

	if err != nil {
		return nil, err
	}

	for _, node := range leafs {
		size := node.Payload.Size
		for node.Parent != nil {
			node.Parent.Payload.Size += size
			node = node.Parent
		}
	}

	return leafs, nil
}
