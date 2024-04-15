FROM golang:1.20 AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

# COPY entrypoint.sh .

COPY . .

RUN go mod download

RUN chmod +x entrypoint.sh
RUN go build -o /app/news-feed-bot ./cmd/

ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=${NFB_DATABASE_DSN:-postgres://postgres:postgres@db:5432/news_feed_bot?sslmode=disable}
ENV NFB_TELEGRAM_BOT_TOKEN=6924410271:AAELRqS3nSMMJkAH3BJRVkjpHgGAjxDflpw
ENV NFB_TELEGRAM_CHANNEL_ID=-1002052354012

EXPOSE 8080

# RUN ls
CMD /app/entrypoint.sh

# CMD ["/app/news-feed-bot"]