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
	failOnErrorParam := flag.Bool("fail", false, "fail on error")
	onlyDirsParam := flag.Bool("only-dirs", false, "show only duplicated directories")
	filterDirsParam := flag.String("filter", "", "comma separted paths")
	flag.Parse()

	dirs := make([]string, 0)
	for _, dir := range strings.Split(strings.Trim(*dirsParam, " "), ",") {
		dir = strings.Trim(dir, " ")
		if dir != "" {
			dirs = append(dirs, dir)
		}
	}

	pathFilters := make([]string, 0)
	for _, path := range strings.Split(strings.Trim(*filterDirsParam, " "), ",") {
		path = strings.Trim(path, " ")
		if path != "" {
			pathFilters = append(pathFilters, strings.ToLower(path))
		}
	}

	src := fsdata.NewMultipleDirsFsDataSource(
		fsdata.MulDirsDataSrcWithDirs(dirs...),
		fsdata.MulDirsDataSrcWithDFailOnError(*failOnErrorParam),
	)

	ixs := []nodes.Indexer[fsdata.FsData]{
		fsdata.NewNameSizeFsIndexer(),
	}
	if *useHashParam {
		ixs = append(ixs, fsdata.NewHashFsIndexer(
			fsdata.WithSuppressHashIndexErrors(!*failOnErrorParam)))
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

		if len(pathFilters) > 0 {
			show := false

			for _, node := range nodes {
				for _, filter := range pathFilters {
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
			if *onlyDirsParam {
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
