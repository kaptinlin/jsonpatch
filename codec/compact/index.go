// Package compact implements a compact array-based codec for JSON Patch operations.
// This codec provides significant space savings compared to standard JSON format
// while maintaining full compatibility with all operation types.
package compact

// Re-export main functions for convenience
var (
	// NewEncoder creates a new compact encoder with the given options
	NewCompactEncoder = NewEncoder

	// NewDecoder creates a new compact decoder with the given options
	NewCompactDecoder = NewDecoder

	// EncodeCompact encodes operations into compact format using default options
	EncodeCompact = Encode

	// DecodeCompact decodes compact format operations using default options
	DecodeCompact = Decode

	// EncodeCompactJSON encodes operations into compact format JSON bytes
	EncodeCompactJSON = EncodeJSON

	// DecodeCompactJSON decodes compact format JSON bytes into operations
	DecodeCompactJSON = DecodeJSON
)
