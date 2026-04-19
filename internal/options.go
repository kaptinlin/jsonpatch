package internal

// JSONPatchOptions configures JSON Patch decoding.
type JSONPatchOptions struct {
	// CreateMatcher overrides regex compilation for pattern operations.
	CreateMatcher CreateRegexMatcher
}

// WithMutate enables or disables in-place mutation.
func WithMutate(mutate bool) Option {
	return func(o *Options) {
		o.Mutate = mutate
	}
}

// WithMatcher sets the regex matcher factory used by pattern operations.
func WithMatcher(createMatcher CreateRegexMatcher) Option {
	return func(o *Options) {
		o.CreateMatcher = createMatcher
	}
}
