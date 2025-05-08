setup:
	@echo Running setup...
	@docker-compose up db -d

install:
	@echo Installing dependencies...
	@go mod tidy

test:
	@echo Running tests...
	@go test -v ./...
	@echo Tests passed!

cover:
	@echo Running test coverage...
	@go test -v ./... -coverprofile=coverage/cover.out
	@go tool cover -html coverage/cover.out -o coverage/cover.html
	@go tool cover -func coverage/cover.out
	@echo Test coverage successfully!

run:
	@go run cmd/main.go

docker-build:
	@echo Running Docker build:
	@docker-compose up app -d