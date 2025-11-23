.PHONY: up down restart logs

up:
	docker-compose up

down:
	docker-compose down

clean:
	docker-compose down -v

logs:
	docker-compose logs -f

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run ./...

fmt:
	golangci-lint run --fix ./...

test:
	go test ./tests -v
