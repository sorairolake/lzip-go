# SPDX-FileCopyrightText: 2024 Shun Sakai
#
# SPDX-License-Identifier: Apache-2.0 OR MIT

alias all := default

# Run default recipe
default: test

# Remove generated artifacts
@clean:
    go clean

# Run tests
@test:
    go test ./...

# Run `golangci-lint run`
@golangci-lint:
    golangci-lint run -E gofmt,goimports

# Run the formatter
fmt: gofmt goimports

# Run `go fmt`
@gofmt:
    go fmt ./...

# Run `goimports`
@goimports:
    fd -e go -x goimports -w

# Run the linter
lint: vet staticcheck

# Run `go vet`
@vet:
    go vet ./...

# Run `staticcheck`
@staticcheck:
    staticcheck ./...

# Run the linter for GitHub Actions workflow files
@lint-github-actions:
    actionlint -verbose

# Run the formatter for the README
@fmt-readme:
    npx prettier -w README.md

# Increment the version
@bump part:
    bump-my-version bump {{part}}
