PROTO_DIR = proto
PACKAGE = $(shell head -1 2048/go.mod | awk '{print $$2}')

.DEFAULT_GOAL := proto
.PHONY: proto

proto: proto_go proto_py

all: clean proto

proto_go:
	@mkdir proto/go; \
	protoc -I${PROTO_DIR} \
	--go_opt=module=${PACKAGE} \
	--go_out=proto/go \
	--go-grpc_opt=module=${PACKAGE} \
	--go-grpc_out=proto/go \
	${PROTO_DIR}/*.proto

proto_py:
	@mkdir proto/python; \
	python -m grpc_tools.protoc  --python_out proto/python --grpc_python_out=proto/python --proto_path proto proto/*.proto
	sed -i -e 's/import game_pb2 as game__pb2/from . import game_pb2 as game__pb2/g' proto/python/game_pb2_grpc.py
	sed -i -e 's/import env_pb2 as env__pb2/from . import env_pb2 as env__pb2/g' proto/python/env_pb2_grpc.py
	touch proto/python/__init__.py


clean:
		rm -f -r proto/go;
		rm -f -r proto/python;