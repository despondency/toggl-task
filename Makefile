run:
	./scripts/run.sh

run-mongo:
	docker run \
		-d \
		--name mongo \
		-p 8081:27017 \
		-e MONGO_INITDB_ROOT_USERNAME=mongouser \
		-e MONGO_INITDB_ROOT_PASSWORD=mongopass \
		mongo & \


test-unit:
	go test -run "Unit" -v ./internal/...

test-integration:
	go test -run "Integration" -v ./tests/...