ifneq (,$(wildcard ./.env))
    include .env
    export
endif

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ./bin v2.0.2

lint:
	./bin/golangci-lint run 

test:
	go test -v -race ./... 

run:
	go run -race cmd/app/main.go

air:
	air

migrate_up:
	go run cmd/migrations/main.go up

migrate_down:
	go run cmd/migrations/main.go down