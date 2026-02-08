// Package json implements JSON codec for JSON Patch operations.
// Provides encoding and decoding functionality for JSON Patch operations with full RFC 6902 support.
package json

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Operation represents a JSON Patch operation in JSON format.
// This unified structure supports all standard and extended operation types.
type Operation = internal.Operation

// PatchOptions contains options for JSON Patch operations.
type PatchOptions = internal.JSONPatchOptions
