// Package binary implements a MessagePack-based binary codec for JSON Patch operations.
//
// Limitations:
//   - Second-order predicates (and, or, not) are NOT supported in binary codec.
//     These operations are skipped during encoding with a warning.
//     Use the JSON or compact codec if you need second-order predicate support.
package binary
