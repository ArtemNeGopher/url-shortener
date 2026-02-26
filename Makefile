.PHONY: proto up down logs test clean migrate-up migrate-down new-migrate psql

proto:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

up:
	docker-compose up --build -d

down:
	docker-compose down -v

logs:
	docker-compose logs -f

test:
	go test -v -race -coverprofile=coverage.out ./...

clean:
	docker-compose down -v
	rm -rf coverage.out
	find . -name "*.pb.go" -delete

migrate-up:
	migrate -path ./services/analytics-service/migrations -database "postgres://urlshortener:password@localhost:5432/urlshortener?sslmode=disable" up

migrate-down:
	docker-compose exec url-service migrate -path /migrations -database "postgres://urlshortener:password@postgres:5432/urlshortener?sslmode=disable" down

new-migrate:
	@read -p "Enter service name: " service; \
	read -p "Enter migration name: " name; \
	migrate create -seq -ext sql -dir ./services/$$service/migrations $$name

bench:
	go test -bench=. -benchmem ./pkg/...

psql:
	docker-compose exec postgres psql -U urlshortener urlshortener
