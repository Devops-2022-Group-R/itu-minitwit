.PHONY: run
run:
	go run ./src

.PHONY: init
init:
	go run ./src "initDb"

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	gcc scripts/flag_tool.c -l sqlite3 -o scripts/flag_tool

.PHONY: clean
clean:
	rm scripts/flag_tool
