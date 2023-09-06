createDb:
	docker exec -it postgres createdb --username=postgress --owner=postgress go_finance

runPostgres:
	sudo docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

postgresPsql:
	sudo docker exec -it postgres psql -U postgres

startPostgres:
	sudo docker start postgres

migrateInit:
	migrate create -ext sql -dir db/migration -seq initial_tables

migrateup:
	migrate --path db/migration -database "postgresql://postgres:postgres@localhost:5432/go_finance?sslmode=disable" -verbose up

migratedrop:
	migrate --path db/migration -database "postgresql://postgres:postgres@localhost:5432/go_finance?sslmode=disable" -verbose drop

test:
	go test -v -cover ./...

server:
	go run main.go

sqlcGen:
	./sql generate

.PHONY: createDb runPostgres migrateup migratedrop test server sqlcGen