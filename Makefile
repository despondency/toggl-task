run:
	./scripts/run.sh

test-unit:
	go test -run "Unit" -v ./internal/...

test-integration:
	go test -run "Integration" -v ./tests/...