gen: gen-sql gen-proto

install-generators:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get github.com/kyleconroy/sqlc/cmd/sqlc@v1.6.0

gen-proto: ## Generate protobuf and grpc definitions
	rm -f ./pkg/pb/*.go
	protoc --go_out=module=github.com/clstb/phi:. --go-grpc_out=module=github.com/clstb/phi:. proto/*.proto
gen-sql: ## Generate models and queries from sql definitions
	sqlc generate
