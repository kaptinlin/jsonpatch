package json

import (
	"testing"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/kaptinlin/jsonpointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/codec/binary"
	"github.com/kaptinlin/jsonpatch/codec/compact"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestCompositeJSONPathsRemainStableAcrossRoundTrips(t *testing.T) {
	t.Parallel()

	wantWire := compositePathFixture()
	wantPaths := []string{
		"/profile",
		"/profile",
		"/profile/flags",
		"/profile/flags/status",
		"/profile/flags/status",
		"/profile/flags/name",
	}
	wire := wantWire

	for round := 1; round <= 3; round++ {
		decoded, err := Decode(wire, PatchOptions{})
		require.NoError(t, err, "round %d decode", round)
		assert.Equal(t, wantPaths, collectPaths(decoded), "round %d absolute paths", round)

		data, err := EncodeJSON(decoded)
		require.NoError(t, err, "round %d encode", round)
		wire = decodeOperationMaps(t, data)
		assert.Equal(t, wantWire, wire, "round %d wire paths", round)
	}
}

func TestCompositeJSONPathProjection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		parent     []string
		child      []string
		wantParent string
		wantChild  string
	}{
		{name: "root parent", child: []string{"profile"}, wantParent: "", wantChild: "/profile"},
		{name: "child equal to parent", parent: []string{"profile"}, child: []string{"profile"}, wantParent: "/profile", wantChild: ""},
		{name: "nested child", parent: []string{"profile"}, child: []string{"profile", "name"}, wantParent: "/profile", wantChild: "/name"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encoded, err := Encode([]internal.Op{op.NewAnd(tc.parent, []any{op.NewDefined(tc.child)})})
			require.NoError(t, err)
			require.Len(t, encoded, 1)
			require.Len(t, encoded[0].Apply, 1)
			assert.Equal(t, tc.wantParent, encoded[0].Path)
			assert.Equal(t, tc.wantChild, encoded[0].Apply[0].Path)
		})
	}
}

func TestCompositeJSONPathProjectionRejectsChildOutsideParent(t *testing.T) {
	t.Parallel()

	encoded, err := Encode([]internal.Op{
		op.NewAnd([]string{"profile"}, []any{op.NewDefined([]string{"account", "name"})}),
	})
	require.Error(t, err)
	assert.Nil(t, encoded)
	assert.ErrorIs(t, err, op.ErrPredicatePathOutsideParent)
}

func TestCompositeCodecsDecodeToSameAbsolutePaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		wire []map[string]any
	}{
		{name: "nested tree", wire: compositePathFixture()},
		{
			name: "relative child repeats parent segment",
			wire: []map[string]any{{
				"op":   "and",
				"path": "/profile",
				"apply": []any{
					map[string]any{"op": "defined", "path": "/profile"},
				},
			}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jsonOps, err := Decode(tc.wire, PatchOptions{})
			require.NoError(t, err)
			want := collectPaths(jsonOps)

			compactWire, err := compact.Encode(jsonOps)
			require.NoError(t, err)
			compactOps, err := compact.Decode(compactWire)
			require.NoError(t, err)
			assert.Equal(t, want, collectPaths(compactOps))

			binaryCodec := binary.New()
			binaryWire, err := binaryCodec.Encode(jsonOps)
			require.NoError(t, err)
			binaryOps, err := binaryCodec.Decode(binaryWire)
			require.NoError(t, err)
			assert.Equal(t, want, collectPaths(binaryOps))
		})
	}
}

func compositePathFixture() []map[string]any {
	return []map[string]any{{
		"op":   "and",
		"path": "/profile",
		"apply": []any{
			map[string]any{"op": "defined", "path": ""},
			map[string]any{
				"op":   "or",
				"path": "/flags",
				"apply": []any{
					map[string]any{
						"op":   "not",
						"path": "/status",
						"apply": []any{
							map[string]any{"op": "defined", "path": ""},
						},
					},
					map[string]any{"op": "defined", "path": "/name"},
				},
			},
		},
	}}
}

func decodeOperationMaps(t *testing.T, data []byte) []map[string]any {
	t.Helper()

	var operations []map[string]any
	require.NoError(t, jsonv2.Unmarshal(data, &operations))
	return operations
}

func collectPaths(operations []internal.Op) []string {
	var paths []string
	for _, operation := range operations {
		paths = appendOperationPaths(paths, operation)
	}
	return paths
}

func appendOperationPaths(paths []string, operation internal.Op) []string {
	paths = append(paths, jsonpointer.Format(operation.Path()...))
	if composite, ok := operation.(internal.SecondOrderPredicateOp); ok {
		for _, child := range composite.Ops() {
			paths = appendOperationPaths(paths, child)
		}
	}
	return paths
}
