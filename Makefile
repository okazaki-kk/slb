.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: test-pretty
test-pretty:
	set -o pipefail && go test -v ./... fmt -json | tparse -all

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: build
	go build -o slb
