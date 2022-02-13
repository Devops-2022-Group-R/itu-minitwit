init:
	go run ./src "initDb"

build:
	gcc scripts/flag_tool.c -l sqlite3 -o scripts/flag_tool

clean:
	rm scripts/flag_tool
