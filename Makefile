build: *.go
	go build -o bin/vivid -v

.PHONY: test

test:
	go test ./...
