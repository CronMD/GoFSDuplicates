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

type resultFilter func(nodes []*nodes.Node[fsdata.FsData]) bool

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

	resultFilters := []resultFilter{
		func(nodes []*nodes.Node[fsdata.FsData]) bool {
			return len(nodes) >= 2
		},
	}
	if len(params.DirsFilters()) > 0 {
		resultFilters = append(resultFilters, func(nodes []*nodes.Node[fsdata.FsData]) bool {
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

			return show
		})
	}
	if params.DirsOnly() {
		resultFilters = append(resultFilters, func(nodes []*nodes.Node[fsdata.FsData]) bool {
			return !nodes[0].Payload.IsFile
		})
	}

	slices.SortFunc(duplicates, func(nodes1, nodes2 []*nodes.Node[fsdata.FsData]) int {
		return int(nodes1[0].Payload.Size) - int(nodes2[0].Payload.Size)
	})
	var ttlDupsSizeSum int64 = 0
	var ttlVolumeSize int64 = 0
	for _, nodes := range duplicates {
		ttlVolumeSize += nodes[0].Payload.Size

		show := true
		for _, filter := range resultFilters {
			if !filter(nodes) {
				show = false
				break
			}
		}
		if !show {
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
		ttlDupsSizeSum += nodes[0].Payload.Size * int64(len(nodes)-1)

		fmt.Println()
	}

	fmt.Printf(
		"Estimated duplicates size ~%s (scanned ~%s)\n",
		humansize.SizeToString(ttlDupsSizeSum),
		humansize.SizeToString(ttlVolumeSize))
}
