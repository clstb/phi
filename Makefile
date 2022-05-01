test:
	go test go/pkg/util/test.go
	go test go/pkg/auth/server/server_test.go

clean_stubs: clean_stubs_auth clean_stubs_tinkgw

clean_stubs_auth:
	rm -f go/pkg/auth/db/db.go
	rm -f go/pkg/auth/db/models.go
	rm -f go/pkg/auth/db/querier.go
	rm -f go/pkg/auth/db/users.sql.go
	rm -f go/pkg/auth/pb/auth.pb.go
	rm -f go/pkg/auth/pb/auth_grpc.pb.go

clean_stubs_tinkgw:
	rm -f go/pkg/tinkgw/db/db.go
	rm -f go/pkg/tinkgw/db/models.go
	rm -f go/pkg/tinkgw/db/querier.go
	rm -f go/pkg/tinkgw/db/users.sql.go
	rm -f go/pkg/tinkgw/pb/tinkgw.pb.go
	rm -f go/pkg/tinkgw/pb/tinkgw_grpc.pb.go

stubs: stubs_auth stubs_tinkgw

stubs_auth:
	go generate github.com/clstb/phi/go/pkg/auth/db
	go generate github.com/clstb/phi/go/pkg/auth/pb

stubs_tinkgw:
	go generate github.com/clstb/phi/go/pkg/tinkgw/db
	go generate github.com/clstb/phi/go/pkg/tinkgw/pb




