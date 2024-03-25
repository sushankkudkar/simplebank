DB_URL=postgresql://root:MBkgwzTaSzjqWpnIjIGq@simple-bank.cvnxivhqxepy.ap-south-1.rds.amazonaws.com:5432/simple_bank

network:
	docker network create bank-network

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

lint:
	golint ./...

server:
	go run main.go

.PHONY: network postgres createdb dropdb migrateup1 migratedown1 migrateup migratedown sqlc lint server
