package main

import (
	"df/internal/sources"
	"flag"
	"fmt"
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
	for leaf := range src.Leafs() {
		fmt.Println(leaf)
	}
}
