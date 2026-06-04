# jsonpatch Agent Guide

This file provides operational guidance for agents working in `github.com/kaptinlin/jsonpatch`. Keep user-facing usage in `README.md`; keep technical contracts in `SPECS/`.

## Project Overview

- **Module:** `github.com/kaptinlin/jsonpatch`
- **Go version:** see `go.mod`
- **Contract source:** `SPECS/`
- **Reference evidence:** `.reference/json-joy/`
- **Core model:** compile operation vocabulary into a reusable `Patch`, then apply it immutably with `Apply` or explicitly in place with `ApplyInPlace`.

## Commands

```bash
task test           # Run all tests with race detection
task golangci-lint  # Run golangci-lint only
task lint           # Run golangci-lint and tidy checks
task fmt            # Format Go code
task vet            # Run go vet
task bench          # Run benchmarks
task verify         # Run deps, fmt, vet, lint, test, vuln
```

## Architecture

| Path | Responsibility |
| --- | --- |
| `patch.go`, `errors.go`, `index.go`, `util.go` | Compiled public API, structured errors, operation constants, document-shape dispatch, compile-time capability policy |
| `op/` | Executable operation implementations and shared apply helpers |
| `internal/` | Shared interfaces, constants, and codec payload types |
| `codec/json/` | JSON-shaped operation encoding and decoding |
| `codec/compact/` | Compact array codec with segment-array paths |
| `codec/binary/` | MessagePack codec for executable operations |
| `tests/` | Scenario, predicate, reference, benchmark, and test utility coverage |
| `examples/` | Runnable user-facing examples |

## Agent Workflow

### Design Phase - Read SPECS First

Before designing or modifying behavior, read the relevant `SPECS/` documents. SPECS record tested contracts; source code and tests remain the executable truth.

1. Identify relevant specs from the SPECS Index.
2. Read those specs before editing behavior.
3. Implement the correct behavior with focused tests.
4. Update the owning spec when the change becomes a lasting contract.
5. Ask the user only when code, tests, specs, and request still leave a risky product decision unresolved.

### Implementation Phase - Check Reference Evidence

When reference behavior matters, inspect `.reference/README.md` and `.reference/json-joy/` before changing code. Treat the TypeScript reference as evidence for vocabulary and edge cases, not as authority over Go API shape.

## SPECS Index

| Spec | Owns |
| --- | --- |
| `SPECS/00-overview.md` | Library scope, mutation default, supported document shapes |
| `SPECS/20-api-specs.md` | Public entry points, compile options, capabilities, RFC 6902 mutating operations, error contract |
| `SPECS/25-predicate-specs.md` | Predicate and second-order predicate behavior |
| `SPECS/26-extended-operation-specs.md` | Extended operation contracts |
| `SPECS/30-data-model-specs.md` | JSON payload, compact payload, type vocabulary, document/result/error types |
| `SPECS/40-architecture-specs.md` | Package boundaries, execution pipeline, codec wire contract |
| `SPECS/50-coding-standards.md` | Implementation, errors, tests, and documentation rules |

## References Index

| Reference | Use |
| --- | --- |
| `.reference/README.md` | Local reference index and maintenance scripts |
| `.reference/json-joy/` | TypeScript JSON Patch+ vocabulary, predicates, extended operations, and codec evidence |

## Agent Operating Rules

- Read nearby code and relevant specs before editing.
- Prefer one clear public story over ambiguous fallback paths.
- Keep edits surgical and avoid unrelated refactors.
- Preserve user changes already in the working tree.
- Prove behavior with tests, not spec-mirror tests.
- Do not add policy-only gates that restate docs or specs.
- Fail loudly with sentinel errors and wrapped context.
- Treat reference projects as evidence, not authority.
- Respect context budgets; summarize long docs instead of copying them into prompts.

## Design Philosophy

- **KISS** - One lifecycle: compile a `Patch`, then `Apply` or `ApplyInPlace`.
- **SRP** - Root owns type-preserving dispatch and compile policy; `op/` owns behavior; `codec/*` owns wire formats.
- **DRY** - Operation vocabulary, codec shapes, and spec contracts should not drift into duplicate truths.
- **YAGNI** - Do not add planners, managers, plugins, or codegen public API without real consumers.
- **Simplicity as art** - Dangerous actions are verbs, not small options; default application stays copy-safe.
- **Never:** accidental complexity, compatibility shims, abstraction theater, or configurability as a substitute for a clear API.

