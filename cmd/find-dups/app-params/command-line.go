package appparams

import (
	"flag"
	"strings"
)

type CmdLineParams struct {
	dirs        []string
	useHash     bool
	failOnError bool
	onlyDirs    bool
	filterDirs  []string
}

func NewCmdLineParams() *CmdLineParams {
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

	return &CmdLineParams{
		dirs:        dirs,
		filterDirs:  pathFilters,
		useHash:     *useHashParam,
		failOnError: *failOnErrorParam,
		onlyDirs:    *onlyDirsParam,
	}
}

func (p *CmdLineParams) Dirs() []string {
	return p.dirs
}

func (p *CmdLineParams) UseHash() bool {
	return p.useHash
}

func (p *CmdLineParams) FailOnError() bool {
	return p.failOnError
}

func (p *CmdLineParams) DirsOnly() bool {
	return p.failOnError
}

func (p *CmdLineParams) DirsFilters() []string {
	return p.filterDirs
}
