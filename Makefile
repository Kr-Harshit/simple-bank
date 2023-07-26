.PHONY: 
	create-network
	run-db
	delete-db
	stop-db 
	start-db 
	remove-db
	create-migration 
	migate-up 
	migrate-one-up
	migrate-one-down
	migrate-down
	sqlc
	test
	mock
	run-server
	build-server 



create-network:
	docker network create simple-bank

####### DATABASE ########

run-db:
	docker run \
	--name simple-bank-db \
	-p 5432:5432 \
	-e POSTGRES_USER=root \
	-e POSTGRES_PASSWORD=secret \
	-e POSTGRES_DB=simple-bank \
	--network simple-bank \
	-d postgres:alpine

delete-db:
	docker rm -f simple-bank-db

stop-db:
	docker stop simple-bank-DB

start-db:
	docker start simple-bank-DB

remove-db: 
	docker exec -it postgresDB dropdb simple_bank

login-db:
	docker exec -it postgresDB psql -U root -W


###### SERVER #######

run-server:
	docker run \
	--name simple-bank-server \
	-p 8080:8080 \
	--network simple-bank \
	-e ENV_DATABASE_SOURCE="postgres://root:secret@simple-bank-DB:5432/simple-bank?sslmode=disable" \
	-e ENV_TOKEN_KEY="2bddd92e3cbe124f65091c39bdad5bf5" \
	-d simplebank:latest

build-server:
	docker build -f docker/Dockerfile -t simplebank:latest .

test-server:
	go test -v --cover ./...

###### TOOLS #########

create-migration:	
	migrate create -ext sql -dir service/db/migration -seq -digits 3 {## give schema name ##}

migrate-up:
	migrate -path service/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose up

migrate-one-up:
	migrate -path service/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose up 1

migrate-down:
	migrate -path service/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose down

migrate-one-down:
	migrate -path service/db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose down 1

sqlc:
	sqlc generate

mock: 
	mockery --dir service/db/gen --output service/db/mocks  --name Store  