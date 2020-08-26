build:
	go build -o calendar ./cmd/calendar/main.go

test:
	go test -v -race ./... -count=1

lint:
	golangci-lint run ./...

run:
	go run ./cmd/calendar/main.go -config ./configs/config.yaml

migrateup:
	goose -dir migrations postgres "user=${CAL_DB_LOGIN} password=${CAL_DB_PASSWORD} dbname=${CAL_DB_NAME} sslmode=disable" up

migratedown:
	goose -dir migrations postgres "user=${CAL_DB_LOGIN} password=${CAL_DB_PASSWORD} dbname=${CAL_DB_NAME} sslmode=disable" down

generate:
	go generate proto/gen.go

integration-tests:
	docker-compose -f ./docker-compose-tests.yaml up --build --abort-on-container-exit --exit-code-from integration_tests && \
	docker-compose -f ./docker-compose-tests.yaml down

clean:
	rm -rf internal/pb/*
	rm -rf internal/api/swaggerui/statik/*
	rm -rf swaggerui/calendar/*.json

.PHONY: build proto
