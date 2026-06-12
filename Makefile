.PHONY: build test release-dry-run clean

build:
	go build -o bin/scanrail ./cmd/scanrail

test:
	go test ./...
	go vet ./...
	node packages/npm/cli/test-wrapper.mjs

release-dry-run:
	go test ./...
	go vet ./...
	node packages/npm/cli/test-wrapper.mjs
	node scripts/build-release.mjs
	npm pack --workspaces --dry-run

clean:
	rm -rf bin dist .scanrail
