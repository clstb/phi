clean:
	rm -rf ./pkg/pb/*
	rm -rf ./pkg/db/*
	rm -rf ./phi

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o phi

test:
	go test ./... -v
	
test-integration:
	go test ./... -tags=integration

install-generators:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get github.com/kyleconroy/sqlc/cmd/sqlc@v1.6.0

gen: gen-sql gen-proto

gen-proto: ## Generate protobuf and grpc definitions
	protoc --go_out=module=github.com/clstb/phi:. --go-grpc_out=module=github.com/clstb/phi:. proto/*.proto
gen-sql: ## Generate models and queries from sql definitions
	sqlc generate
