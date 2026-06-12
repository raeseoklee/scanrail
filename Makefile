.PHONY: build test release-dry-run npm-publish-dry-run npm-publish experiment-scanner-spike tape-scanner-spike tape-headers-demo tape-mcp-verification clean

build:
	go build -o bin/scanrail ./cmd/scanrail

test:
	go test ./...
	go vet ./...
	node packages/npm/cli/test-wrapper.mjs
	node packages/npm/scanrail/test-wrapper.mjs

release-dry-run:
	go test ./...
	go vet ./...
	node packages/npm/cli/test-wrapper.mjs
	node packages/npm/scanrail/test-wrapper.mjs
	node scripts/build-release.mjs
	npm pack --workspaces --dry-run

npm-publish-dry-run:
	npm run publish:dry-run

npm-publish:
	npm run publish:npm

experiment-scanner-spike:
	node experiments/scanner-adapter-spike/run.mjs

tape-scanner-spike:
	vhs experiments/scanner-adapter-spike/tapes/scanner-adapter-spike.tape

tape-headers-demo:
	vhs examples/headers-demo/tapes/headers-demo.tape

tape-mcp-verification:
	vhs examples/mcp-verification/tapes/mcp-verification.tape

clean:
	rm -rf bin dist .scanrail
