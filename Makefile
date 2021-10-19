deps:
	go mod tidy
	bazel run  //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
	bazel run  //:gazelle -- -build_tags=integration
