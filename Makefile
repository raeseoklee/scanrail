.PHONY: build test release-dry-run experiment-scanner-spike tape-scanner-spike clean

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

experiment-scanner-spike:
	node experiments/scanner-adapter-spike/run.mjs

tape-scanner-spike:
	vhs experiments/scanner-adapter-spike/tapes/scanner-adapter-spike.tape

clean:
	rm -rf bin dist .scanrail
