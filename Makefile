.PHONY: all
all: build test

.PHONY: build
build: lc-go

.PHONY: install
install:
	go install .
	cd "$$(dirname "$$(command -v lc-go)")" && ln -s lc-go lc

.PHONY: clean
clean:
	rm -f ls-go

lc-go: main.go
	go build

.PHONY: test
test:
	go test
