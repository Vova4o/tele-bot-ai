#!/bin/sh

set -e 

echo $GOOSE_DBSTRING
goose -dir internal/storage/migrations up

/app/news-feed-bot