startdb:
	docker run --name postgres_bank_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres_bank_db createdb --username=root --owner=root bank_db

dropdb:
	docker exec -it postgres_bank_db dropdb bank_db
	
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose down


.PHONY: startdb createdb dropdb migrateup migratedown