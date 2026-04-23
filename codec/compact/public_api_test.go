package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestPublicEncoderDecoderAPI(t *testing.T) {
	t.Parallel()

	encoder := NewEncoder(WithStringOpcode(true))
	decoder := NewDecoder()

	encoded, err := encoder.Encode(op.NewAdd([]string{"items", "0"}, "x"))
	require.NoError(t, err)
	assert.Equal(t, Op{"add", "/items/0", "x"}, encoded)

	decoded, err := decoder.Decode(encoded)
	require.NoError(t, err)
	assert.Equal(t, internal.OpAddType, decoded.Op())
	assert.Equal(t, []string{"items", "0"}, decoded.Path())
}

func TestPublicSliceAndJSONHelpers(t *testing.T) {
	t.Parallel()

	ops := []internal.Op{
		op.NewAdd([]string{"name"}, "Ada"),
		op.NewCopy([]string{"alias"}, []string{"name"}),
	}

	encoder := NewEncoder()
	encoded, err := encoder.EncodeSlice(ops)
	require.NoError(t, err)
	require.Len(t, encoded, 2)

	decoder := NewDecoder()
	decoded, err := decoder.DecodeSlice(encoded)
	require.NoError(t, err)
	require.Len(t, decoded, 2)
	assert.Equal(t, ops[0].Op(), decoded[0].Op())
	assert.Equal(t, ops[1].Op(), decoded[1].Op())

	jsonData, err := EncodeJSON(ops, WithStringOpcode(true))
	require.NoError(t, err)
	decodedFromJSON, err := DecodeJSON(jsonData)
	require.NoError(t, err)
	require.Len(t, decodedFromJSON, 2)
	assert.Equal(t, ops[0].Op(), decodedFromJSON[0].Op())
	assert.Equal(t, ops[1].Op(), decodedFromJSON[1].Op())
}

func TestDecodeErrorsAndOptions(t *testing.T) {
	t.Parallel()

	options := &Options{}
	WithStringOpcode(true)(options)
	assert.True(t, options.StringOpcode)

	t.Run("test_string defaults missing position to zero", func(t *testing.T) {
		t.Parallel()

		decoded, err := Decode([]Op{{CodeTestString, "/name", "Ad"}})
		require.NoError(t, err)
		require.Len(t, decoded, 1)

		testString, ok := decoded[0].(*op.TestStringOperation)
		require.True(t, ok)
		assert.Equal(t, 0, testString.Pos)
		assert.Equal(t, "Ad", testString.Str)
	})

	_, err := Decode([]Op{{CodeAdd}})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrMinLength)

	_, err = Decode([]Op{{CodeAdd, 1, "x"}})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPathNotString)
}
