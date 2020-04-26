NAME=gdp
BIN=bin
VERSION=$(shell git describe --tags --abbrev=0)
OS=darwin
ARCH=amd64

all: setup mod test build
setup:
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports
mod:
	go mod download
test:
	go test -v -cover ./...
build:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BIN)/$(NAME) -v
clean:
	go clean
	rm -rf $(BIN)/
compress:
	tar cvzf $(BIN)/$(NAME)_$(VERSION)_$(OS)_$(ARCH).tar.gz $(BIN)/$(NAME)
