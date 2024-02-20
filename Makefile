# Define variables for Docker commands to simplify usage
DOCKER_COMPOSE = sudo docker-compose
MIGRATE = sudo docker run --network host migrate/migrate

# Project variables
BINARY_NAME = ventrata_octo
DB_MIGRATIONS_PATH = ./migrations
DB_CONNECTION_STRING = postgres://postgres:postgres@localhost:5432/ventrata_octo?sslmode=disable

.PHONY: build run test clean docker-build docker-up docker-down migrate-up migrate-down

# Build the Go binary.
build:
	@echo "Building..."
	go build -o ${BINARY_NAME} main.go

# Run the Go application.
run: build
	@echo "Running..."
	./${BINARY_NAME}

# Run tests.
test:
	@echo "Testing..."
	go test ./...

# Clean up binaries.
clean:
	@echo "Cleaning..."
	go clean
	rm -f ${BINARY_NAME}

# Build the Docker container for the application.
docker-build:
	@echo "Building Docker image..."
	${DOCKER_COMPOSE} build app

# Start all services defined in the Docker Compose file.
docker-up:
	@echo "Starting Docker containers..."
	${DOCKER_COMPOSE} up -d

# Stop all services defined in the Docker Compose file.
docker-down:
	@echo "Stopping Docker containers..."
	${DOCKER_COMPOSE} down

# Apply database migrations.
migrate-up:
	@echo "Applying database migrations..."
	${MIGRATE} -path=${DB_MIGRATIONS_PATH} -database "${DB_CONNECTION_STRING}" up

# Revert database migrations.
migrate-down:
	@echo "Reverting database migrations..."
	${MIGRATE} -path=${DB_MIGRATIONS_PATH} -database "${DB_CONNECTION_STRING}" down