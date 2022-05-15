PROTO_DIR = go/proto
PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')

.DEFAULT_GOAL := proto
.PHONY: proto


proto:
	@protoc -I${PROTO_DIR} --go_opt=module=${PACKAGE} --go_out=go/ --go-grpc_opt=module=${PACKAGE} --go-grpc_out=go/ ${PROTO_DIR}/*.proto


clean:
		rm ${PROTO_DIR}/*.go
