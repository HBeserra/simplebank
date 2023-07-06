services:
	docker-compose -f ./services-compose.yaml up -d

createdb:
	docker exec -it simplebank-postgresdb createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simplebank-postgresdb dropdb simple_bank

install_migrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
	sudo mv migrate /usr/bin/migrate

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

psql:
	docker-compose -f ./services-compose.yaml exec -it postgres psql -U root -d simple_bank

test:
	go test -v -cover -short ./...

.PHONY: services createdb dropdb migrateup migratedown sqlc test psql install_migrate