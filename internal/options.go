package internal

// JSONPatchOptions configures JSON Patch decoding.
type JSONPatchOptions struct {
	// CreateMatcher overrides regex compilation for pattern operations.
	CreateMatcher CreateRegexMatcher
}
