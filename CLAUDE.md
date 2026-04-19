# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

- **Module:** `github.com/kaptinlin/jsonpatch`
- **Go Version:** 1.26.2
- Canonical technical contracts live in `SPECS/`.
- Keep `README.md` user-facing and keep `CLAUDE.md` operational.
- The library applies RFC 6902 operations, predicate operations, and extended operations while preserving the caller's document type whenever the result can be converted back safely.

## Commands

```bash
# Run all tests with race detection
task test

# Run golangci-lint and tidy checks
task lint

# Run markdownlint for docs and specs
task markdownlint

# Format Go code
task fmt

# Run go vet
task vet
```

## Architecture

| Path | Responsibility |
| --- | --- |
| `jsonpatch.go`, `index.go`, `validate.go` | Public API, document-shape dispatch, validation |
| `op/` | Executable operation implementations |
| `internal/` | Shared contracts, options, payload and result types |
| `codec/json/` | JSON operation encoding and decoding |
| `codec/compact/` | Compact array codec |
| `codec/binary/` | MessagePack codec |
| `tests/` | Scenario and predicate coverage |
| `examples/` | Runnable usage examples |

## Agent Workflow

### Design Phase — Read SPECS First

Before designing or modifying behavior, read the relevant `SPECS/` documents first.
SPECS define the current contracts for API behavior, compatibility, architecture, testing, and documentation ownership.
Do not invent new behavior that contradicts `SPECS/`.

**Workflow:**

1. Identify the relevant specs from the index below.
2. Read those specs completely before changing behavior.
3. Update the owning spec in the same change when behavior changes.
4. Ask the user before proceeding if the current specs do not cover the case.

### Implementation Phase — Check Compatibility References

When behavior needs to match json-joy, inspect `.reference/README.md` and `.reference/json-joy/` before changing code.
Use the current Go implementation as the source of truth when the reference submodule is not initialized.

## SPECS Index

- `SPECS/00-overview.md` — library scope, mutation default, supported document shapes
- `SPECS/20-api-specs.md` — public entry points and RFC 6902 mutating operations
- `SPECS/25-predicate-specs.md` — predicate and second-order predicate behavior
- `SPECS/26-extended-operation-specs.md` — extended operation contracts
- `SPECS/30-data-model-specs.md` — payload and result types
- `SPECS/40-architecture-specs.md` — package boundaries and execution pipeline
- `SPECS/50-coding-standards.md` — compatibility, errors, tests, and documentation rules

## Design Philosophy

- **KISS** — Keep one public story for each layer: JSON-shaped operations through `ApplyPatch`, executable operations through `ApplyOp` and `ApplyOps`, codec work through `codec/*`.
- **SRP** — The root package owns type-preserving dispatch, `op/` owns operation behavior, and `codec/*` owns wire formats.
- **Simplicity as art** — Mutation is opt-in through `WithMutate(true)`. The default path stays predictable and copy-safe.
- **Errors as teachers** — Stable sentinel errors and wrapped context should tell callers whether they hit bad payloads, invalid pointers, or runtime document-state failures.
- **APIs as language** — Prefer operation names and option names that read like patch vocabulary, not framework plumbing.
- **Never:** accidental complexity, feature gravity, abstraction theater, configurability cope.

## API Design Principles

- **Progressive Disclosure** — Use `ApplyPatch` for JSON-shaped operations, `ApplyOp` and `ApplyOps` for executable operations, and `codec/*` only when callers need wire-format control.

## Coding Rules

### Must Follow

- Use Go 1.26.2 features when they simplify code.
- Follow [Google Go Best Practices](https://google.github.io/go-style/best-practices).
- Follow [Google Go Style Decisions](https://google.github.io/go-style/decisions).
- Preserve the caller's document type whenever the result can be converted back safely.
- Keep mutation opt-in. Require `WithMutate(true)` for in-place updates.
- Keep the library pure: return errors and leave logging to callers.
- Use sentinel errors for stable failure classes and match with `errors.Is`.
- Keep technical behavior in `SPECS/`, not in `README.md` or `CLAUDE.md`.

### Domain Patterns

See `SPECS/` for detailed rules:

- `SPECS/20-api-specs.md` — public API contracts
- `SPECS/25-predicate-specs.md` — predicate and negation rules
- `SPECS/26-extended-operation-specs.md` — extended operation contracts
- `SPECS/40-architecture-specs.md` — package boundaries and execution pipeline
- `SPECS/50-coding-standards.md` — testing, error, and documentation rules

## Testing

- Use table-driven tests for multi-case behavior.
- Use `t.Parallel()` for top-level tests and independent subtests.
- Add success and failure coverage for each new behavior.
- Run `task test`, `task lint`, and `task markdownlint` after code or docs changes.
- Add example-oriented tests when `README.md` or public usage snippets change.

## Dependencies

- `github.com/go-json-experiment/json` — JSON parsing and marshaling
- `github.com/kaptinlin/deepclone` — clone support for immutable application
- `github.com/kaptinlin/jsonpointer` — JSON Pointer formatting and traversal
- `github.com/tinylib/msgp` — MessagePack support for the binary codec

## Dependency Issue Reporting

When you encounter a bug, limitation, or unexpected behavior in a dependency library:

1. **Do NOT** work around it by reimplementing the dependency's functionality.
2. **Do NOT** skip the dependency and silently replace it with project-local code.
3. **Do** create a report file: `reports/<dependency-name>.md`.
4. **Do** include the dependency version, trigger scenario, expected behavior, actual behavior, errors, and a non-implemented workaround idea.
5. **Do** continue with tasks that are unaffected by the dependency issue.

## Agent Skills

Use the matching local skill in `.agents/skills/` (mirrored in `.claude/skills/`) when it fits the task:

- `agent-md-writing` — regenerate `CLAUDE.md` and refresh the `AGENTS.md` symlink
- `readme-writing` — regenerate `README.md`
- `library-docs-maintaining` — refresh top-level library docs together
- `library-specs-maintaining` — update `SPECS/` contracts
- `library-test-covering` — add or expand test coverage
- `golangci-linting` — configure or troubleshoot linting
- `go-best-practices` — review Go API and implementation style
- `dependency-selecting` — choose kaptinlin ecosystem dependencies
- `committing` — prepare a conventional commit
- `releasing` — guide release work

## Forbidden

- No documentation masquerading as code — do not encode spec prose in values or helpers that no program reads at runtime.
- No working around dependency bugs — report them in `reports/` instead.
- No assuming every `string` input is a JSON document; only strings starting with `{` or `[` are parsed as JSON.
- No implicit mutation — use `WithMutate(true)` when mutation is required.
- No `panic` or logging in library code.
- No moving canonical behavior from `SPECS/` back into `README.md` or `CLAUDE.md`.
