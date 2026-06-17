.PHONY: test-unit test-e2e test-all

test-unit:
	go test ./internal/domain/... -v -cover

test-e2e:
	go test ./test/e2e/... -v

test-all: test-unit test-e2e
