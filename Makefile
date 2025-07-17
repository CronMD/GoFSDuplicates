test:
	go test -count 1 ./...

build-win:
	GOOS=windows GOARCH=amd64 go build -C cmd/find-dups/ -o ../../build/dedup.exe
