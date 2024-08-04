# Build stage
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY ./internal/ ./internal/
COPY ./cmd/ ./cmd/

RUN go build -o /app/bin/app ./cmd/app/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/app /app/bin/app

CMD ["/app/bin/app"]
