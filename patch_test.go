package jsonpatch_test

import (
	"errors"
	"regexp/syntax"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestCompileJSONCapabilityPolicy(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.CompileJSON([]byte(`[{"op":"inc","path":"/count","inc":1}]`))
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, jsonpatch.ErrUnsupportedCapability)

	var patchErr *jsonpatch.Error
	require.True(t, errors.As(err, &patchErr))
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "inc", patchErr.Op())
	assert.Equal(t, "/count", patchErr.Path())
	assert.Equal(t, "json", patchErr.Codec())

	patch, err = jsonpatch.CompileJSON(
		[]byte(`[{"op":"inc","path":"/count","inc":1}]`),
		jsonpatch.WithCapabilities(jsonpatch.RFC6902, jsonpatch.Extended),
	)
	require.NoError(t, err)

	result, err := jsonpatch.Apply(patch, map[string]any{"count": float64(1)})
	require.NoError(t, err)
	assert.Equal(t, float64(2), result.Doc["count"])
}

func TestCompileJSONPayloadErrorIncludesOperationContext(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.CompileJSON([]byte(`[{"op":"add","path":"/name"}]`))
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, jsonpatch.ErrPayloadInvalid)

	var patchErr *jsonpatch.Error
	require.True(t, errors.As(err, &patchErr))
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "add", patchErr.Op())
	assert.Equal(t, "/name", patchErr.Path())
	assert.Equal(t, "json", patchErr.Codec())
}

func TestCompileRejectsInvalidRegex(t *testing.T) {
	t.Parallel()

	t.Run("JSON operation", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.CompileJSON(
			[]byte(`[{"op":"matches","path":"/name","value":"["}]`),
			jsonpatch.WithCapabilities(jsonpatch.RegexPredicate),
		)

		require.ErrorIs(t, err, jsonpatch.ErrPayloadInvalid)
		assert.Nil(t, patch)
		var syntaxErr *syntax.Error
		require.ErrorAs(t, err, &syntaxErr)
	})

	t.Run("Go operation", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.CompileOps(
			[]jsonpatch.Op{op.NewMatches([]string{"name"}, "[", false, nil)},
			jsonpatch.WithCapabilities(jsonpatch.RegexPredicate),
		)

		require.ErrorIs(t, err, jsonpatch.ErrPayloadInvalid)
		assert.Nil(t, patch)
		var syntaxErr *syntax.Error
		require.ErrorAs(t, err, &syntaxErr)
	})

	t.Run("custom matcher factory", func(t *testing.T) {
		t.Parallel()

		factoryErr := errors.New("custom matcher rejected pattern")
		operation := op.NewMatches([]string{"name"}, "custom", false,
			func(string, bool) (internal.RegexMatcher, error) {
				return nil, factoryErr
			},
		)
		patch, err := jsonpatch.CompileOps(
			[]jsonpatch.Op{operation},
			jsonpatch.WithCapabilities(jsonpatch.RegexPredicate),
		)

		require.ErrorIs(t, err, jsonpatch.ErrPayloadInvalid)
		require.ErrorIs(t, err, factoryErr)
		assert.Nil(t, patch)
	})
}

func TestPatchApplyReportsStepFactsAndPreservesInput(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(
		op.NewReplace([]string{"name"}, "Grace"),
		op.NewCopy([]string{"displayName"}, []string{"name"}),
	)
	require.NoError(t, err)

	doc := map[string]any{"name": "Ada"}
	result, err := jsonpatch.Apply(patch, doc)
	require.NoError(t, err)

	assert.Equal(t, "Ada", doc["name"])
	assert.Equal(t, "Grace", result.Doc["name"])
	assert.Equal(t, "Grace", result.Doc["displayName"])
	require.Len(t, result.Steps, 2)

	assert.Equal(t, 0, result.Steps[0].Index())
	assert.Equal(t, "replace", result.Steps[0].Op())
	assert.Equal(t, "/name", result.Steps[0].Path())
	assert.Empty(t, result.Steps[0].From())
	assert.Equal(t, "Ada", result.Steps[0].Old())

	assert.Equal(t, 1, result.Steps[1].Index())
	assert.Equal(t, "copy", result.Steps[1].Op())
	assert.Equal(t, "/displayName", result.Steps[1].Path())
	assert.Equal(t, "/name", result.Steps[1].From())
	assert.Nil(t, result.Steps[1].Old())
}

