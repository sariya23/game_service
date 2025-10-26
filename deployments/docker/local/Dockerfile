FROM golang:1.25-alpine as base
FROM base as dev
WORKDIR /app
COPY . .
RUN apk update
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN apk update && apk add --no-cache postgresql-client
RUN go mod download
RUN go build -o api cmd/main.go
RUN apk add make