NAME=gdp
BIN=bin
VERSION=$(shell git describe --tags --abbrev=0)
OS=darwin
ARCH=amd64

all: setup deps test build
setup:
	go get github.com/golang/dep/cmd/dep
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
deps:
	dep ensure -v
update:
	dep ensure -v -update
test:
	go test -v -cover ./...
build:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BIN)/$(NAME) -v
clean:
	go clean
	rm -rf vendor/
	rm -rf $(BIN)/
compress:
	tar cvzf $(BIN)/$(NAME)_$(VERSION)_$(OS)_$(ARCH).tar.gz $(BIN)/$(NAME)
