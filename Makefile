# SPDX-License-Identifier: Apache-2.0

.PHONY: all
all: mod compile lint format test

.PHONY: mod
mod:
	go mod tidy

.PHONY: compile
compile:
	go build ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: format
format:
	for pkg in `go list -f '{{.Dir}}' ./...`; do gofmt -s -w $${pkg}; done

.PHONY: test
test:
	testOpt="-count=$${count:-1}"; \
	[ -z "$(run)" ] || testOpt="$$testOpt -run $(run)"; \
	go test -coverprofile=coverage.html -v $$testOpt ./...

.PHONY: coverage
coverage:
	go tool cover -html=coverage.html
