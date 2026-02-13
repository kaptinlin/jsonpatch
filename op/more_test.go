package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestMore_Basic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		doc           any
		path          []string
		value         float64
		expectError   bool
		expectedError error
	}{
		{
			name:        "greater_than_success",
			doc:         map[string]any{"score": 85.5},
			path:        []string{"score"},
			value:       80.0,
			expectError: false,
		},
		{
			name:        "greater_than_failure",
			doc:         map[string]any{"score": 25.0},
			path:        []string{"score"},
			value:       30.0,
			expectError: true,
		},
		{
			name:        "equal_failure",
			doc:         map[string]any{"score": 25.0},
			path:        []string{"score"},
			value:       25.0,
			expectError: true,
		},
		{
			name:        "integer_comparison",
			doc:         map[string]any{"age": 30},
			path:        []string{"age"},
			value:       25.0,
			expectError: false,
		},
		{
			name:          "non_numeric_value",
			doc:           map[string]any{"name": "John"},
			path:          []string{"name"},
			value:         10.0,
			expectError:   true,
			expectedError: ErrNotNumber,
		},
		{
			name:          "missing_path",
			doc:           map[string]any{"score": 85.5},
			path:          []string{"missing"},
			value:         80.0,
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
		{
			name: "nested_path",
			doc: map[string]any{
				"user": map[string]any{
					"stats": map[string]any{
						"score": 95.0,
					},
				},
			},
			path:        []string{"user", "stats", "score"},
			value:       90.0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			moreOp := NewMore(tt.path, tt.value)
			result, err := moreOp.Apply(tt.doc)

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

func TestMore_Constructor(t *testing.T) {
	t.Parallel()
	path := []string{"age"}
	value := 18.0

	moreOp := NewMore(path, value)
	if diff := cmp.Diff(path, moreOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, value, moreOp.Value, "Value")
	if got := moreOp.Op(); got != internal.OpMoreType {
		assert.Equal(t, internal.OpMoreType, got, "Op()")
	}
	if got := moreOp.Code(); got != internal.OpMoreCode {
		assert.Equal(t, internal.OpMoreCode, got, "Code()")
	}
}

func TestMore_ToJSON(t *testing.T) {
	t.Parallel()
	moreOp := NewMore([]string{"age"}, 18.0)
	got, err := moreOp.ToJSON()

	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}
	if got.Op != string(internal.OpMoreType) {
		assert.Equal(t, string(internal.OpMoreType), got.Op, "ToJSON().Op")
	}
	if got.Path != "/age" {
		assert.Equal(t, "/age", got.Path, "ToJSON().Path")
	}
	if got.Value != 18 { // Expect int, not float64
		assert.Equal(t, 18, got.Value, "ToJSON().Value")
	}
}

func TestMore_ToCompact(t *testing.T) {
	t.Parallel()
	moreOp := NewMore([]string{"age"}, 18.0)
	compact, err := moreOp.ToCompact()
	if err != nil {
		t.Errorf("ToCompact() failed: %v", err)
	}
	want := []any{internal.OpMoreCode, []string{"age"}, 18.0}
	assert.Equal(t, want, compact)
}
