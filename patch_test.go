package jsonpatch_test

import (
	"errors"
	"maps"
	"math"
	"regexp/syntax"
	"slices"
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

func TestCompileRejectsNestedRegexWithoutCapability(t *testing.T) {
	t.Parallel()

	operation := op.NewAnd(nil,
		op.NewMatches([]string{"name"}, "^Ada$", false, nil),
	)
	patch, err := jsonpatch.CompileOps(
		[]jsonpatch.Op{operation},
		jsonpatch.WithCapabilities(jsonpatch.Predicate),
	)

	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, jsonpatch.ErrUnsupportedCapability)

	var patchErr *jsonpatch.Error
	require.ErrorAs(t, err, &patchErr)
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "matches", patchErr.Op())
	assert.Equal(t, "/name", patchErr.Path())
}

func TestCompileRejectsDeeplyNestedRegexWithoutCapability(t *testing.T) {
	t.Parallel()

	operation := op.NewOr(nil,
		op.NewAnd(nil,
			op.NewMatches([]string{"profile", "name"}, "^Ada$", false, nil),
		),
	)
	patch, err := jsonpatch.CompileOps(
		[]jsonpatch.Op{operation},
		jsonpatch.WithCapabilities(jsonpatch.Predicate),
	)

	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, jsonpatch.ErrUnsupportedCapability)

	var patchErr *jsonpatch.Error
	require.ErrorAs(t, err, &patchErr)
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "matches", patchErr.Op())
	assert.Equal(t, "/profile/name", patchErr.Path())
}

func TestCompileRejectsRegexInsideNotWithoutCapability(t *testing.T) {
	t.Parallel()

	operation := op.NewNot(nil,
		op.NewMatches([]string{"name"}, "^Ada$", false, nil),
	)
	patch, err := jsonpatch.CompileOps(
		[]jsonpatch.Op{operation},
		jsonpatch.WithCapabilities(jsonpatch.Predicate),
	)

	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, jsonpatch.ErrUnsupportedCapability)
}

func TestCompileJSONAppliesCapabilityPolicyToNestedPredicates(t *testing.T) {
	t.Parallel()

	data := []byte(`[{"op":"and","path":"/profile","apply":[{"op":"matches","path":"/name","value":"^Ada$"}]}]`)
	patch, err := jsonpatch.CompileJSON(
		data,
		jsonpatch.WithCapabilities(jsonpatch.Predicate),
	)
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, jsonpatch.ErrUnsupportedCapability)

	var patchErr *jsonpatch.Error
	require.ErrorAs(t, err, &patchErr)
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "matches", patchErr.Op())
	assert.Equal(t, "/profile/name", patchErr.Path())
	assert.Equal(t, "json", patchErr.Codec())

	patch, err = jsonpatch.CompileJSON(
		data,
		jsonpatch.WithCapabilities(jsonpatch.Predicate, jsonpatch.RegexPredicate),
	)
	require.NoError(t, err)
	assert.Equal(t, 1, patch.Len())
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
	operands := []op.PredicateOp{operand}
	operation := op.NewAnd(nil, operands...)

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

func TestApplyInPlaceFailedMovePreservesDocument(t *testing.T) {
	t.Parallel()

	doc := map[string]any{
		"source": "Ada",
		"target": "scalar",
	}
	want := map[string]any{
		"source": "Ada",
		"target": "scalar",
	}
	patch, err := jsonpatch.Compile(
		op.NewMove([]string{"target", "name"}, []string{"source"}),
	)
	require.NoError(t, err)

	err = jsonpatch.ApplyInPlace(patch, &doc)
	require.Error(t, err)
	assert.ErrorIs(t, err, jsonpatch.ErrRuntimeConflict)
	assert.ErrorIs(t, err, op.ErrCannotAddToValue)
	if diff := cmp.Diff(want, doc); diff != "" {
		t.Errorf("failed move mutated document (-want +got):\n%s", diff)
	}
}

func TestApplyInPlaceKeepsPriorOperationsWhenMoveFails(t *testing.T) {
	t.Parallel()

	doc := map[string]any{
		"source": "Ada",
		"status": "before",
		"target": "scalar",
	}
	patch, err := jsonpatch.Compile(
		op.NewReplace([]string{"status"}, "updated"),
		op.NewMove([]string{"target", "name"}, []string{"source"}),
	)
	require.NoError(t, err)

	err = jsonpatch.ApplyInPlace(patch, &doc)
	require.Error(t, err)
	assert.Equal(t, map[string]any{
		"source": "Ada",
		"status": "updated",
		"target": "scalar",
	}, doc)
}

