run:
	docker run \
    -d \
    --name mongo \
    -p 8081:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=mongouser \
    -e MONGO_INITDB_ROOT_PASSWORD=mongopass \
    mongo & \
    go run ./cmd/main.go