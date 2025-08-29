OUT_FILE=dedup
BUILD_DIR=build
BUILD_OUT_FILE=$(BUILD_DIR)/$(OUT_FILE)
WIN_BUILD_OUT_FILE=$(BUILD_OUT_FILE).exe

DEPLOY_DIR = /Volumes/data/install/

test:
	go test -count 1 ./...

build-win:
	GOOS=windows GOARCH=amd64 go build -C cmd/find-dups/ -o ../../$(WIN_BUILD_OUT_FILE)

deploy-win: build-win
	cp $(WIN_BUILD_OUT_FILE) $(DEPLOY_DIR)
