#!/bin/bash

if [ "$1" = "init" ]; then
    if [ -f "/tmp/minitwit.db" ]; then 
        echo "Database already exists."
        exit 1
    fi
    echo "Putting a database to /tmp/minitwit.db..."
    go run ./src "initDb"
elif [ "$1" = "start" ]; then
    echo "Starting minitwit..."
    nohup go run ./src > /tmp/out.log 2>&1 &
elif [ "$1" = "stop" ]; then
    echo "Stopping minitwit..."
    pkill -f src
elif [ "$1" = "inspectdb" ]; then
    ./scripts/flag_tool -i | less
elif [ "$1" = "flag" ]; then
    ./scripts/flag_tool "$@"
else
  echo "I do not know this command..."
fi
