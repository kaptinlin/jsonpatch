package internal

// JSONPatchOptions contains options for JSON Patch decoding.
type JSONPatchOptions struct {
	CreateMatcher CreateRegexMatcher
}

// WithMutate sets whether the patch should modify the original document.
func WithMutate(mutate bool) Option {
	return func(o *Options) {
		o.Mutate = mutate
	}
}

// WithMatcher sets a custom regex matcher factory for pattern operations.
func WithMatcher(createMatcher CreateRegexMatcher) Option {
	return func(o *Options) {
		o.CreateMatcher = createMatcher
	}
}
