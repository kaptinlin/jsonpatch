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

func TestApplyPreservesAdditionalDocumentShapes(t *testing.T) {
	t.Parallel()

	t.Run("named string remains named type", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, namedJSONString("Ada"), []jsoncodec.Operation{{Op: "replace", Path: "", Value: "Grace"}})
		require.NoError(t, err)
		assert.Equal(t, namedJSONString("Grace"), result.Doc)
		require.Len(t, result.Steps, 1)
		assert.Equal(t, "", result.Steps[0].Path())
	})

	t.Run("plain string root replacement stays string", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, "Ada", []jsoncodec.Operation{{Op: "replace", Path: "", Value: "Grace"}})
		require.NoError(t, err)
		assert.Equal(t, "Grace", result.Doc)
		require.Len(t, result.Steps, 1)
		assert.Equal(t, "Ada", result.Steps[0].Old())
	})

	t.Run("explicit JSON text root null result stays JSON text", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, jsonpatch.JSONText(`{"name":"Ada"}`), []jsoncodec.Operation{{Op: "replace", Path: "", Value: nil}})
		require.NoError(t, err)
		assert.Equal(t, jsonpatch.JSONText("null"), result.Doc)
		require.Len(t, result.Steps, 1)
		assert.Equal(t, map[string]any{"name": "Ada"}, result.Steps[0].Old())
	})

	t.Run("primitive root replacement stays primitive", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, 1, []jsoncodec.Operation{{Op: "replace", Path: "", Value: 2}})
		require.NoError(t, err)
		assert.Equal(t, 2, result.Doc)
		require.Len(t, result.Steps, 1)
		assert.Equal(t, 1, result.Steps[0].Old())
	})

	t.Run("nil document becomes structured value", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, any(nil), []jsoncodec.Operation{{Op: "add", Path: "", Value: map[string]any{"name": "Ada"}}})
		require.NoError(t, err)
		want := map[string]any{"name": "Ada"}
		if diff := cmp.Diff(want, result.Doc); diff != "" {
			t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
		}
		require.Len(t, result.Steps, 1)
		assert.True(t, result.Steps[0].Applied())
	})

	t.Run("interface slice is patched as primitive shape", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, any([]any{"Ada"}), []jsoncodec.Operation{{Op: "add", Path: "/-", Value: "Grace"}})
		require.NoError(t, err)
		want := []any{"Ada", "Grace"}
		if diff := cmp.Diff(want, result.Doc); diff != "" {
			t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
		}
		require.Len(t, result.Steps, 1)
		assert.Equal(t, "/-", result.Steps[0].Path())
	})
}

func TestCompileOpsPreservesAdditionalDocumentShapes(t *testing.T) {
	t.Parallel()

	t.Run("JSON bytes remain bytes", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps(t, []byte(`{"name":"Ada"}`), op.NewAdd([]string{"role"}, "admin"))
		require.NoError(t, err)
		assertJSONEqual(t, []byte(`{"name":"Ada","role":"admin"}`), result.Doc)
		require.Len(t, result.Steps, 1)
		assert.Equal(t, "/role", result.Steps[0].Path())
	})

	t.Run("explicit JSON text remains JSON text", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps(t, jsonpatch.JSONText(`{"name":"Ada"}`), op.NewAdd([]string{"role"}, "admin"))
		require.NoError(t, err)
		assertJSONEqual(t, `{"name":"Ada","role":"admin"}`, string(result.Doc))
		require.Len(t, result.Steps, 1)
	})

	t.Run("primitive root replacement stays primitive", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps(t, true, op.NewReplace(nil, false))
		require.NoError(t, err)
		assert.False(t, result.Doc)
		require.Len(t, result.Steps, 1)
		assert.Equal(t, true, result.Steps[0].Old())
	})

	t.Run("struct remains struct", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps(t, account{Name: "Ada", Age: 36}, op.NewReplace([]string{"name"}, "Grace"))
		require.NoError(t, err)
		want := account{Name: "Grace", Age: 36}
		if diff := cmp.Diff(want, result.Doc); diff != "" {
			t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
		}
		require.Len(t, result.Steps, 1)
		assert.Equal(t, "Ada", result.Steps[0].Old())
	})

	t.Run("interface slice remains slice value", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps(t, any([]any{"Ada"}), op.NewAdd([]string{"-"}, "Grace"))
		require.NoError(t, err)
		want := []any{"Ada", "Grace"}
		if diff := cmp.Diff(want, result.Doc); diff != "" {
			t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
		}
		require.Len(t, result.Steps, 1)
	})
}

func TestApplyReportsConversionFailures(t *testing.T) {
	t.Parallel()

	t.Run("operations cannot convert null result to concrete map", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, map[string]any{"name": "Ada"}, []jsoncodec.Operation{{Op: "replace", Path: "", Value: nil}})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
	})

	t.Run("ops cannot convert value result to concrete map", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps(t, map[string]any{"name": "Ada"}, op.NewReplace(nil, "Ada"))
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
	})

	t.Run("primitive conversion failure is classified", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps(t, 1, op.NewReplace(nil, "one"))
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
	})

	t.Run("plain string null conversion failure is classified", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, "Ada", []jsoncodec.Operation{{Op: "replace", Path: "", Value: nil}})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
	})
}

func TestApplyWrapsDecodeAndRuntimeErrors(t *testing.T) {
	t.Parallel()

	t.Run("decode error keeps codec sentinel", func(t *testing.T) {
		t.Parallel()

		result, err := applyOperations(t, map[string]any{"name": "Ada"}, []jsoncodec.Operation{{Op: "unknown", Path: "/name"}})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsoncodec.ErrCodecOpUnknown)
	})

	t.Run("runtime operation error keeps operation sentinel", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps(t, map[string]any{"name": "Ada"}, op.NewRemove([]string{"missing"}))
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, op.ErrPathNotFound)
	})
}

func applyOperations[T jsonpatch.Document](t *testing.T, doc T, operations []jsoncodec.Operation) (*jsonpatch.Result[T], error) {
	t.Helper()
	patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		return nil, err
	}
	return jsonpatch.Apply(patch, doc)
}

func applyOperationsInPlace[T jsonpatch.Document](t *testing.T, doc T, operations []jsoncodec.Operation) (*jsonpatch.Result[T], error) {
	t.Helper()
	patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		return nil, err
	}
	if err := jsonpatch.ApplyInPlace(patch, &doc); err != nil {
		return nil, err
	}
	return &jsonpatch.Result[T]{Doc: doc}, nil
}

func applyOps[T jsonpatch.Document](t *testing.T, doc T, operations ...jsonpatch.Op) (*jsonpatch.Result[T], error) {
	t.Helper()
	patch, err := jsonpatch.CompileOps(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		return nil, err
	}
	return jsonpatch.Apply(patch, doc)
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
