PROTO_DIR = proto
PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')

.DEFAULT_GOAL := proto
.PHONY: proto


proto: proto_go proto_py


proto_py:
	@python -m grpc_tools.protoc           				\
	--python_out fava/src/fava/file/proto       		\
	--grpc_python_out=fava/src/fava/file/proto 		 \
	--proto_path src/fava/file/proto      				 \
	proto/ledger.proto;								\
	sed -i -e 's/import ledger_pb2 as ledger__pb2/from . import ledger_pb2 as ledger__pb2/g' src/fava/file/proto/ledger_pb2_grpc.py



proto_go:
	@protoc -I${PROTO_DIR} \
	--go_opt=module=${PACKAGE} \
	--go_out=. \
	--go-grpc_opt=module=${PACKAGE} \
	--go-grpc_out=. \
	${PROTO_DIR}/*.proto

clean:
		rm -f ${PROTO_DIR}/*.go; \
		rm -r -f ui/node_modules; \
		rm -f  ui/pnpm-lock.yaml



fava_frontend:
	@cd fava/frontend; \
	npm i;        \
	npm run build


clean_frontend:
	@rm -f fava/src/fava/static/app.js;   \
	rm -f fava/src/fava/static/app.css;   \
	rm -f fava/src/fava/static/*.woff;    \
	rm -r -f fava/frontend/node_modules;  \
	rm -f fava/frontend/package-lock.json


clean_proto:
	@rm -f src/fava/file/proto/*.py; \
	rm -f src/fava/file/proto/*.py-e; \
