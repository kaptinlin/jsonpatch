// Package binary implements a MessagePack-based binary codec for JSON Patch operations.
//
// Limitations:
//   - Second-order predicates (and, or, not) are NOT supported in binary codec.
//     These operations return an error during encoding.
//     Use the JSON or compact codec if you need second-order predicate support.
package binary
