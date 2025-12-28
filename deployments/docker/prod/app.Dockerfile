FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api cmd/main.go

FROM alpine:3.20 AS app
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/api /app/api
COPY config ./config
EXPOSE 47900 46900
