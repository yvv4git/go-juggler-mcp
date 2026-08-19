.PHONY: build test vet fmt lint mock mocks

MOCKGEN = go run go.uber.org/mock/mockgen

build:
	go build -o go-juggler-mcp .

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l -w .

lint:
	go vet ./...
	go test ./...

mock:
	$(MOCKGEN) -package=core -destination=internal/core/server_mock.go github.com/yvv4git/go-juggler-mcp/internal/ports BrowserClient,Filesystem

mocks: mock