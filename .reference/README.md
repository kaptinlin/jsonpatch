# JSON Patch Reference

This directory contains reference materials for the Go implementation of JSON Patch, including the original TypeScript source code from json-joy.

## Directory Structure

```
.reference/
├── json-joy/          # Original TypeScript source (git submodule)
├── scripts/           # Development scripts
└── README.md          # This file
```

## Quick Start

```bash
# Initialize the TypeScript reference
git submodule update --init --recursive

# Update to latest version
./.reference/scripts/update-reference.sh

# Check for API changes
./.reference/scripts/check-api-changes.sh
```

## TypeScript Source Mapping

| TypeScript | Go Implementation |
|------------|-------------------|
| `src/json-patch/index.ts` | `jsonpatch.go` |
| `src/json-patch/types.ts` | `internal/` |
| `src/json-patch/applyPatch/` | Core API functions |
| `src/json-patch/op/` | Operation implementations |

## Scripts

- `update-reference.sh` - Update the json-joy submodule to latest version
- `check-api-changes.sh` - Check for API compatibility changes

## Resources

- [json-joy Repository](https://github.com/streamich/json-joy)
- [JSON Patch RFC 6902](https://tools.ietf.org/html/rfc6902) 
