.PHONY: 
	createNetwork
	createPostgres
	deletePostgres
	stopPostgres 
	restartPostgres 
	createDb 
	dropDb 
	createMigration 
	migateUp 
	migrateUp1
	migrateDown
	migrateDown1
	sqlc
	test
	mock 

createNetwork:
	docker network create simple-bank

createPostgres:
	docker run --name postgresDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret --network simple-bank -d postgres:alpine

deletePostgres:
	docker rm -f postgresDB

stopPostgres:
	docker stop postgresDB

restartPostgres:
	docker start postgresDB

loginPostgres:
	docker exec -it postgresDB psql -U root -W

createDb:
	docker exec -it postgresDB createdb --username=root --owner=root simple_bank 

dropDb: 
	docker exec -it postgresDB dropdb simple_bank

createMigration:	
	migrate create -ext sql -dir service/db/migration -seq -digits 3 {## give schema name ##}

migrateUp:
	migrate -path service/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose up

migrateUp1:
	migrate -path service/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose up 1

migrateDown:
	migrate -path service/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose down

migrateDown1:
	migrate -path service/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose down 1
	

sqlc:
	sqlc generate

server:
	docker run --name apiserver -p 8080:8080 --network simple-bank -e DB_SOURCE="postgres://root:secret@postgresDB:5432/simple_bank?sslmode=disable" -d simplebank:latest

test:
	go test -v --cover ./...

mock: 
	mockery --dir service/db/gen --output service/db/mocks  --name Store  