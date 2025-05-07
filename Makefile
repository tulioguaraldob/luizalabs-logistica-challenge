setup:
	@echo Running setup...
	@docker-compose up -d

install:
	@echo Installing dependencies...
	@go mod tidy

run:
	@go run cmd/main.go