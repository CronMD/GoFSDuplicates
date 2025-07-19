package main

import (
	humansize "df/cmd/find-dups/human-size"
	"df/internal/nodes"
	"df/internal/payloads/fsdata"
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"
)

func main() {
	dirsParam := flag.String("dirs", "", "comma separated dir paths")
	useHashParam := flag.Bool("hash", false, "check files hashes")
	failOnError := flag.Bool("fail", false, "fail on error")
	flag.Parse()

	dirs := make([]string, 0)
	for _, dir := range strings.Split(strings.Trim(*dirsParam, " "), ",") {
		dir = strings.Trim(dir, " ")
		if dir != "" {
			dirs = append(dirs, dir)
		}
	}

	src := fsdata.NewMultipleDirsFsDataSource(dirs...)

	ixs := []nodes.Indexer[fsdata.FsData]{
		fsdata.NewNameSizeFsIndexer(),
	}
	if *useHashParam {
		ixs = append(ixs, fsdata.NewHashFsIndexer(0.02, !*failOnError))
	}

	finder := nodes.NewDupFinder(ixs)

	duplicates, err := finder.FindFromSources(src)
	if err != nil {
		log.Fatal(err)
	}

	slices.SortFunc(duplicates, func(nodes1, nodes2 []*nodes.Node[fsdata.FsData]) int {
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
