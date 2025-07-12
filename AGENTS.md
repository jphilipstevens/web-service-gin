# Agent Guidelines

## Overview
This repository hosts a Go web service template built with the Gin framework. Application code resides under `app/` with domain packages (e.g., `albums`), infrastructure modules (`db`, `cache`, middleware), and dependency wiring. Example usage and sample configuration live under `example/`. Configuration helpers can be found in `config/`.

## Architecture Notes
- `app/` contains well-defined modules. Each module exposes a minimal interface and keeps implementation details private. For example, the `cache` package defines an interface and a Redis implementation while isolating tracing and error mapping logic.
- `example/` demonstrates domain features and provides an example `config.yaml` used by tests.
- Tests reside alongside packages using Go's standard `*_test.go` pattern. Utility test helpers are under `testUtils/`.

## Development Practices
- Follow the principles from *A Philosophy of Software Design*: prefer deep modules with simple interfaces, hide implementation details, and keep each package focused on a single responsibility.
- Optimize for clarity and changeability over cleverness. Avoid temporal coupling and aim for composable designs.
- Keep modules small but meaningful. When adding new functionality, consider whether it belongs in an existing package or merits its own package.

## Workflow Requirements
1. **Code formatting**: run `gofmt -w $(find . -name '*.go' -not -path './vendor/*')` before committing.
2. **Testing**: run `go test ./...` and ensure all tests pass.
3. **Add tests** for new functionality when feasible.

Adhering to these guidelines will help maintain a clean, easy to understand codebase suitable for expert maintainers.
