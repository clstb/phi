PROTO_DIR = proto
PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')

.DEFAULT_GOAL := proto
.PHONY: proto


proto: proto_go proto_py


proto_py:
	@python -m grpc_tools.protoc --python_out=proto --grpc_python_out=proto --proto_path=proto proto/shared.proto
	@python -m grpc_tools.protoc --python_out=proto --grpc_python_out=proto --proto_path=proto proto/ledger.proto
	@sed -i -e 's/import ledger_pb2 as ledger__pb2/from . import ledger_pb2 as ledger__pb2/g' proto/ledger_pb2_grpc.py
	@sed -i -e 's/import shared_pb2 as shared__pb2/from . import shared_pb2 as shared__pb2/g' proto/ledger_pb2_grpc.py
	@sed -i -e 's/import shared_pb2 as shared__pb2/from . import shared_pb2 as shared__pb2/g' proto/ledger_pb2.py
	@mv proto/*.py fava/src/fava/file



proto_go:
	@protoc -I${PROTO_DIR} \
	--go_opt=module=${PACKAGE} \
	--go_out=. \
	--go-grpc_opt=module=${PACKAGE} \
	--go-grpc_out=. \
	${PROTO_DIR}/*.proto
