SRC_FILES := $(wildcard **/*.go) $(wildcard *.go)
EXCLUDE_SRC_FILES := %_test.go
SRC_FILES := $(filter-out $(EXCLUDE_SRC_FILES), $(SRC_FILES))
TEST_FILES := $(wildcard **/*_test.go)
GENERATED_FILES := vivian/grammar_pigeon.go

bin/vivid: $(GENERATED_FILES) $(SRC_FILES)
	go build -o bin/vivid -v

build: bin/vivid

vivian/grammar_pigeon.go: vivian/grammar.peg
	go get -u github.com/mna/pigeon
	pigeon vivian/grammar.peg | goimports > vivian/grammar_pigeon.go

.last_test: $(GENERATED_FILES) $(SRC_FILES) $(TEST_FILES)
	go test ./...
	touch .last_test

test: .last_test
