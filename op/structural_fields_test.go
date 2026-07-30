package op

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitComputedFieldsOverrideProps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		pos   int
		props map[string]any
		want  any
	}{
		{
			name:  "plain string",
			value: "Ada",
			pos:   1,
			props: map[string]any{"text": "override", "mark": true},
			want: []any{
				map[string]any{"text": "A", "mark": true},
				map[string]any{"text": "da", "mark": true},
			},
		},
		{
			name:  "Slate text",
			value: map[string]any{"text": "Ada", "bold": true},
			pos:   1,
			props: map[string]any{"text": "override", "bold": false, "mark": true},
			want: []any{
				map[string]any{"text": "A", "bold": false, "mark": true},
				map[string]any{"text": "da", "bold": false, "mark": true},
			},
		},
		{
			name: "Slate element",
			value: map[string]any{
				"type":     "paragraph",
				"children": []any{"A", "B"},
			},
			pos:   1,
			props: map[string]any{"children": []any{"override"}, "type": "quote"},
			want: []any{
				map[string]any{"type": "quote", "children": []any{"A"}},
				map[string]any{"type": "quote", "children": []any{"B"}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := NewSplit(nil, tc.pos, tc.props).Apply(tc.value)
			require.NoError(t, err)
			assert.Equal(t, tc.want, result.Doc)
		})
	}
}

func TestMergePropsAreMetadataOnly(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		one  any
		two  any
		want any
	}{
		{name: "string", one: "A", two: "da", want: "Ada"},
		{name: "number", one: 1, two: 2, want: float64(3)},
		{
			name: "Slate text",
			one:  map[string]any{"text": "A", "bold": true},
			two:  map[string]any{"text": "da", "italic": true},
			want: map[string]any{"text": "Ada", "bold": true, "italic": true},
		},
		{
			name: "Slate element",
			one:  map[string]any{"type": "p", "children": []any{"A"}},
			two:  map[string]any{"align": "center", "children": []any{"B"}},
			want: map[string]any{"type": "p", "align": "center", "children": []any{"A", "B"}},
		},
		{name: "fallback", one: true, two: nil, want: []any{true, nil}},
	}

	props := map[string]any{
		"text":     "override",
		"children": []any{"override"},
		"injected": true,
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := NewMerge([]string{"1"}, 0, props).Apply([]any{tc.one, tc.two})
			require.NoError(t, err)
			assert.Equal(t, []any{tc.want}, result.Doc)
		})
	}
}

func TestSplitThenMergeRestoresSupportedValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		pos   int
	}{
		{name: "string", value: "Ada", pos: 1},
		{name: "number", value: 7.0, pos: 3},
		{name: "Slate text", value: map[string]any{"text": "Ada", "bold": true}, pos: 1},
		{name: "Slate element", value: map[string]any{"type": "p", "children": []any{"A", "B"}}, pos: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			split, err := NewSplit(nil, tc.pos, nil).Apply(tc.value)
			require.NoError(t, err)
			merged, err := NewMerge([]string{"1"}, tc.pos, nil).Apply(split.Doc)
			require.NoError(t, err)
			assert.Equal(t, []any{tc.value}, merged.Doc)
		})
	}
}
