package jsonpatch_test

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/op"
)

type namedJSONString string

type account struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestApplyPatchPreservesAdditionalDocumentShapes(t *testing.T) {
	t.Parallel()

	t.Run("named string remains named type", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyPatch(namedJSONString("Ada"), []jsonpatch.Operation{{Op: "replace", Path: "", Value: "Grace"}})
		require.NoError(t, err)
		assert.Equal(t, namedJSONString("Grace"), result.Doc)
		require.Len(t, result.Res, 1)
		assert.Equal(t, namedJSONString("Grace"), result.Res[0].Doc)
	})

	t.Run("plain string root replacement stays string", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyPatch("Ada", []jsonpatch.Operation{{Op: "replace", Path: "", Value: "Grace"}})
		require.NoError(t, err)
		assert.Equal(t, "Grace", result.Doc)
		require.Len(t, result.Res, 1)
		assert.Equal(t, "Grace", result.Res[0].Doc)
	})

	t.Run("primitive root replacement stays primitive", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyPatch(1, []jsonpatch.Operation{{Op: "replace", Path: "", Value: 2}})
		require.NoError(t, err)
		assert.Equal(t, 2, result.Doc)
		require.Len(t, result.Res, 1)
		assert.Equal(t, 2, result.Res[0].Doc)
	})

	t.Run("nil document becomes structured value", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyPatch(any(nil), []jsonpatch.Operation{{Op: "add", Path: "", Value: map[string]any{"name": "Ada"}}})
		require.NoError(t, err)
		want := map[string]any{"name": "Ada"}
		if diff := cmp.Diff(want, result.Doc); diff != "" {
			t.Errorf("ApplyPatch() document mismatch (-want +got):\n%s", diff)
		}
		require.Len(t, result.Res, 1)
		if diff := cmp.Diff(want, result.Res[0].Doc); diff != "" {
			t.Errorf("ApplyPatch() operation result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("interface slice is patched as primitive shape", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyPatch(any([]any{"Ada"}), []jsonpatch.Operation{{Op: "add", Path: "/-", Value: "Grace"}})
		require.NoError(t, err)
		want := []any{"Ada", "Grace"}
		if diff := cmp.Diff(want, result.Doc); diff != "" {
			t.Errorf("ApplyPatch() document mismatch (-want +got):\n%s", diff)
		}
		require.Len(t, result.Res, 1)
		if diff := cmp.Diff(want, result.Res[0].Doc); diff != "" {
			t.Errorf("ApplyPatch() operation result mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestApplyOpsPreservesAdditionalDocumentShapes(t *testing.T) {
	t.Parallel()

	t.Run("JSON bytes remain bytes", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyOps([]byte(`{"name":"Ada"}`), []jsonpatch.Op{op.NewAdd([]string{"role"}, "admin")})
		require.NoError(t, err)
		assertJSONEqual(t, []byte(`{"name":"Ada","role":"admin"}`), result.Doc)
		require.Len(t, result.Res, 1)
		assertJSONEqual(t, []byte(`{"name":"Ada","role":"admin"}`), result.Res[0].Doc)
	})

	t.Run("JSON string remains string", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyOps(`{"name":"Ada"}`, []jsonpatch.Op{op.NewAdd([]string{"role"}, "admin")})
		require.NoError(t, err)
		assertJSONEqual(t, `{"name":"Ada","role":"admin"}`, result.Doc)
		require.Len(t, result.Res, 1)
		assertJSONEqual(t, `{"name":"Ada","role":"admin"}`, result.Res[0].Doc)
	})

	t.Run("primitive root replacement stays primitive", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyOps(true, []jsonpatch.Op{op.NewReplace(nil, false)})
		require.NoError(t, err)
		assert.False(t, result.Doc)
		require.Len(t, result.Res, 1)
		assert.False(t, result.Res[0].Doc)
	})

	t.Run("struct remains struct", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyOps(account{Name: "Ada", Age: 36}, []jsonpatch.Op{op.NewReplace([]string{"name"}, "Grace")})
		require.NoError(t, err)
		want := account{Name: "Grace", Age: 36}
		if diff := cmp.Diff(want, result.Doc); diff != "" {
			t.Errorf("ApplyOps() document mismatch (-want +got):\n%s", diff)
		}
		require.Len(t, result.Res, 1)
		if diff := cmp.Diff(want, result.Res[0].Doc); diff != "" {
			t.Errorf("ApplyOps() operation result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("interface slice remains slice value", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyOps(any([]any{"Ada"}), []jsonpatch.Op{op.NewAdd([]string{"-"}, "Grace")})
		require.NoError(t, err)
		want := []any{"Ada", "Grace"}
		if diff := cmp.Diff(want, result.Doc); diff != "" {
			t.Errorf("ApplyOps() document mismatch (-want +got):\n%s", diff)
		}
		require.Len(t, result.Res, 1)
		if diff := cmp.Diff(want, result.Res[0].Doc); diff != "" {
			t.Errorf("ApplyOps() operation result mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestApplyAPIsReportConversionFailures(t *testing.T) {
	t.Parallel()

	t.Run("patch cannot convert null result to concrete map", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyPatch(map[string]any{"name": "Ada"}, []jsonpatch.Operation{{Op: "replace", Path: "", Value: nil}})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
	})

	t.Run("ops cannot convert value result to concrete map", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyOps(map[string]any{"name": "Ada"}, []jsonpatch.Op{op.NewReplace(nil, "Ada")})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
	})

	t.Run("primitive conversion failure is classified", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyOps(1, []jsonpatch.Op{op.NewReplace(nil, "one")})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
	})
}

func TestApplyPatchWrapsDecodeAndRuntimeErrors(t *testing.T) {
	t.Parallel()

	t.Run("decode error keeps validation sentinel", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyPatch(map[string]any{"name": "Ada"}, []jsonpatch.Operation{{Op: "unknown", Path: "/name"}})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsoncodec.ErrCodecOpUnknown)
	})

	t.Run("runtime operation error keeps operation sentinel", func(t *testing.T) {
		t.Parallel()

		result, err := jsonpatch.ApplyOps(map[string]any{"name": "Ada"}, []jsonpatch.Op{op.NewRemove([]string{"missing"})})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, op.ErrPathNotFound)
	})
}

func assertJSONEqual(t *testing.T, want, got any) {
	t.Helper()

	wantValue := decodeJSONComparable(t, want)
	gotValue := decodeJSONComparable(t, got)
	if diff := cmp.Diff(wantValue, gotValue); diff != "" {
		t.Errorf("JSON value mismatch (-want +got):\n%s", diff)
	}
}

func decodeJSONComparable(t *testing.T, value any) any {
	t.Helper()

	switch v := value.(type) {
	case []byte:
		var decoded any
		require.NoError(t, json.Unmarshal(v, &decoded))
		return decoded
	case string:
		if len(v) == 0 || (v[0] != '{' && v[0] != '[') {
			return v
		}
		var decoded any
		require.NoError(t, json.Unmarshal([]byte(v), &decoded))
		return decoded
	default:
		return value
	}
}
