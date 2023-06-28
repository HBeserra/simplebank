services:
	docker-compose -f ./services-compose.yaml up -d

createdb:
	docker exec -it simplebank-postgresdb createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simplebank-postgresdb dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

.PHONY: services createdb dropdb migrateup migratedown sqlc test