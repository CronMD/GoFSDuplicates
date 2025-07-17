package main

import (
	humansize "df/cmd/find-dups/human-size"
	"df/internal/indexers"
	"df/internal/nodes"
	"df/internal/sources"
	"flag"
	"fmt"
	"slices"
	"strings"
)

func main() {
	dirsParam := flag.String("dirs", "", "comma separated dir paths")
	flag.Parse()

	dirs := make([]string, 0)
	for _, dir := range strings.Split(strings.Trim(*dirsParam, " "), ",") {
		dir = strings.Trim(dir, " ")
		if dir != "" {
			dirs = append(dirs, dir)
		}
	}

	src := sources.NewMultipleDirsFsDataSource(dirs...)

	finder := nodes.NewDupFinder([]nodes.Indexer[sources.FsData]{
		indexers.NewNameSizeFsIndexer(),
	})

	duplicates := finder.FindFromSources(src)
	slices.SortFunc(duplicates, func(nodes1, nodes2 []*nodes.Node[sources.FsData]) int {
		return int(nodes1[0].Payload.Size) - int(nodes2[0].Payload.Size)
	})
	for _, nodes := range duplicates {
		if len(nodes) < 2 {
			continue
		}

		if nodes[0].Payload.IsFile {
			fmt.Println("File:")
		} else {
			fmt.Println("Dir:")
		}

		for _, node := range nodes {
			fmt.Printf(
				"\t%s %s\n",
				node.Payload.Path,
				humansize.SizeToString(node.Payload.Size))
		}

		fmt.Println()
	}
}
