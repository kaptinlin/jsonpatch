package jsonpatch_test

import (
	"errors"
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

func TestPatchApplyReportsStepFactsAndPreservesInput(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(op.NewReplace([]string{"name"}, "Grace"))
	require.NoError(t, err)

	doc := map[string]any{"name": "Ada"}
	result, err := jsonpatch.Apply(patch, doc)
	require.NoError(t, err)

	assert.Equal(t, "Ada", doc["name"])
	assert.Equal(t, "Grace", result.Doc["name"])
	require.Len(t, result.Steps, 1)

	step := result.Steps[0]
	assert.Equal(t, 0, step.Index())
	assert.Equal(t, "replace", step.Op())
	assert.Equal(t, "/name", step.Path())
	assert.Empty(t, step.From())
	assert.Equal(t, "Ada", step.Old())
	assert.True(t, step.Applied())
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

func TestCompilePreservesCustomRegexMatcher(t *testing.T) {
	t.Parallel()

	matcherUsed := false
	operation := op.NewMatches([]string{"name"}, "not-a-real-pattern-contract", false,
		func(_ string, _ bool) internal.RegexMatcher {
			return func(value string) bool {
				matcherUsed = true
				return value == "custom"
			}
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

	_, err = jsonpatch.Apply(patch, map[string]any{"name": "Ada"})
	require.Error(t, err)
	assert.ErrorIs(t, err, jsonpatch.ErrTestFailed)

	var patchErr *jsonpatch.Error
	require.True(t, errors.As(err, &patchErr))
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "test", patchErr.Op())
	assert.Equal(t, "/name", patchErr.Path())
	assert.Empty(t, patchErr.Codec())
	assert.Empty(t, cmp.Diff(jsonpatch.ErrTestFailed, patchErr.Kind(), cmp.Comparer(errors.Is)))
}

func TestCompileRejectsOperationWithoutClone(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(applyOnlyOp{})
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, jsonpatch.ErrPayloadInvalid)

	var patchErr *jsonpatch.Error
	require.True(t, errors.As(err, &patchErr))
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "add", patchErr.Op())
	assert.Equal(t, "/name", patchErr.Path())
	assert.Error(t, patchErr.Cause())
}

type applyOnlyOp struct{}

func (applyOnlyOp) Op() jsonpatch.OpType {
	return jsonpatch.OpAddType
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
