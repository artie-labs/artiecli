.PHONY: build
build:
	goreleaser build --clean

.PHONY: release
release:
	goreleaser release --clean
