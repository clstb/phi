package pb

//go:generate protoc --go_out=module=github.com/clstb/phi/go/pkg/services/auth/pb:. --go-grpc_out=module=github.com/clstb/phi/go/pkg/services/auth/pb:. --grpc-gateway_out=module=github.com/clstb/phi/go/pkg/services/auth/pb:.  auth.proto
