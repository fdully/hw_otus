build:
	go build -o calendar ./cmd/calendar/main.go

test:
	go test -race ./... -count=1

lint:
	golangci-lint run ./...

run:
	go run ./cmd/calendar/main.go -config ./configs/config.yaml

migrate_up:
	goose -dir migrations postgres "user=${CAL_DB_LOGIN} password=${CAL_DB_PASSWORD} dbname=${CAL_DB_NAME} sslmode=disable" up

migrate_down:
	goose -dir migrations postgres "user=${CAL_DB_LOGIN} password=${CAL_DB_PASSWORD} dbname=${CAL_DB_NAME} sslmode=disable" down

.PHONY: build
