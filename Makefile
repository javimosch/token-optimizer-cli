.PHONY: build release clean install

BINARY=token-optimizer-cli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GOFLAGS=-ldflags="-s -w" -trimpath

build:
	go build $(GOFLAGS) -o $(BINARY) .

release: build
	tar czf $(BINARY)-linux-amd64.tar.gz $(BINARY)
	sha256sum $(BINARY)-linux-amd64.tar.gz > $(BINARY)-linux-amd64.tar.gz.sha256

install: build
	cp $(BINARY) /usr/local/bin/

clean:
	rm -f $(BINARY) *.tar.gz *.sha256
