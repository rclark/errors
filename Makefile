.PHONY: init doc test

init:
	go mod tidy
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@v1.1.0
	go install gotest.tools/gotestsum@v1.12.0
	go install github.com/axw/gocov/gocov@v1.1.0
	go install github.com/matm/gocov-html/cmd/gocov-html@v1.4.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	@if [ ! -e .git/hooks/pre-commit ]; then \
		chmod 755 .githooks/pre-commit; \
		ln -s $(PWD)/.githooks/pre-commit .git/hooks/pre-commit; \
	fi

doc:
	@gomarkdoc \
		--output readme.md  \
		--template-file file=templates/file.md \
		--template-file package=templates/package.md \
		.

test:
	gotestsum --format testname -- -coverprofile=coverage.out ./...
	gocov convert coverage.out | gocov-html > coverage.html
