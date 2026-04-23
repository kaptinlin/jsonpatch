package binary

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestCodec_RoundTripViaPublicAPI(t *testing.T) {
	t.Parallel()

	codec := New()
	require.NotNil(t, codec)

	ops := []internal.Op{
		op.NewAdd([]string{"profile", "name"}, "Ada"),
		op.NewTestString([]string{"profile", "name"}, "Ad", 0, false, false),
		op.NewMatches([]string{"profile", "name"}, "^Ada$", false, nil),
		op.NewTestStringLenWithNot([]string{"profile", "name"}, 3, false),
		op.NewMerge([]string{"nodes"}, 1, map[string]any{"joined": true}),
	}

	data, err := codec.Encode(ops)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	decoded, err := codec.Decode(data)
	require.NoError(t, err)
	require.Len(t, decoded, len(ops))

	for i := range ops {
		got, err := decoded[i].ToJSON()
		require.NoError(t, err)
		want, err := ops[i].ToJSON()
		require.NoError(t, err)
		assert.Equal(t, want, got)
	}
}

func TestCodec_DecodeRejectsInvalidData(t *testing.T) {
	t.Parallel()

	codec := New()
	decoded, err := codec.Decode([]byte{0xc1})
	require.Error(t, err)
	assert.Nil(t, decoded)
}
