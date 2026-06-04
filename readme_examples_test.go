package jsonpatch_test

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/codec/binary"
	"github.com/kaptinlin/jsonpatch/codec/compact"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/op"
)

type readmeUser struct {
	Name   string   `json:"name"`
	Active bool     `json:"active"`
	Roles  []string `json:"roles"`
}

func TestREADMEExamples(t *testing.T) {
	t.Parallel()

	t.Run("quick start compiles and applies immutably", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"name": "John",
			"tags": []any{"golang"},
		}

		patch, err := jsonpatch.Compile(
			op.NewTest([]string{"name"}, "John"),
			op.NewReplace([]string{"name"}, "Jane"),
			op.NewAdd([]string{"email"}, "jane@example.com"),
		)
		require.NoError(t, err)

		result, err := jsonpatch.Apply(patch, doc)
		require.NoError(t, err)
		assert.Equal(t, "Jane", result.Doc["name"])
		assert.Equal(t, "jane@example.com", result.Doc["email"])
		assert.Equal(t, "John", doc["name"])
	})

	t.Run("executable operations compile into a patch", func(t *testing.T) {
		t.Parallel()

		user := readmeUser{Name: "John", Active: true, Roles: []string{"admin"}}

		patch, err := jsonpatch.Compile(
			op.NewTest([]string{"active"}, true),
			op.NewReplace([]string{"name"}, "Jane"),
			op.NewAdd([]string{"roles", "-"}, "owner"),
		)
		require.NoError(t, err)

		result, err := jsonpatch.Apply(patch, user)
		require.NoError(t, err)
		assert.Equal(t, "Jane", result.Doc.Name)
		assert.Equal(t, []string{"admin", "owner"}, result.Doc.Roles)
		assert.Equal(t, "John", user.Name)
	})

	t.Run("capabilities gate predicate operations", func(t *testing.T) {
		t.Parallel()

		data := []byte(`[{"op":"contains","path":"/name","value":"Ada"}]`)
		_, err := jsonpatch.CompileJSON(data)
		require.Error(t, err)
		assert.ErrorIs(t, err, jsonpatch.ErrUnsupportedCapability)

		patch, err := jsonpatch.CompileJSON(data,
			jsonpatch.WithCapabilities(jsonpatch.RFC6902, jsonpatch.Predicate),
		)
		require.NoError(t, err)

		result, err := jsonpatch.Apply(patch, map[string]any{"name": "Ada Lovelace"})
		require.NoError(t, err)
		assert.Equal(t, "Ada Lovelace", result.Doc["name"])
	})

	t.Run("JSONText marks string documents as JSON", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.CompileJSON([]byte(`[{"op":"replace","path":"/name","value":"Jane"}]`))
		require.NoError(t, err)

		result, err := jsonpatch.Apply(patch, jsonpatch.JSONText(`{"name":"John"}`))
		require.NoError(t, err)

		var got map[string]any
		require.NoError(t, json.Unmarshal([]byte(result.Doc), &got))
		assert.Equal(t, "Jane", got["name"])
	})

	t.Run("ApplyInPlace writes result back to the input variable", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"name": "John"}
		patch, err := jsonpatch.Compile(op.NewReplace([]string{"name"}, "Jane"))
		require.NoError(t, err)

		require.NoError(t, jsonpatch.ApplyInPlace(patch, &doc))
		assert.Equal(t, "Jane", doc["name"])
	})

	t.Run("structured errors expose failure context", func(t *testing.T) {
		t.Parallel()

		patch, err := jsonpatch.Compile(op.NewTest([]string{"name"}, "Ada"))
		require.NoError(t, err)

		result, err := jsonpatch.Apply(patch, map[string]any{"name": "Grace"})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, jsonpatch.ErrTestFailed)

		var patchErr *jsonpatch.Error
		require.ErrorAs(t, err, &patchErr)
		assert.Equal(t, 0, patchErr.Index())
		assert.Equal(t, "test", patchErr.Op())
		assert.Equal(t, "/name", patchErr.Path())
		assert.NotNil(t, patchErr.Cause())
	})

	t.Run("JSON operation values compile through the JSON codec", func(t *testing.T) {
		t.Parallel()

		operations := []jsoncodec.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
		}

		patch, err := jsonpatch.CompileOperations(operations)
		require.NoError(t, err)

		result, err := jsonpatch.Apply(patch, map[string]any{"name": "John"})
		require.NoError(t, err)
		assert.Equal(t, "Jane", result.Doc["name"])
	})

	t.Run("compact codec round-trips executable operations", func(t *testing.T) {
		t.Parallel()

		ops := []jsonpatch.Op{
			op.NewAdd([]string{"name"}, "Jane"),
			op.NewInc([]string{"version"}, 1),
		}

		encoded, err := compact.EncodeJSON(ops)
		require.NoError(t, err)

		decoded, err := compact.DecodeJSON(encoded)
		require.NoError(t, err)
		require.Len(t, decoded, 2)
		assert.Equal(t, jsonpatch.OpAddType, decoded[0].Op())
		assert.Equal(t, []string{"name"}, decoded[0].Path())
		assert.Equal(t, jsonpatch.OpIncType, decoded[1].Op())
		assert.Equal(t, []string{"version"}, decoded[1].Path())
	})

	t.Run("binary codec round-trips executable operations", func(t *testing.T) {
		t.Parallel()

		ops := []jsonpatch.Op{
			op.NewAdd([]string{"name"}, "Jane"),
			op.NewInc([]string{"version"}, 1),
		}

		codec := binary.New()
		data, err := codec.Encode(ops)
		require.NoError(t, err)

		decoded, err := codec.Decode(data)
		require.NoError(t, err)
		require.Len(t, decoded, 2)
		assert.Equal(t, jsonpatch.OpAddType, decoded[0].Op())
		assert.Equal(t, jsonpatch.OpIncType, decoded[1].Op())
	})

	t.Run("binary codec round-trips second-order predicates", func(t *testing.T) {
		t.Parallel()

		codec := binary.New()
		ops := []jsonpatch.Op{
			op.NewNot(op.NewTest([]string{"active"}, true)),
		}
		data, err := codec.Encode(ops)
		require.NoError(t, err)

		decoded, err := codec.Decode(data)
		require.NoError(t, err)
		require.Len(t, decoded, 1)
		assert.Equal(t, jsonpatch.OpNotType, decoded[0].Op())
	})
}
