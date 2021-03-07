bin/vivid: **/*.go
	go build -o bin/vivid -v

build: bin/vivid

vivian/grammar_pigeon.go: vivian/grammar.peg
	go get -u github.com/mna/pigeon
	pigeon vivian/grammar.peg | goimports > vivian/grammar_pigeon.go

.PHONY: test

test: build
	go test ./...
