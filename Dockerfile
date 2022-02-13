# syntax=docker/dockerfile:1

FROM golang:1.17

RUN apt-get update \
    && apt-get install -y sqlite3 libsqlite3-dev

RUN git config --global core.autocrlf true

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o main ./src

CMD ["/app/main"]
