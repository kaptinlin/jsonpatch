# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

- **Module:** `github.com/kaptinlin/jsonpatch`
- **Go Version:** 1.26
- Canonical technical contracts live in `SPECS/`.
- Keep `README.md` user-facing and keep `CLAUDE.md` operational.

## Commands

```bash
# Run all tests with race detection
task test

# Run linter (golangci-lint + tidy check)
task lint

# Run markdown lint for docs and specs
task markdownlint

# Format Go code
task fmt

# Run go vet
task vet

# Full verification
task verify
```

### Targeted Commands

```bash
# Test specific operation implementations
go test -race ./op -run TestAddOp
go test -race ./op -run TestRemoveOp

# Test codec implementations
go test -race ./codec/json/...
go test -race ./codec/compact/...
go test -race ./codec/binary/...

# Benchmark a specific package
go test -bench=. -benchmem ./codec/json
```

## SPECS Index

- `SPECS/00-overview.md` — library scope, mutation default, supported document shapes
- `SPECS/20-api-specs.md` — public entry points and RFC 6902 mutating operations
- `SPECS/25-predicate-specs.md` — predicate and second-order predicate behavior
- `SPECS/26-extended-operation-specs.md` — extended operation contracts
- `SPECS/30-data-model-specs.md` — payload and result types
- `SPECS/40-architecture-specs.md` — package boundaries and execution pipeline
- `SPECS/50-coding-standards.md` — compatibility, errors, tests, and documentation rules

When behavior changes, update the affected spec in the same change.

## Agent Skills

Package-local skills available in `.agents/skills/`:

- **agent-md-creating** - Generate CLAUDE.md for Go projects
- **code-simplifying** - Refine code for clarity and consistency
- **committing** - Create conventional commits
- **dependency-selecting** - Select Go dependencies from kaptinlin ecosystem
- **go-best-practices** - Google Go coding best practices
- **linting** - Set up and run golangci-lint v2
- **modernizing** - Go 1.20-1.26 modernization guide
- **ralphy-initializing** - Initialize Ralphy AI coding loop
- **ralphy-todo-creating** - Create Ralphy TODO.yaml task files
- **readme-creating** - Generate README.md for Go libraries
- **releasing** - Guide release process for Go packages
- **testing** - Write Go tests with best practices

Use the `Skill` tool when the relevant skill is available in the current Claude session.
