FROM golang:1.20-alpine AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN chmod +x entrypoint.sh
RUN go build -o /app/news-feed-bot ./cmd/

EXPOSE 8080

CMD ["/app/entrypoint.sh"]

# CMD ["/app/news-feed-bot"]