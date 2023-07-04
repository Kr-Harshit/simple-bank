.PHONY: 
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

createPostgres:
	docker run --name postgresDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

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
	go run main.go

test:
	go test -v --cover ./...

mock: 
	mockgen -package mock -destination ./db/mock/Store.go github.com/KHarshit1203/simple-bank/db/gen Store