func TestApplyInPlaceKeepsSuccessfulPrefixAcrossDocumentShapes(t *testing.T) {
	t.Parallel()

	prefixPatch, err := jsonpatch.Compile(
		op.NewReplace([]string{"status"}, "updated"),
		op.NewRemove([]string{"missing"}),
	)
	require.NoError(t, err)

	t.Run("root replacement", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.Compile(
			op.NewReplace(nil, "updated"),
			op.NewRemove([]string{"missing"}),
		)
		require.NoError(t, err)
		doc := "before"
		err = jsonpatch.ApplyInPlace(patch, &doc)
		require.ErrorIs(t, err, op.ErrCannotRemoveFromValue)
		assert.Equal(t, "updated", doc)
	})

	t.Run("slice growth", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.Compile(
			op.NewAdd([]string{"-"}, "second"),
			op.NewRemove([]string{"9"}),
		)
		require.NoError(t, err)
		doc := []any{"first"}
		err = jsonpatch.ApplyInPlace(patch, &doc)
		require.ErrorIs(t, err, op.ErrIndexOutOfRange)
		assert.Equal(t, []any{"first", "second"}, doc)
	})

	t.Run("JSON text", func(t *testing.T) {
		t.Parallel()

		doc := jsonpatch.JSONText(`{"status":"before"}`)
		err := jsonpatch.ApplyInPlace(prefixPatch, &doc)
		require.ErrorIs(t, err, op.ErrPathNotFound)
		assert.JSONEq(t, `{"status":"updated"}`, string(doc))
	})

	t.Run("JSON bytes", func(t *testing.T) {
		t.Parallel()

		doc := []byte(`{"status":"before"}`)
		err := jsonpatch.ApplyInPlace(prefixPatch, &doc)
		require.ErrorIs(t, err, op.ErrPathNotFound)
		assert.JSONEq(t, `{"status":"updated"}`, string(doc))
	})

	t.Run("typed map", func(t *testing.T) {
		t.Parallel()

		doc := map[string]string{"status": "before"}
		err := jsonpatch.ApplyInPlace(prefixPatch, &doc)
		require.ErrorIs(t, err, op.ErrPathNotFound)
		assert.Equal(t, map[string]string{"status": "updated"}, doc)
	})

	t.Run("struct", func(t *testing.T) {
		t.Parallel()

		type document struct {
			Status string `json:"status"`
		}
		doc := document{Status: "before"}
		err := jsonpatch.ApplyInPlace(prefixPatch, &doc)
		require.ErrorIs(t, err, op.ErrPathNotFound)
		assert.Equal(t, document{Status: "updated"}, doc)
	})
}

func TestApplyInPlaceFirstFailurePreservesOriginalCarrier(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(op.NewRemove([]string{"missing"}))
	require.NoError(t, err)

	assertFailure := func(t *testing.T, err error) {
		t.Helper()
		require.ErrorIs(t, err, op.ErrPathNotFound)

		var patchErr *jsonpatch.Error
		require.ErrorAs(t, err, &patchErr)
		assert.Equal(t, 0, patchErr.Index())
		assert.Equal(t, "remove", patchErr.Op())
		assert.Equal(t, "/missing", patchErr.Path())
	}

	t.Run("JSON bytes remain byte-identical", func(t *testing.T) {
		t.Parallel()

		doc := []byte("{\n  \"status\": \"before\"\n}\n")
		want := slices.Clone(doc)
		err := jsonpatch.ApplyInPlace(patch, &doc)
		assertFailure(t, err)
		assert.Equal(t, want, doc)
	})

	t.Run("struct retains fields outside JSON", func(t *testing.T) {
		t.Parallel()

		type document struct {
			Status string `json:"status"`
			Local  string `json:"-"`
		}
		doc := document{Status: "before", Local: "keep"}
		want := doc
		err := jsonpatch.ApplyInPlace(patch, &doc)
		assertFailure(t, err)
		assert.Equal(t, want, doc)
	})
}

func TestApplyInPlaceReportsPrefixConversionFailure(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(
		op.NewReplace([]string{"status"}, 1),
		op.NewRemove([]string{"missing"}),
	)
	require.NoError(t, err)
	doc := map[string]string{"status": "before"}
	err = jsonpatch.ApplyInPlace(patch, &doc)
	require.Error(t, err)
	assert.ErrorIs(t, err, op.ErrPathNotFound)
	assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
	assert.Equal(t, map[string]string{"status": "before"}, doc)
}

