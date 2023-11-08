up:
	echo "Starting Docker images..."
	docker compose up

down:
	echo "Stopping Docker images..."
	docker compose down

lint:
	golangci-lint run

test:
	go test ./... -v
