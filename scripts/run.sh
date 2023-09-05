#!/usr/bin/env bash

export DATABASE_USER=mongouser
export DATABASE_PASSWORD=mongopass
export DATABASE_URI=localhost:8081
export ENVIRONMENT=test
export PORT=8084

docker run \
    -d \
    --name mongo \
    -p 8081:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=mongouser \
    -e MONGO_INITDB_ROOT_PASSWORD=mongopass \
    mongo & \

go run ./cmd/main.go