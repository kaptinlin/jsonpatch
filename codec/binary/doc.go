// Package binary implements a MessagePack-based binary codec for JSON Patch operations.
//
// The binary codec uses the same operation tree as the compact codec: paths are
// encoded as segment arrays, optional false fields are omitted, and composite
// predicate children use paths relative to the containing predicate path.
package binary
