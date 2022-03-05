VERSION ?= $(shell git rev-parse HEAD | cut -c 1-8)
DATABASE_URL ?= None

deps:
	go mod tidy
	bazel run  //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
	bazel run //:gazelle

build:
	bazel build --define version=$(VERSION) //...

test:
ifndef DATABASE_URL
$(error DATABASE_URL is not set)
endif
	bazel test --define version=$(VERSION) --test_output=errors --test_env DATABASE_URL=$(DATABASE_URL) //...

push:
	bazel run --define version=$(VERSION) --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //:push

release:
	gox -osarch="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64" ./go/cmd/phi/