func TestCompileClonesOperationPayloads(t *testing.T) {
	t.Parallel()

	payload := map[string]any{"name": "Ada"}
	operation := op.NewAdd([]string{"user"}, payload)

	patch, err := jsonpatch.Compile(operation)
	require.NoError(t, err)

	payload["name"] = "mutated"
	operation.Value = map[string]any{"name": "also-mutated"}

	result, err := jsonpatch.Apply(patch, map[string]any{})
	require.NoError(t, err)

	want := map[string]any{
		"user": map[string]any{"name": "Ada"},
	}
	if diff := cmp.Diff(want, result.Doc); diff != "" {
		t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
	}
}

func TestOperationContractIncludesClone(t *testing.T) {
	t.Parallel()

	path := []string{"name"}
	var operation jsonpatch.Op = op.NewAdd(path, "Ada")
	cloned, err := operation.Clone()
	require.NoError(t, err)

	path[0] = "changed"
	assert.Equal(t, []string{"name"}, cloned.Path())
}

func TestCompileClonesCompositePredicateOperands(t *testing.T) {
	t.Parallel()

	operand := op.NewContains([]string{"name"}, "Ada")
	operands := []any{operand}
	operation := op.NewAnd(nil, operands)

	patch, err := jsonpatch.CompileOps([]jsonpatch.Op{operation}, jsonpatch.WithCapabilities(jsonpatch.Predicate))
	require.NoError(t, err)

	operand.Value = "Grace"
	operands[0] = op.NewContains([]string{"name"}, "Grace")
	operation.Operations = nil

	_, err = jsonpatch.Apply(patch, map[string]any{"name": "Ada Lovelace"})
	require.NoError(t, err)
}

func TestPatchSamePathCopyUsesAddSemantics(t *testing.T) {
	t.Parallel()

	t.Run("object member", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.Compile(op.NewCopy([]string{"name"}, []string{"name"}))
		require.NoError(t, err)

		doc := map[string]any{"name": "Ada"}
		result, err := jsonpatch.Apply(patch, doc)
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"name": "Ada"}, result.Doc)
		assert.Equal(t, map[string]any{"name": "Ada"}, doc)
	})

	t.Run("array element", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.Compile(op.NewCopy([]string{"0"}, []string{"0"}))
		require.NoError(t, err)

		doc := []any{map[string]any{"name": "Ada"}}
		result, err := jsonpatch.Apply(patch, doc)
		require.NoError(t, err)
		require.Len(t, result.Doc, 2)

		copied := result.Doc[0].(map[string]any)
		source := result.Doc[1].(map[string]any)
		copied["name"] = "Grace"
		assert.Equal(t, "Ada", source["name"])
		assert.Equal(t, []any{map[string]any{"name": "Ada"}}, doc)
	})

	t.Run("root", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.Compile(op.NewCopy(nil, nil))
		require.NoError(t, err)

		doc := map[string]any{"name": "Ada"}
		result, err := jsonpatch.Apply(patch, doc)
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"name": "Ada"}, result.Doc)
		assert.Equal(t, map[string]any{"name": "Ada"}, doc)
	})
}

func TestCompilePreservesCustomRegexMatcher(t *testing.T) {
	t.Parallel()

	matcherUsed := false
	operation := op.NewMatches([]string{"name"}, "not-a-real-pattern-contract", false,
		func(_ string, _ bool) (internal.RegexMatcher, error) {
			return func(value string) bool {
				matcherUsed = true
				return value == "custom"
			}, nil
		},
	)

	patch, err := jsonpatch.CompileOps([]jsonpatch.Op{operation}, jsonpatch.WithCapabilities(jsonpatch.RegexPredicate))
	require.NoError(t, err)

	operation.Pattern = "mutated-after-compile"
	result, err := jsonpatch.Apply(patch, map[string]any{"name": "custom"})
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "custom"}, result.Doc)
	assert.True(t, matcherUsed)
}

