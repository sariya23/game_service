FROM golang:1.23.1-alpine3.19 as base
FROM base as dev
WORKDIR /app
COPY . .
RUN apk update
RUN go mod download
RUN go build -o api cmd/main.go
RUN apk add make