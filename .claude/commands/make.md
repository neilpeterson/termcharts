---
description: Run make targets (build, test, cover, lint, install, release, clean, help, deps)
allowed-tools: Bash(make:*)
argument-hint: [target]
---

Run the make build process for the termcharts project.

Available make targets:
- build: Build the binary
- test: Run tests
- cover: Run tests with coverage report
- lint: Run golangci-lint
- install: Install binary to $GOPATH/bin
- deps: Download and tidy dependencies
- release: Cross-compile for all platforms
- clean: Remove build artifacts
- help: Display help message

If arguments are provided, run `make $ARGUMENTS`.
If no arguments provided, run `make build`.

Show the output to the user.
