.PHONY: init doc test

init:
	go mod tidy
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@if [ ! -e .git/hooks/pre-commit ]; then \
		chmod 755 .githooks/pre-commit; \
		ln -s $(PWD)/.githooks/pre-commit .git/hooks/pre-commit; \
	fi

doc:
	@gomarkdoc \
		--output usage.md .

lint:
	golangci-lint run ./...

test:
	go test ./...
