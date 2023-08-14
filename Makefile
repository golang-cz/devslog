.PHONY: fmt check-fmt lint vet test

test:
	@go test

test-coverage:
	@go test -cover -coverprofile=coverage.out ./... -json

test-coverage-inspect: test-coverage
	@go tool cover -html=coverage.out

test-tparse:
	@go test -cover -coverprofile=coverage.out ./... -json | tparse -all
