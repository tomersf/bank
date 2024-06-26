startdb:
	docker run --name postgres_bank_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres_bank_db createdb --username=root --owner=root bank_db

dropdb:
	docker exec -it postgres_bank_db dropdb bank_db

stopdb:
	docker stop postgres_bank_db
	
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/tomersf/bank/db/sqlc Store

.PHONY: startdb createdb dropdb stopdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock