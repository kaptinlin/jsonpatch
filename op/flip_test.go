package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestFlip_Apply(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		doc      any
		expected any
		oldValue any
		wantErr  bool
	}{
		{
			name:     "boolean true to false",
			path:     []string{"flag"},
			doc:      map[string]any{"flag": true},
			expected: map[string]any{"flag": false},
			oldValue: true,
		},
		{
			name:     "boolean false to true",
			path:     []string{"flag"},
			doc:      map[string]any{"flag": false},
			expected: map[string]any{"flag": true},
			oldValue: false,
		},
		{
			name:     "number 0 to true",
			path:     []string{"count"},
			doc:      map[string]any{"count": 0},
			expected: map[string]any{"count": true},
			oldValue: 0,
		},
		{
			name:     "number 5 to false",
			path:     []string{"count"},
			doc:      map[string]any{"count": 5},
			expected: map[string]any{"count": false},
			oldValue: 5,
		},
		{
			name:     "empty string to true",
			path:     []string{"text"},
			doc:      map[string]any{"text": ""},
			expected: map[string]any{"text": true},
			oldValue: "",
		},
		{
			name:     "non-empty string to false",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello"},
			expected: map[string]any{"text": false},
			oldValue: "hello",
		},
		{
			name:     "nil to true",
			path:     []string{"value"},
			doc:      map[string]any{"value": nil},
			expected: map[string]any{"value": true},
			oldValue: nil,
		},
		{
			name:     "empty array to false",
			path:     []string{"items"},
			doc:      map[string]any{"items": []any{}},
			expected: map[string]any{"items": false},
			oldValue: []any{},
		},
		{
			name:     "non-empty array to false",
			path:     []string{"items"},
			doc:      map[string]any{"items": []any{1, 2, 3}},
			expected: map[string]any{"items": false},
			oldValue: []any{1, 2, 3},
		},
		{
			name:     "empty map to false",
			path:     []string{"config"},
			doc:      map[string]any{"config": map[string]any{}},
			expected: map[string]any{"config": false},
			oldValue: map[string]any{},
		},
		{
			name:     "non-empty map to false",
			path:     []string{"config"},
			doc:      map[string]any{"config": map[string]any{"key": "value"}},
			expected: map[string]any{"config": false},
			oldValue: map[string]any{"key": "value"},
		},
		{
			name:     "nested path",
			path:     []string{"user", "active"},
			doc:      map[string]any{"user": map[string]any{"active": true}},
			expected: map[string]any{"user": map[string]any{"active": false}},
			oldValue: true,
		},
		{
			name:     "array element",
			path:     []string{"flags", "0"},
			doc:      map[string]any{"flags": []any{true, false, true}},
			expected: map[string]any{"flags": []any{false, false, true}},
			oldValue: true,
		},
		{
			name:     "root level boolean",
			path:     []string{},
			doc:      true,
			expected: false,
			oldValue: true,
		},
		{
			name:     "root level number",
			path:     []string{},
			doc:      42,
			expected: false,
			oldValue: 42,
		},
		{
			name:     "path not found creates true",
			path:     []string{"nonexistent"},
			doc:      map[string]any{"flag": true},
			expected: map[string]any{"flag": true, "nonexistent": true},
			oldValue: nil,
		},
		{
			name:    "invalid path for array",
			path:    []string{"items", "invalid"},
			doc:     map[string]any{"items": []any{1, 2, 3}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flipOp := NewFlip(tt.path)

			docCopy, err := DeepClone(tt.doc)
			if err != nil {
				t.Fatalf("DeepClone() error: %v", err)
			}

			result, err := flipOp.Apply(docCopy)

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

func TestFlip_Constructor(t *testing.T) {
	path := []string{"user", "active"}
	flipOp := NewFlip(path)

	if diff := cmp.Diff(path, flipOp.Path()); diff != "" {
		t.Errorf("NewFlip() Path mismatch (-want +got):\n%s", diff)
	}
	if got := flipOp.Op(); got != internal.OpFlipType {
		t.Errorf("Op() = %v, want %v", got, internal.OpFlipType)
	}
	if got := flipOp.Code(); got != internal.OpFlipCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpFlipCode)
	}
}

func TestFlip_ComplexTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{"float64 zero", 0.0, true},
		{"float64 non-zero", 3.14, false},
		{"int8 zero", int8(0), true},
		{"int8 non-zero", int8(1), false},
		{"uint zero", uint(0), true},
		{"uint non-zero", uint(1), false},
		{"float32 zero", float32(0.0), true},
		{"float32 non-zero", float32(1.0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flipOp := NewFlip([]string{"value"})
			doc := map[string]any{"value": tt.value}

			result, err := flipOp.Apply(doc)
			if err != nil {
				t.Fatalf("Apply() unexpected error: %v", err)
			}

			resultDoc := result.Doc.(map[string]any)
			if resultDoc["value"] != tt.expected {
				t.Errorf("Apply() value = %v, want %v", resultDoc["value"], tt.expected)
			}
		})
	}
}
