package jsonpatch_test

import (
	"testing"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/codec/binary"
	"github.com/kaptinlin/jsonpatch/codec/compact"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestRemoveOldValuePresenceRoundTripsAcrossCodecs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		oldValue any
		hasOld   bool
	}{
		{name: "absent"},
		{name: "string", oldValue: "Ada", hasOld: true},
		{name: "object", oldValue: map[string]any{"name": "Ada"}, hasOld: true},
		{name: "explicit null", oldValue: nil, hasOld: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var operation *op.RemoveOperation
			if tt.hasOld {
				operation = op.NewRemoveWithOldValue([]string{"name"}, tt.oldValue)
			} else {
				operation = op.NewRemove([]string{"name"})
			}

			t.Run("JSON", func(t *testing.T) {
				t.Parallel()

				encoded, err := jsoncodec.EncodeJSON([]internal.Op{operation})
				require.NoError(t, err)
				var wire []map[string]any
				require.NoError(t, jsonv2.Unmarshal(encoded, &wire))
				oldValue, hasOld := wire[0]["oldValue"]
				assert.Equal(t, tt.hasOld, hasOld)
				if diff := cmp.Diff(tt.oldValue, oldValue); tt.hasOld && diff != "" {
					t.Errorf("oldValue mismatch (-want +got):\n%s", diff)
				}

				decoded, err := jsoncodec.DecodeJSON(encoded, jsoncodec.PatchOptions{})
				require.NoError(t, err)
				assertRemoveOldValue(t, decoded[0], tt.hasOld, tt.oldValue)
			})

			for name, options := range map[string][]compact.Option{
				"numeric compact": nil,
				"string compact":  {compact.WithStringOpcode(true)},
			} {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					encoded, err := compact.Encode([]internal.Op{operation}, options...)
					require.NoError(t, err)
					decoded, err := compact.Decode(encoded)
					require.NoError(t, err)
					assertRemoveOldValue(t, decoded[0], tt.hasOld, tt.oldValue)
				})
			}

			t.Run("binary", func(t *testing.T) {
				t.Parallel()

				codec := binary.New()
				encoded, err := codec.Encode([]internal.Op{operation})
				require.NoError(t, err)
				decoded, err := codec.Decode(encoded)
				require.NoError(t, err)
				assertRemoveOldValue(t, decoded[0], tt.hasOld, tt.oldValue)
			})
		})
	}
}

func TestReplaceOldValuePresenceRoundTripsAcrossCodecs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		oldValue any
		hasOld   bool
	}{
		{name: "absent"},
		{name: "string", oldValue: "Ada", hasOld: true},
		{name: "object", oldValue: map[string]any{"name": "Ada"}, hasOld: true},
		{name: "explicit null", oldValue: nil, hasOld: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var operation *op.ReplaceOperation
			if tt.hasOld {
				operation = op.NewReplaceWithOldValue([]string{"name"}, "Grace", tt.oldValue)
			} else {
				operation = op.NewReplace([]string{"name"}, "Grace")
			}

			t.Run("JSON", func(t *testing.T) {
				t.Parallel()

				encoded, err := jsoncodec.EncodeJSON([]internal.Op{operation})
				require.NoError(t, err)
				assertJSONOldValue(t, encoded, tt.hasOld, tt.oldValue)
				decoded, err := jsoncodec.DecodeJSON(encoded, jsoncodec.PatchOptions{})
				require.NoError(t, err)
				assertOperationOldValue(t, decoded[0], tt.hasOld, tt.oldValue)
			})

			for name, options := range map[string][]compact.Option{
				"numeric compact": nil,
				"string compact":  {compact.WithStringOpcode(true)},
			} {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					encoded, err := compact.Encode([]internal.Op{operation}, options...)
					require.NoError(t, err)
					decoded, err := compact.Decode(encoded)
					require.NoError(t, err)
					assertOperationOldValue(t, decoded[0], tt.hasOld, tt.oldValue)
				})
			}

			t.Run("binary", func(t *testing.T) {
				t.Parallel()

				codec := binary.New()
				encoded, err := codec.Encode([]internal.Op{operation})
				require.NoError(t, err)
				decoded, err := codec.Decode(encoded)
				require.NoError(t, err)
				assertOperationOldValue(t, decoded[0], tt.hasOld, tt.oldValue)
			})
		})
	}
}

func assertRemoveOldValue(t *testing.T, operation internal.Op, wantHas bool, wantValue any) {
	t.Helper()

	remove, ok := operation.(*op.RemoveOperation)
	require.True(t, ok)
	assert.Equal(t, wantHas, remove.HasOldValue)
	if diff := cmp.Diff(wantValue, remove.OldValue); wantHas && diff != "" {
		t.Errorf("oldValue mismatch (-want +got):\n%s", diff)
	}
}

func assertOperationOldValue(t *testing.T, operation internal.Op, wantHas bool, wantValue any) {
	t.Helper()

	encoded, err := jsoncodec.EncodeJSON([]internal.Op{operation})
	require.NoError(t, err)
	assertJSONOldValue(t, encoded, wantHas, wantValue)
}

func assertJSONOldValue(t *testing.T, encoded []byte, wantHas bool, wantValue any) {
	t.Helper()

	var wire []map[string]any
	require.NoError(t, jsonv2.Unmarshal(encoded, &wire))
	oldValue, hasOld := wire[0]["oldValue"]
	assert.Equal(t, wantHas, hasOld)
	if diff := cmp.Diff(wantValue, oldValue); wantHas && diff != "" {
		t.Errorf("oldValue mismatch (-want +got):\n%s", diff)
	}
}
