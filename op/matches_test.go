package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
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
					assert.Fail(t, "Apply() succeeded, want error")
				}
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					assert.Equal(t, tt.expectedError, err, "Apply() error")
				}
				assert.Equal(t, internal.OpResult[any]{}, result)
			} else {
				if err != nil {
					t.Errorf("Apply() failed: %v", err)
				}
				if result.Doc == nil {
					assert.Fail(t, "Apply() result.Doc = nil, want non-nil")
				}
				assert.Equal(t, tt.doc, result.Doc)
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
		assert.Equal(t, pattern, matchesOp.Pattern, "Pattern")
	}
	assert.Equal(t, ignoreCase, matchesOp.IgnoreCase, "IgnoreCase")
	if got := matchesOp.Op(); got != internal.OpMatchesType {
		assert.Equal(t, internal.OpMatchesType, got, "Op()")
	}
	if got := matchesOp.Code(); got != internal.OpMatchesCode {
		assert.Equal(t, internal.OpMatchesCode, got, "Code()")
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
		assert.Fail(t, "Test() with invalid pattern = true, want false")
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
		assert.Equal(t, string(internal.OpMatchesType), got.Op, "ToJSON().Op")
	}
	if got.Path != "/email" {
		assert.Equal(t, "/email", got.Path, "ToJSON().Path")
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
	assert.Equal(t, want, compact)
}

func TestMatches_ToCompact_WithoutIgnoreCase(t *testing.T) {
	t.Parallel()
	matchesOp := NewMatches([]string{"name"}, "john", false, nil)
	compact, err := matchesOp.ToCompact()
	if err != nil {
		t.Errorf("ToCompact() failed: %v", err)
	}
	want := []any{internal.OpMatchesCode, []string{"name"}, "john", false}
	assert.Equal(t, want, compact)
}
