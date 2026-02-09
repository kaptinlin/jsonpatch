package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestMore_Basic(t *testing.T) {
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
			moreOp := NewMore(tt.path, tt.value)
			result, err := moreOp.Apply(tt.doc)

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

func TestMore_Constructor(t *testing.T) {
	path := []string{"age"}
	value := 18.0

	moreOp := NewMore(path, value)
	if diff := cmp.Diff(path, moreOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
	if moreOp.Value != value {
		t.Errorf("Value = %v, want %v", moreOp.Value, value)
	}
	if got := moreOp.Op(); got != internal.OpMoreType {
		t.Errorf("Op() = %v, want %v", got, internal.OpMoreType)
	}
	if got := moreOp.Code(); got != internal.OpMoreCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpMoreCode)
	}
}

func TestMore_ToJSON(t *testing.T) {
	moreOp := NewMore([]string{"age"}, 18.0)
	got, err := moreOp.ToJSON()

	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}
	if got.Op != string(internal.OpMoreType) {
		t.Errorf("ToJSON().Op = %q, want %q", got.Op, string(internal.OpMoreType))
	}
	if got.Path != "/age" {
		t.Errorf("ToJSON().Path = %q, want %q", got.Path, "/age")
	}
	if got.Value != 18 { // Expect int, not float64
		t.Errorf("ToJSON().Value = %v, want %v", got.Value, 18)
	}
}

func TestMore_ToCompact(t *testing.T) {
	moreOp := NewMore([]string{"age"}, 18.0)
	compact, err := moreOp.ToCompact()
	if err != nil {
		t.Errorf("ToCompact() failed: %v", err)
	}
	want := []any{internal.OpMoreCode, []string{"age"}, 18.0}
	if diff := cmp.Diff(want, compact); diff != "" {
		t.Errorf("ToCompact() mismatch (-want +got):\n%s", diff)
	}
}
