gen-proto: ## Generate protobuf and grpc definitions
	rm -f ./pkg/pb/*.go
	protoc --go_out=module=github.com/clstb/phi:. --go-grpc_out=module=github.com/clstb/phi:. proto/*.proto
