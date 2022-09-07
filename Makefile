
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
	for dir in $$(go list -f '{{.Dir}}' ./...); do \
		(cd $$dir && gofmt -s -w $$(go list -f '{{.GoFiles}} {{.TestGoFiles}}' | tr -d '[]')); \
	done

.PHONY: test
test:
	[ -n "$(run)" ] && r="-run $(run)"; \
	go test -coverprofile=coverage.html -v ./... $$r

.PHONY: coverage
coverage:
	go tool cover -html=coverage.html