## API Design Principles

- **Progressive disclosure:** use `Compile`/`CompileJSON`, then `Apply`; reach for `ApplyInPlace` only when mutation is intentional; use `codec/*` only for wire-format control.
- **Capability honesty:** default compilation accepts RFC 6902 only; predicates, regex predicates, and extended operations require explicit capabilities.
- **Explicit document shape:** plain `string` is scalar text; use `JSONText` or `[]byte` for JSON text.
- **Structured failure:** callers match stable failure classes with `errors.Is` and inspect `*Error` for index, op, path, from, codec, and cause.

## Coding Rules

### Must Follow

- Use the Go version declared in `go.mod`; reach for newer features only when they simplify code.
- Follow Google Go Best Practices: https://google.github.io/go-style/best-practices
- Follow Google Go Style Decisions: https://google.github.io/go-style/decisions
- Preserve the caller's document type whenever conversion back is safe.
- Keep mutation opt-in through `ApplyInPlace`.
- Keep the library pure: return errors and leave logging to callers.
- Use sentinel errors for stable failure classes and match with `errors.Is`.
- Keep compile and execution errors separate: compile checks payload shape and capability policy; execution checks document state.
- Record durable behavior in `SPECS/`; keep `README.md` user-facing and `CLAUDE.md` operational.

### Go 1.26 Patterns In Use

| Pattern | Where |
| --- | --- |
| Generics and type-preserving constraints | Root `Apply`, `Result`, `Document` |
| `testing.B.Loop()` | Benchmarks |
| `for range N` | Test and helper loops |
| `maps.Clone` / `maps.Copy` | Examples and tests |

## Testing

- Use table-driven tests for multi-case behavior.
- Use `t.Parallel()` for top-level tests and independent subtests.
- Add success and failure coverage for new behavior.
- Protect codec wire contracts with focused golden tests.
- Do not compare error strings; use `errors.Is` and `errors.As`.
- Add example-oriented tests when `README.md` or public usage snippets change.
- Run `task test` and `task lint` after code or docs changes; report tidy-only failures separately from lint failures.

## Dependencies

- `github.com/go-json-experiment/json` - JSON parsing and marshaling
- `github.com/kaptinlin/deepclone` - clone support for immutable application
- `github.com/kaptinlin/jsonpointer` - JSON Pointer formatting and traversal
- `github.com/tinylib/msgp` - MessagePack support for the binary codec

## Dependency Issue Reporting

When you encounter a bug, limitation, or unexpected behavior in a dependency library:

1. Do not work around it by reimplementing the dependency's functionality.
2. Do not skip or silently replace the dependency with project-local code.
3. Create `reports/<dependency-name>.md`.
4. Include dependency version, trigger scenario, expected behavior, actual behavior, relevant errors, and a non-implemented workaround idea.
5. Continue tasks that are unaffected by the dependency issue.

## Agent Skills

Use the matching local skill in `.agents/skills/` when it fits the task.

| Skill | When to Use |
| --- | --- |
| `agent-md-writing` | Refresh `CLAUDE.md` and the `AGENTS.md` symlink |
| `readme-writing` | Refresh user-facing README usage documentation |
| `library-docs-maintaining` | Refresh top-level library docs together |
| `library-specs-maintaining` | Consolidate or maintain SPECS documents |
| `library-test-covering` | Add or expand behavior coverage |
| `library-legacy-pruning` | Delete deprecated APIs, legacy shims, and old compatibility paths |
| `golangci-linting` | Configure or troubleshoot golangci-lint v2 |
| `go-best-practices` | Review Go API and implementation style |
| `code-simplifying` | Simplify recently changed code without changing behavior |
| `committing` | Prepare conventional commits |
| `releasing` | Guide release work |

## Forbidden

- No documentation masquerading as code; do not encode spec prose in values or helpers that no program consumes.
- No policy-only gate scripts whose only job is to restate docs or specs.
- No spec mirror tests when behavior is already covered by stronger tests.
- No working around dependency bugs; report them in `reports/`.
- No assuming `string` input is a JSON document; use `JSONText` or `[]byte`.
- No implicit mutation; use `ApplyInPlace`.
- No reintroducing legacy helper entry points or runtime mutation options.
- No `panic` or logging in library code.
- No untested API or architecture promises in `README.md` or `CLAUDE.md`.
