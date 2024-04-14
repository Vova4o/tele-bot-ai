#!/bin/bash

set -e 

echo $GOOSE_DBSTRING
goose -dir internal/storage/migrations up

/main