func TestApplyInPlaceConvertsAtDocumentBoundary(t *testing.T) {
	t.Parallel()

	t.Run("named typed document can recover before final conversion", func(t *testing.T) {
		t.Parallel()

		type document map[string]string
		patch, err := jsonpatch.Compile(
			op.NewReplace([]string{"status"}, 1),
			op.NewReplace([]string{"status"}, "updated"),
		)
		require.NoError(t, err)

		doc := document{"status": "before"}
		err = jsonpatch.ApplyInPlace(patch, &doc)
		require.NoError(t, err)
		assert.Equal(t, document{"status": "updated"}, doc)
	})

	t.Run("JSON text can recover before final encoding", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.Compile(
			op.NewAdd([]string{"bad"}, math.Inf(1)),
			op.NewRemove([]string{"bad"}),
		)
		require.NoError(t, err)

		doc := jsonpatch.JSONText(`{"status":"before"}`)
		err = jsonpatch.ApplyInPlace(patch, &doc)
		require.NoError(t, err)
		assert.JSONEq(t, `{"status":"before"}`, string(doc))
	})

	t.Run("JSON bytes preserve both failures when prefix cannot encode", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.Compile(
			op.NewAdd([]string{"bad"}, math.Inf(1)),
			op.NewRemove([]string{"missing"}),
		)
		require.NoError(t, err)

		doc := []byte(`{"status":"before"}`)
		err = jsonpatch.ApplyInPlace(patch, &doc)
		require.Error(t, err)
		assert.ErrorIs(t, err, op.ErrPathNotFound)
		assert.ErrorIs(t, err, jsonpatch.ErrConversionFailed)
		assert.JSONEq(t, `{"status":"before"}`, string(doc))

		var patchErr *jsonpatch.Error
		require.ErrorAs(t, err, &patchErr)
		assert.Equal(t, 1, patchErr.Index())
		assert.Equal(t, "remove", patchErr.Op())
		assert.Equal(t, "/missing", patchErr.Path())
	})
}

func TestApplyInPlaceJSONBytesKeepsPrefixAndRestoresFailedMove(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.Compile(
		op.NewReplace([]string{"status"}, "updated"),
		op.NewMove([]string{"target", "name"}, []string{"source"}),
	)
	require.NoError(t, err)
	doc := []byte(`{"source":"Ada","status":"before","target":"scalar"}`)
	err = jsonpatch.ApplyInPlace(patch, &doc)
	require.ErrorIs(t, err, op.ErrCannotAddToValue)
	assert.JSONEq(t, `{"source":"Ada","status":"updated","target":"scalar"}`, string(doc))
}

func TestApplyFailedMovePreservesInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		doc  map[string]any
		path []string
	}{
		{
			name: "missing destination parent",
			doc:  map[string]any{"source": "Ada"},
			path: []string{"missing", "name"},
		},
		{
			name: "scalar destination parent",
			doc:  map[string]any{"source": "Ada", "target": "scalar"},
			path: []string{"target", "name"},
		},
		{
			name: "array destination out of range",
			doc:  map[string]any{"source": "Ada", "items": []any{"first"}},
			path: []string{"items", "2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			want := maps.Clone(tc.doc)
			patch, err := jsonpatch.Compile(op.NewMove(tc.path, []string{"source"}))
			require.NoError(t, err)
			result, err := jsonpatch.Apply(patch, tc.doc)
			require.Error(t, err)
			assert.Nil(t, result)
			if diff := cmp.Diff(want, tc.doc); diff != "" {
				t.Errorf("failed immutable move mutated input (-want +got):\n%s", diff)
			}
		})
	}
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

func TestPatchApplyClassifiesPredicateMissingPathAsTestFailure(t *testing.T) {
	t.Parallel()

	patch, err := jsonpatch.CompileOps(
		[]jsonpatch.Op{op.NewType([]string{"missing"}, "string")},
		jsonpatch.WithCapabilities(jsonpatch.Predicate),
	)
	require.NoError(t, err)

	result, err := jsonpatch.Apply(patch, map[string]any{"name": "Ada"})
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, jsonpatch.ErrTestFailed)
	assert.ErrorIs(t, err, op.ErrPathNotFound)
	assert.NotErrorIs(t, err, jsonpatch.ErrRuntimeConflict)

	var patchErr *jsonpatch.Error
	require.ErrorAs(t, err, &patchErr)
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "type", patchErr.Op())
	assert.Equal(t, "/missing", patchErr.Path())
}

