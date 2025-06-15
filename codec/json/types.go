// Package json implements JSON codec for JSON Patch operations.
// Provides encoding and decoding functionality for JSON Patch operations with full RFC 6902 support.
package json

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Operation represents a JSON Patch operation in JSON format.
// This unified structure supports all standard and extended operation internal.
type Operation = internal.Operation

// CompactOperation represents a compact format operation.
type CompactOperation = internal.CompactOperation

// JsonPatchOptions contains options for JSON Patch operations.
type JsonPatchOptions = internal.JsonPatchOptions
