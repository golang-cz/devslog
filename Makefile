.PHONY: fmt check-fmt lint vet test

test:
	@go test -race

test-coverage:
	@go test -race -cover -coverprofile=coverage.txt ./... -json

test-coverage-inspect: test-coverage
	@go tool cover -html=coverage.txt

test-tparse:
	@go test -race -cover -coverprofile=coverage.txt ./... -json | tparse -all