func TestPatchApplyInPlaceIsExplicit(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(op.NewReplace([]string{"name"}, "Grace"))
	require.NoError(t, err)

	doc := map[string]any{"name": "Ada"}
	err = jsonpatch.ApplyInPlace(patch, &doc)
	require.NoError(t, err)
	assert.Equal(t, "Grace", doc["name"])
}

func TestPatchApplyDistinguishesJSONTextFromScalarString(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(op.NewAdd([]string{"role"}, "admin"))
	require.NoError(t, err)

	_, err = jsonpatch.Apply(patch, `{"name":"Ada"}`)
	require.Error(t, err)
	assert.ErrorIs(t, err, jsonpatch.ErrRuntimeConflict)

	result, err := jsonpatch.Apply(patch, jsonpatch.JSONText(`{"name":"Ada"}`))
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"Ada","role":"admin"}`, string(result.Doc))

	_, err = jsonpatch.Apply(patch, jsonpatch.JSONText(`{"name":`))
	require.Error(t, err)
	assert.ErrorIs(t, err, jsonpatch.ErrPayloadInvalid)
}

func TestPatchApplyClassifiesStringAliasAsScalar(t *testing.T) {
	t.Parallel()

	type scalarString string

	patch, err := jsonpatch.Compile(op.NewAdd([]string{"role"}, "admin"))
	require.NoError(t, err)

	_, err = jsonpatch.Apply(patch, scalarString(`{"name":"Ada"}`))
	require.Error(t, err)
	assert.ErrorIs(t, err, jsonpatch.ErrRuntimeConflict)
}

func TestPatchApplyClassifiesByteAliasAsJSONText(t *testing.T) {
	t.Parallel()

	type jsonBytes []byte

	patch, err := jsonpatch.Compile(op.NewAdd([]string{"role"}, "admin"))
	require.NoError(t, err)

	result, err := jsonpatch.Apply(patch, jsonBytes(`{"name":"Ada"}`))
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"Ada","role":"admin"}`, string(result.Doc))
}

func TestPatchApplyStructuredErrors(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(op.NewTest([]string{"name"}, "Grace"))
	require.NoError(t, err)

	result, err := jsonpatch.Apply(patch, map[string]any{"name": "Ada"})
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, jsonpatch.ErrTestFailed)

	var patchErr *jsonpatch.Error
	require.True(t, errors.As(err, &patchErr))
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "test", patchErr.Op())
	assert.Equal(t, "/name", patchErr.Path())
	assert.Empty(t, patchErr.Codec())
	assert.Empty(t, cmp.Diff(jsonpatch.ErrTestFailed, patchErr.Kind(), cmp.Comparer(errors.Is)))
}

func FuzzApplyPositionedOperationsDoesNotPanic(f *testing.F) {
	f.Add("abc", "a", 0, 1)
	f.Add("", "", -1, 1)
	f.Add("text", "value", 100, 100)

	f.Fuzz(func(_ *testing.T, text, needle string, testPos, mergePos int) {
		testPatch, err := jsonpatch.CompileOps(
			[]jsonpatch.Op{op.NewTestString(nil, needle, float64(testPos), false, false)},
			jsonpatch.WithCapabilities(jsonpatch.Predicate),
		)
		if err == nil {
			_, _ = jsonpatch.Apply(testPatch, text)
		}

		mergePatch, err := jsonpatch.CompileOps(
			[]jsonpatch.Op{op.NewMerge([]string{"1"}, float64(mergePos), nil)},
			jsonpatch.WithCapabilities(jsonpatch.Extended),
		)
		if err == nil {
			_, _ = jsonpatch.Apply(mergePatch, []any{text, needle})
		}
	})
}
