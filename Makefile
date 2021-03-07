build: *.peg *.go
	go build -o bin/vivid -v

*.peg:
	go get -u github.com/mna/pigeon
	pigeon vivian/grammar.peg | goimports > vivian/grammar_pigeon.go

.PHONY: test


test:
	go test ./...