func TestPatchApplyPredicateFailuresPreserveSpecificCause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		op      jsonpatch.Op
		doc     map[string]any
		wantErr error
	}{
		{name: "wrong target type", op: op.NewContains([]string{"count"}, "5"), doc: map[string]any{"count": 5.0}, wantErr: op.ErrNotString},
		{name: "type mismatch", op: op.NewType([]string{"name"}, "number"), doc: map[string]any{"name": "Ada"}, wantErr: op.ErrTypeMismatch},
		{name: "test_type mismatch", op: op.NewTestType([]string{"name"}, "number"), doc: map[string]any{"name": "Ada"}, wantErr: op.ErrTypeMismatch},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			patch, err := jsonpatch.CompileOps(
				[]jsonpatch.Op{tc.op},
				jsonpatch.WithCapabilities(jsonpatch.Predicate),
			)
			require.NoError(t, err)
			result, err := jsonpatch.Apply(patch, tc.doc)
			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, jsonpatch.ErrTestFailed)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.NotErrorIs(t, err, jsonpatch.ErrTypeMismatch)
			assert.NotErrorIs(t, err, jsonpatch.ErrRuntimeConflict)
		})
	}
}

func TestPatchApplyClassifiesRuntimeValueTypeFailure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		op      jsonpatch.Op
		doc     any
		wantErr error
		wantOp  string
		path    string
	}{
		{name: "str_ins", op: op.NewStrIns([]string{"count"}, 0, "x"), doc: map[string]any{"count": 5.0}, wantErr: op.ErrNotString, wantOp: "str_ins", path: "/count"},
		{name: "str_del", op: op.NewStrDel([]string{"count"}, 0, 1), doc: map[string]any{"count": 5.0}, wantErr: op.ErrNotString, wantOp: "str_del", path: "/count"},
		{name: "inc root", op: op.NewInc(nil, 1), doc: map[string]any{"count": 5.0}, wantErr: op.ErrNotNumber, wantOp: "inc"},
		{name: "inc nested", op: op.NewInc([]string{"count"}, 1), doc: map[string]any{"count": "five"}, wantErr: op.ErrNotNumber, wantOp: "inc", path: "/count"},
		{name: "extend", op: op.NewExtend([]string{"count"}, map[string]any{"next": 6}, false), doc: map[string]any{"count": 5.0}, wantErr: op.ErrNotObject, wantOp: "extend", path: "/count"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			patch, err := jsonpatch.CompileOps(
				[]jsonpatch.Op{tc.op},
				jsonpatch.WithCapabilities(jsonpatch.Extended),
			)
			require.NoError(t, err)
			result, err := jsonpatch.Apply(patch, tc.doc)
			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, jsonpatch.ErrTypeMismatch)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.NotErrorIs(t, err, jsonpatch.ErrRuntimeConflict)

			var patchErr *jsonpatch.Error
			require.ErrorAs(t, err, &patchErr)
			assert.Equal(t, 0, patchErr.Index())
			assert.Equal(t, tc.wantOp, patchErr.Op())
			assert.Equal(t, tc.path, patchErr.Path())
		})
	}
}

func TestPatchApplyKeepsNonTypeFailuresAsRuntimeConflicts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		op      jsonpatch.Op
		doc     any
		wantErr error
	}{
		{name: "missing path", op: op.NewStrIns([]string{"missing"}, 0, "x"), doc: map[string]any{"text": "Ada"}, wantErr: op.ErrPathNotFound},
		{name: "index out of range", op: op.NewMerge([]string{"items", "2"}, 0, nil), doc: map[string]any{"items": []any{"Ada"}}, wantErr: op.ErrIndexOutOfRange},
		{name: "invalid position", op: op.NewStrIns([]string{"value"}, 1, "x"), doc: map[string]any{"value": nil}, wantErr: op.ErrInvalidPosition},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			patch, err := jsonpatch.CompileOps(
				[]jsonpatch.Op{tc.op},
				jsonpatch.WithCapabilities(jsonpatch.Extended),
			)
			require.NoError(t, err)
			result, err := jsonpatch.Apply(patch, tc.doc)
			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, jsonpatch.ErrRuntimeConflict)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.NotErrorIs(t, err, jsonpatch.ErrTypeMismatch)
		})
	}
}

func FuzzApplyPositionedOperationsDoesNotPanic(f *testing.F) {
	f.Add("abc", "a", 0, 1)
	f.Add("", "", -1, 1)
	f.Add("text", "value", 100, 100)

	f.Fuzz(func(_ *testing.T, text, needle string, testPos, mergePos int) {
		testPatch, err := jsonpatch.CompileOps(
			[]jsonpatch.Op{op.NewTestString(nil, needle, testPos, false)},
			jsonpatch.WithCapabilities(jsonpatch.Predicate),
		)
		if err == nil {
			_, _ = jsonpatch.Apply(testPatch, text)
		}

		mergePatch, err := jsonpatch.CompileOps(
			[]jsonpatch.Op{op.NewMerge([]string{"1"}, mergePos, nil)},
			jsonpatch.WithCapabilities(jsonpatch.Extended),
		)
		if err == nil {
			_, _ = jsonpatch.Apply(mergePatch, []any{text, needle})
		}
	})
}
