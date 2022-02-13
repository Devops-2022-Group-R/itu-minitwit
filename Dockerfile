# syntax=docker/dockerfile:1

FROM golang:1.17

RUN apt-get update \
    && apt-get install -y sqlite3 libsqlite3-dev

WORKDIR /app
COPY . /app

RUN go build -o minitwit ./src

CMD ["/app/minitwit"]
