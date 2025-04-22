ifneq (,$(wildcard ./.env))
    include .env
    export
endif

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ./bin v2.0.2

lint:
	./bin/golangci-lint run

test:
	go test ./... -v

run:
	go run cmd/app/main.go

air:
	air