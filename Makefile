#********************************************************************************#
#				Targets for developing, testing, buiulding project
#********************************************************************************#


APPLICATION_NAME = simple-bank

#********************************************************************************#
#				Build, Test and Lint
#********************************************************************************#

build-binary:
	CGO_ENABLED=0 GOPROXY=direct GOARCH=amd64 GOOS=linux go build -o $(APPLICATION_NAME) main.go

test:
	go test -v --cover ./...

dep:
	go mod download

vet:
	go vet

lint:
	docker run -t --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.53.3 golangci-lint run -v


#********************************************************************************#
#				local debugging
#********************************************************************************#

DEBUG_DATABASE_SOURCE = postgres://root:secret@localhost:5432/$(APPLICATION_NAME)?sslmode=disable
MIGRATION_VERSION =
db_status := $(shell docker ps -f name=simple-bank-debug-db -q )

export ENV_DATABASE_SOURCE=$(DEBUG_DATABASE_SOURCE)
export ENV_TOKEN_KEY=$(shell openssl rand -hex 16)
export ENV_DATABASE_MIGRATE_SOURCE=file://service/db/migration

debug-server: debug-clean debug-database
	CGO_ENABLED=0 GOPROXY=direct go run main.go

debug-database:
	docker run --name $(APPLICATION_NAME)-debug-db -p 5432:5432 -e POSTGRES_DB=$(APPLICATION_NAME) -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine 

debug-database-login:
	docker exec -it $(APPLICATION_NAME)-debug-db psql -U root -w $(APPLICATION_NAME)

debug-clean:
	docker rm -f $(APPLICATION_NAME)-debug-db

debug-migrate-up:
	migrate -path service/db/migration -database $(DEBUG_DATABASE_SOURCE) --verbose up $(MIGRATION_VERSION)

debug-migrate-down:
	migrate -path service/db/migration -database $(DEBUG_DATABASE_SOURCE) --verbose down $(MIGRATION_VERSION)


#********************************************************************************#
#				local docker-compose testing
#********************************************************************************#

compose-test:
	docker compose up

compose-clean:
	docker compose down
	docker rmi simplebank-api


#********************************************************************************#
#				development tools
#********************************************************************************#

SCHEMA_NAME = "init"

create-migration:	
	migrate create -ext sql -dir service/db/migration -seq -digits 3 $(SCHEMA_NAME)

sqlc:
	sqlc generate

mock: 
	mockery --dir service/db/gen --output service/db/mocks  --name Store  