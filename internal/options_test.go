package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithMutate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		val  bool
	}{
		{"enable mutate", true},
		{"disable mutate", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var opts Options
			WithMutate(tt.val)(&opts)
			if opts.Mutate != tt.val {
				assert.Fail(t, fmt.Sprintf("WithMutate(%v): got %v, want %v",
					tt.val, opts.Mutate, tt.val))
			}
		})
	}
}

func TestWithMatcher(t *testing.T) {
	t.Parallel()
	called := false
	matcher := func(_ string, _ bool) RegexMatcher {
		called = true
		return func(_ string) bool { return true }
	}

	var opts Options
	WithMatcher(matcher)(&opts)

	if opts.CreateMatcher == nil {
		require.FailNow(t, "WithMatcher: CreateMatcher is nil")
	}

	opts.CreateMatcher("test", false)
	if !called {
		assert.Fail(t, "WithMatcher: factory was not called")
	}
}
