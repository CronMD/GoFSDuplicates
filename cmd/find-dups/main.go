package main

import (
	appparams "df/cmd/find-dups/app-params"
	humansize "df/cmd/find-dups/human-size"
	"df/internal/nodes"
	"df/internal/payloads/fsdata"
	"fmt"
	"log"
	"slices"
	"strings"
)

type AppParams interface {
	Dirs() []string
	UseHash() bool
	FailOnError() bool
	DirsOnly() bool
	DirsFilters() []string
}

var params AppParams = appparams.NewCmdLineParams()

func main() {
	src := fsdata.NewMultipleDirsFsDataSource(
		fsdata.MulDirsDataSrcWithDirs(params.Dirs()...),
		fsdata.MulDirsDataSrcWithDFailOnError(params.FailOnError()),
	)

	ixs := []nodes.Indexer[fsdata.FsData]{
		fsdata.NewNameSizeFsIndexer(),
	}
	if params.UseHash() {
		ixs = append(ixs, fsdata.NewHashFsIndexer(
			fsdata.WithSuppressHashIndexErrors(!params.FailOnError())))
	}

	finder := nodes.NewDupFinder(ixs)

	duplicates, err := finder.FindFromSources(src)
	if err != nil {
		log.Fatal(err)
	}

	slices.SortFunc(duplicates, func(nodes1, nodes2 []*nodes.Node[fsdata.FsData]) int {
		return int(nodes1[0].Payload.Size) - int(nodes2[0].Payload.Size)
	})
	var ttlDupsSizeSum int64 = 0
	for _, nodes := range duplicates {
		if len(nodes) < 2 {
			continue
		}

		if len(params.DirsFilters()) > 0 {
			show := false

			for _, node := range nodes {
				for _, filter := range params.DirsFilters() {
					if strings.Index(strings.ToLower(node.Payload.Path), filter) == 0 {
						show = true
						break
					}
				}
				if show {
					break
				}
			}

			if !show {
				continue
			}
		}

		if nodes[0].Payload.IsFile {
			if params.DirsOnly() {
				continue
			}

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
		ttlDupsSizeSum += nodes[0].Payload.Size * int64(len(nodes)-1)

		fmt.Println()
	}

	fmt.Printf(
		"Estimated duplicates size %s\n",
		humansize.SizeToString(ttlDupsSizeSum))
}
