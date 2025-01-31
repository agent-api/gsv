# Displays this help message
help:
  @echo "Available commands:"
  @just --list

# Runs all tests
test: unit-test

# Runs all in-code, unit tests
unit-test:
  go test ./...

# Lints Go code via golangci-lint within Docker
lint:
  docker run \
    -t \
    --rm \
    -v "$(pwd)/:/app" \
    -w /app \
    golangci/golangci-lint:v1.60 \
    golangci-lint run -v

# Formats Go code via goimports
format:
  find . -type f -name "*.go" -exec goimports -local github.com/open-sauced/pizza-cli -w {} \;

# Installs the dev tools for working with this project. Requires "go", "just", and "docker"
install-dev-tools:
  #!/usr/bin/env sh

  go install golang.org/x/tools/cmd/goimports@latest

# Runs all the dev tasks (like formatting, linting, building, etc.)
dev: format lint test
