package internal

import "testing"

func TestWithMutate(t *testing.T) {
	tests := []struct {
		name string
		val  bool
	}{
		{"enable mutate", true},
		{"disable mutate", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts Options
			WithMutate(tt.val)(&opts)
			if opts.Mutate != tt.val {
				t.Errorf("WithMutate(%v): got %v, want %v",
					tt.val, opts.Mutate, tt.val)
			}
		})
	}
}

func TestWithMatcher(t *testing.T) {
	called := false
	matcher := func(_ string, _ bool) RegexMatcher {
		called = true
		return func(_ string) bool { return true }
	}

	var opts Options
	WithMatcher(matcher)(&opts)

	if opts.CreateMatcher == nil {
		t.Fatal("WithMatcher: CreateMatcher is nil")
	}

	opts.CreateMatcher("test", false)
	if !called {
		t.Error("WithMatcher: factory was not called")
	}
}
