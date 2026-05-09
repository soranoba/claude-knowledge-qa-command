.PHONY: build install test clean

BINARY      = knowledge-qa
BINARY_PATH = ./bin/$(BINARY)
SKILL_DIR   ?= $(HOME)/.claude/commands

build:
	mkdir -p bin
	go build -o $(BINARY_PATH) .

install: build
	mkdir -p $(SKILL_DIR)/bin
	install -m 755 $(BINARY_PATH) $(SKILL_DIR)/bin/$(BINARY)
	install -m 644 .claude/commands/knowledge-qa.md $(SKILL_DIR)/knowledge-qa.md

test:
	go test ./...

clean:
	rm -f $(BINARY_PATH)
