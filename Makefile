.PHONY: protos lint

# Generate Go bindings from proto file
protos:
	protoc --go_out=. --go_opt=paths=source_relative --go_opt=Mcfbd/internal/proto/cfbd.proto=github.com/clintrovert/cfbd-go/cfbd cfbd/internal/proto/cfbd.proto
	mv cfbd/internal/proto/cfbd.pb.go cfbd/generated.go

# Lint all non-test Go files using golangci-lint
lint:
	@echo "Installing/updating golangci-lint to latest version..."; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest
	@echo "Linting non-test Go files..."
	@LINT_OUTPUT=$$($$(go env GOPATH)/bin/golangci-lint run --skip-files='.*_test\.go$$' --skip-dirs='internal/test|internal/examples' ./cfbd/... 2>&1); \
	NON_TEST_ERRORS=$$(echo "$$LINT_OUTPUT" | grep "\.go:" | grep -v "_test.go:" || true); \
	if [ -n "$$NON_TEST_ERRORS" ]; then \
		echo "$$NON_TEST_ERRORS"; \
		echo "$$LINT_OUTPUT" | tail -1; \
		exit 1; \
	else \
		echo "All non-test files passed linting"; \
	fi
