package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestMatches_Basic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		doc           any
		path          []string
		pattern       string
		ignoreCase    bool
		expectError   bool
		expectedError error
	}{
		{
			name:        "simple pattern match",
			doc:         map[string]any{"text": "hello123"},
			path:        []string{"text"},
			pattern:     "hello\\d+",
			ignoreCase:  false,
			expectError: false,
		},
		{
			name:        "pattern no match",
			doc:         map[string]any{"text": "hello"},
			path:        []string{"text"},
			pattern:     "\\d+",
			ignoreCase:  false,
			expectError: true,
		},
		{
			name:        "case sensitive match",
			doc:         map[string]any{"text": "Hello"},
			path:        []string{"text"},
			pattern:     "hello",
			ignoreCase:  false,
			expectError: true,
		},
		{
			name:        "case insensitive match",
			doc:         map[string]any{"text": "Hello"},
			path:        []string{"text"},
			pattern:     "hello",
			ignoreCase:  true,
			expectError: false,
		},
		{
			name:          "non string value",
			doc:           map[string]any{"number": 123},
			path:          []string{"number"},
			pattern:       "\\d+",
			ignoreCase:    false,
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "missing path",
			doc:           map[string]any{"text": "hello"},
			path:          []string{"missing"},
			pattern:       "hello",
			ignoreCase:    false,
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			matchesOp := NewMatches(tt.path, tt.pattern, tt.ignoreCase, nil)

			result, err := matchesOp.Apply(tt.doc)

			if tt.expectError {
				if err == nil {
					t.Error("Apply() succeeded, want error")
				}
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("Apply() error = %v, want %v", err, tt.expectedError)
				}
				if diff := cmp.Diff(internal.OpResult[any]{}, result); diff != "" {
					t.Errorf("Apply() result mismatch (-want +got):\n%s", diff)
				}
			} else {
				if err != nil {
					t.Errorf("Apply() failed: %v", err)
				}
				if result.Doc == nil {
					t.Error("Apply() result.Doc = nil, want non-nil")
				}
				if diff := cmp.Diff(tt.doc, result.Doc); diff != "" {
					t.Errorf("Apply() result.Doc mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestMatches_Constructor(t *testing.T) {
	t.Parallel()
	path := []string{"email"}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	ignoreCase := false

	matchesOp := NewMatches(path, pattern, ignoreCase, nil)

	if diff := cmp.Diff(path, matchesOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
	if matchesOp.Pattern != pattern {
		t.Errorf("Pattern = %q, want %q", matchesOp.Pattern, pattern)
	}
	if matchesOp.IgnoreCase != ignoreCase {
		t.Errorf("IgnoreCase = %v, want %v", matchesOp.IgnoreCase, ignoreCase)
	}
	if got := matchesOp.Op(); got != internal.OpMatchesType {
		t.Errorf("Op() = %v, want %v", got, internal.OpMatchesType)
	}
	if got := matchesOp.Code(); got != internal.OpMatchesCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpMatchesCode)
	}
}

func TestMatches_InvalidPattern(t *testing.T) {
	t.Parallel()
	path := []string{"email"}
	invalidPattern := `[invalid-regex`

	// Invalid patterns create a matcher that always returns false
	// (aligned with json-joy's behavior)
	matchesOp := NewMatches(path, invalidPattern, false, nil)
	if matchesOp == nil {
		t.Fatal("NewMatches() = nil, want non-nil")
	}

	result, _ := matchesOp.Test(map[string]any{"email": "test@example.com"})
	if result {
		t.Error("Test() with invalid pattern = true, want false")
	}
}

func TestMatches_ToJSON(t *testing.T) {
	t.Parallel()
	matchesOp := NewMatches([]string{"email"}, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, true, nil)

	got, err := matchesOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}

	if got.Op != string(internal.OpMatchesType) {
		t.Errorf("ToJSON().Op = %q, want %q", got.Op, string(internal.OpMatchesType))
	}
	if got.Path != "/email" {
		t.Errorf("ToJSON().Path = %q, want %q", got.Path, "/email")
	}
	if got.Value != `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` {
		t.Errorf("ToJSON().Value = %v, want pattern string", got.Value)
	}
	if got.IgnoreCase != true {
		t.Errorf("ToJSON().IgnoreCase = %v, want true", got.IgnoreCase)
	}
}

func TestMatches_ToCompact(t *testing.T) {
	t.Parallel()
	matchesOp := NewMatches([]string{"name"}, "john", true, nil)
	compact, err := matchesOp.ToCompact()
	if err != nil {
		t.Errorf("ToCompact() failed: %v", err)
	}
	want := []any{internal.OpMatchesCode, []string{"name"}, "john", true}
	if diff := cmp.Diff(want, compact); diff != "" {
		t.Errorf("ToCompact() mismatch (-want +got):\n%s", diff)
	}
}

func TestMatches_ToCompact_WithoutIgnoreCase(t *testing.T) {
	t.Parallel()
	matchesOp := NewMatches([]string{"name"}, "john", false, nil)
	compact, err := matchesOp.ToCompact()
	if err != nil {
		t.Errorf("ToCompact() failed: %v", err)
	}
	want := []any{internal.OpMatchesCode, []string{"name"}, "john", false}
	if diff := cmp.Diff(want, compact); diff != "" {
		t.Errorf("ToCompact() mismatch (-want +got):\n%s", diff)
	}
}
