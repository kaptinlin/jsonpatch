package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpInc_Apply(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		doc      interface{}
		inc      float64
		expected interface{}
		oldValue interface{}
		wantErr  bool
	}{
		{
			name:     "inc int field",
			path:     []string{"count"},
			doc:      map[string]interface{}{"count": 1},
			inc:      2,
			expected: map[string]interface{}{"count": 3.0},
			oldValue: 1.0,
		},
		{
			name:     "inc float field",
			path:     []string{"score"},
			doc:      map[string]interface{}{"score": 1.5},
			inc:      0.5,
			expected: map[string]interface{}{"score": 2.0},
			oldValue: 1.5,
		},
		{
			name:     "dec int field",
			path:     []string{"count"},
			doc:      map[string]interface{}{"count": 5},
			inc:      -3,
			expected: map[string]interface{}{"count": 2.0},
			oldValue: 5.0,
		},
		{
			name:     "inc nested field",
			path:     []string{"user", "age"},
			doc:      map[string]interface{}{"user": map[string]interface{}{"age": 20}},
			inc:      1,
			expected: map[string]interface{}{"user": map[string]interface{}{"age": 21.0}},
			oldValue: 20.0,
		},
		{
			name:     "inc array element",
			path:     []string{"nums", "1"},
			doc:      map[string]interface{}{"nums": []interface{}{1, 2, 3}},
			inc:      10,
			expected: map[string]interface{}{"nums": []interface{}{1, 12.0, 3}},
			oldValue: 2.0,
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
			name:    "path not found",
			path:    []string{"notfound"},
			doc:     map[string]interface{}{"count": 1},
			inc:     1,
			wantErr: true,
		},
		{
			name:    "not a number",
			path:    []string{"str"},
			doc:     map[string]interface{}{"str": "abc"},
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
			op := NewInc(tt.path, tt.inc)
			docCopy, err := DeepClone(tt.doc)
			require.NoError(t, err)

			result, err := op.Apply(docCopy)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Doc)
			assert.Equal(t, tt.oldValue, result.Old)
		})
	}
}

func TestOpInc_Op(t *testing.T) {
	op := NewInc([]string{"count"}, 1)
	assert.Equal(t, internal.OpIncType, op.Op())
}

func TestOpInc_Code(t *testing.T) {
	op := NewInc([]string{"count"}, 1)
	assert.Equal(t, internal.OpIncCode, op.Code())
}

func TestOpInc_NewOpInc(t *testing.T) {
	path := []string{"user", "score"}
	inc := 3.5
	op := NewInc(path, inc)
	assert.Equal(t, path, op.path)
	assert.Equal(t, inc, op.Inc)
	assert.Equal(t, internal.OpIncType, op.Op())
	assert.Equal(t, internal.OpIncCode, op.Code())
}

func TestOpInc_ToJSON(t *testing.T) {
	op := NewInc([]string{"count"}, 5.5)
	jsonOp, err := op.ToJSON()
	require.NoError(t, err)

	assert.Equal(t, "inc", jsonOp.Op)
	assert.Equal(t, "/count", jsonOp.Path)
	assert.Equal(t, 5.5, jsonOp.Inc)
}
