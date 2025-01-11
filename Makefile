.PHONY: all
all: lc-go test

.PHONY: clean
clean:
	rm -f ls-go

lc-go: main.go
	go build

.PHONY: test
test:
	go test
