PG_DSN="postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
LOCAL_BIN := $(CURDIR)/bin

bin-deps:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

lint:
	$(LOCAL_BIN)/golangci-lint run --new-from-rev=origin/main ./...

test:
	$(info #Running integration tests...)
	PG_DSN=${PG_DSN} \
	go test -v -covermode=count -coverprofile=coverage.out ./... && go tool cover -html coverage.out -o coverage.html