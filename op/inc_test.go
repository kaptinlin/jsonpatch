package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestInc_Apply(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		path     []string
		doc      any
		inc      float64
		expected any
		oldValue any
		wantErr  bool
	}{
		{
			name:     "inc int field",
			path:     []string{"count"},
			doc:      map[string]any{"count": 1},
			inc:      2,
			expected: map[string]any{"count": 3.0},
			oldValue: 1,
		},
		{
			name:     "inc float field",
			path:     []string{"score"},
			doc:      map[string]any{"score": 1.5},
			inc:      0.5,
			expected: map[string]any{"score": 2.0},
			oldValue: 1.5,
		},
		{
			name:     "dec int field",
			path:     []string{"count"},
			doc:      map[string]any{"count": 5},
			inc:      -3,
			expected: map[string]any{"count": 2.0},
			oldValue: 5,
		},
		{
			name:     "inc nested field",
			path:     []string{"user", "age"},
			doc:      map[string]any{"user": map[string]any{"age": 20}},
			inc:      1,
			expected: map[string]any{"user": map[string]any{"age": 21.0}},
			oldValue: 20,
		},
		{
			name:     "inc array element",
			path:     []string{"nums", "1"},
			doc:      map[string]any{"nums": []any{1, 2, 3}},
			inc:      10,
			expected: map[string]any{"nums": []any{1, 12.0, 3}},
			oldValue: 2,
		},
		{
			name:     "inc root int",
			path:     []string{},
			doc:      100,
			inc:      23,
			expected: 123.0,
			oldValue: 100.0,
		},
		{
			name:     "inc root float",
			path:     []string{},
			doc:      1.5,
			inc:      0.5,
			expected: 2.0,
			oldValue: 1.5,
		},
		{
			name:     "path not found creates value",
			path:     []string{"notfound"},
			doc:      map[string]any{"count": 1},
			inc:      5,
			expected: map[string]any{"count": 1, "notfound": 5.0},
			oldValue: nil,
		},
		{
			name:    "not a number",
			path:    []string{"str"},
			doc:     map[string]any{"str": "abc"},
			inc:     1,
			wantErr: true,
		},
		{
			name:    "root not a number",
			path:    []string{},
			doc:     "abc",
			inc:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			incOp := NewInc(tt.path, tt.inc)
			docCopy, err := DeepClone(tt.doc)
			if err != nil {
				t.Fatalf("DeepClone() error: %v", err)
			}

			result, err := incOp.Apply(docCopy)

			if tt.wantErr {
				if err == nil {
					t.Error("Apply() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Apply() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.expected, result.Doc); diff != "" {
				t.Errorf("Apply() Doc mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.oldValue, result.Old); diff != "" {
				t.Errorf("Apply() Old mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestInc_Constructor(t *testing.T) {
	t.Parallel()
	path := []string{"user", "score"}
	inc := 3.5
	incOp := NewInc(path, inc)
	if diff := cmp.Diff(path, incOp.path); diff != "" {
		t.Errorf("NewInc() path mismatch (-want +got):\n%s", diff)
	}
	if incOp.Inc != inc {
		t.Errorf("NewInc() Inc = %v, want %v", incOp.Inc, inc)
	}
	if got := incOp.Op(); got != internal.OpIncType {
		t.Errorf("Op() = %v, want %v", got, internal.OpIncType)
	}
	if got := incOp.Code(); got != internal.OpIncCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpIncCode)
	}
}

func TestInc_ToJSON(t *testing.T) {
	t.Parallel()
	incOp := NewInc([]string{"count"}, 5.5)
	got, err := incOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	if got.Op != "inc" {
		t.Errorf("ToJSON() Op = %q, want %q", got.Op, "inc")
	}
	if got.Path != "/count" {
		t.Errorf("ToJSON() Path = %q, want %q", got.Path, "/count")
	}
	if got.Inc != 5.5 {
		t.Errorf("ToJSON() Inc = %v, want %v", got.Inc, 5.5)
	}
}
