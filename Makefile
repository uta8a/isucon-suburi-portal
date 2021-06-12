ci:
	golangci-lint run ./...
dev:
	docker-compose -f ./local-dev/dev/docker-compose.yml up --build
down:
	docker-compose -f ./local-dev/dev/docker-compose.yml down
