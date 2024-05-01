FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
COPY ./internal/ ./internal/
COPY ./cmd/ ./cmd/

RUN go mod tidy
RUN go mod download

RUN go build -o ./bin/app/main ./cmd/app/main.go
CMD ["./bin/app/main"]
