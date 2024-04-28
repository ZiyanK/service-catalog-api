FILE=./database/data_model.sql
DATAMODEL=`cat $(FILE)`

db-up:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres

db-down:
	@if docker ps -a --format '{{.Names}}' | grep -q '^postgres$$'; then \
        docker stop postgres; \
        docker rm postgres; \
        echo "Container stopped."; \
    else \
        echo "Container not found, skipping stop command."; \
    fi

db-connect:
	docker exec -it postgres psql -U root -p 5432 -d postgres

db-migrate:
	docker exec -it postgres psql -U root -p 5432 -d postgres -exec "CREATE DATABASE service_catalog;" -exec "\c service_catalog;" -exec "$(DATAMODEL)"

db-drop:
	docker exec -it postgres psql -U root -p 5432 -d postgres -exec "DROP DATABASE service-catalog;"

run:
	go run app/cmd/main.go

build:
	go build -o ./service-catalog-api app/cmd/main.go

test:
	go test -v ./...

format:
	go fmt ./...

spec-up:
	docker run -d --rm --name swagger -p 8080:8080 -v $$PWD/spec:/spec -e SWAGGER_JSON=/spec/openapi.yaml swaggerapi/swagger-ui

spec-down:
	docker stop swagger