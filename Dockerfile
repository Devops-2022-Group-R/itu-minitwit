# syntax=docker/dockerfile:1

FROM golang:1.17 AS builder
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o minitwit ./src

FROM alpine:3.15 AS runner
ENV ENVIRONMENT=PRODUCTION
COPY --from=builder /app/minitwit ./minitwit
CMD ["./minitwit"]
