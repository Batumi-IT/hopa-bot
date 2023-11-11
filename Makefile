up:
	echo "Starting Docker images..."
	docker compose up

up-build:
	echo "Starting Docker images..."
	docker compose up --build --remove-orphans

down:
	echo "Stopping Docker images..."
	docker compose down

lint:
	golangci-lint run

test:
	go test ./... -v
