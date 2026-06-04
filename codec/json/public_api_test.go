package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestDecodeTestTypeStringArray(t *testing.T) {
	t.Parallel()

	ops, err := Decode([]map[string]any{{
		"op":   "test_type",
		"path": "/value",
		"type": []string{"string", "number"},
	}}, PatchOptions{})
	require.NoError(t, err)
	require.Len(t, ops, 1)

	encoded, err := Encode(ops)
	require.NoError(t, err)
	assert.Equal(t, "test_type", encoded[0].Op)
	assert.Equal(t, "/value", encoded[0].Path)
	assert.Equal(t, []string{"string", "number"}, encoded[0].Type)
	assert.Nil(t, encoded[0].Value)
}

func TestDecodeCompositeNotAndApplyValidation(t *testing.T) {
	t.Parallel()

	ops, err := Decode([]map[string]any{{
		"op":   "not",
		"path": "/user",
		"apply": []any{
			map[string]any{"op": "defined", "path": "/name"},
		},
	}}, PatchOptions{})
	require.NoError(t, err)
	require.Len(t, ops, 1)

	encoded, err := Encode(ops)
	require.NoError(t, err)
	assert.Equal(t, "not", encoded[0].Op)
	require.Len(t, encoded[0].Apply, 1)
	assert.Equal(t, "/user/name", encoded[0].Apply[0].Path)

	_, err = Decode([]map[string]any{{
		"op":    "not",
		"path":  "/user",
		"apply": true,
	}}, PatchOptions{})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotOpMissingApply)
}

func TestDecodeJSONRejectsInvalidPayload(t *testing.T) {
	t.Parallel()

	ops, err := DecodeJSON([]byte(`{"op":"add"}`), PatchOptions{})
	require.Error(t, err)
	assert.Nil(t, ops)
}

func TestEncodeRejectsOperationWithoutJSONProjection(t *testing.T) {
	t.Parallel()

	encoded, err := Encode([]internal.Op{applyOnlyOp{}})
	require.Error(t, err)
	assert.Nil(t, encoded)
	assert.Contains(t, err.Error(), "cannot encode to JSON")
}

type applyOnlyOp struct{}

func (applyOnlyOp) Op() internal.OpType {
	return internal.OpAddType
}

func (applyOnlyOp) Path() []string {
	return []string{"name"}
}

func (applyOnlyOp) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}

func (applyOnlyOp) Validate() error {
	return nil
}
