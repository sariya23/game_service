FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY . .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api cmd/main.go


FROM alpine:3.20 AS app

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api /app/api

CMD ["/app/api"]


FROM golang:1.25-alpine AS migrations

WORKDIR /app

RUN apk add --no-cache postgresql-client git make

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.21.1

COPY migrations ./migrations

ENTRYPOINT ["goose"]
