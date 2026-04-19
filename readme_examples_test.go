package jsonpatch_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/codec/binary"
	"github.com/kaptinlin/jsonpatch/codec/compact"
	"github.com/kaptinlin/jsonpatch/op"
)

type readmeUser struct {
	Name   string   `json:"name"`
	Active bool     `json:"active"`
	Roles  []string `json:"roles"`
}

func TestREADMEExamples(t *testing.T) {
	t.Parallel()

	t.Run("quick start preserves document shape and immutability", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"name": "John",
			"tags": []any{"golang"},
		}
		patch := []jsonpatch.Operation{
			{Op: "test", Path: "/name", Value: "John"},
			{Op: "replace", Path: "/name", Value: "Jane"},
			{Op: "add", Path: "/email", Value: "jane@example.com"},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		require.NoError(t, err)
		assert.Equal(t, "Jane", result.Doc["name"])
		assert.Equal(t, "jane@example.com", result.Doc["email"])
		assert.Equal(t, "John", doc["name"])
	})

	t.Run("executable operations work with ApplyOps", func(t *testing.T) {
		t.Parallel()

		user := readmeUser{Name: "John", Active: true, Roles: []string{"admin"}}
		ops := []jsonpatch.Op{
			op.NewTest([]string{"active"}, true),
			op.NewReplace([]string{"name"}, "Jane"),
			op.NewAdd([]string{"roles", "-"}, "owner"),
		}

		result, err := jsonpatch.ApplyOps(user, ops)
		require.NoError(t, err)
		assert.Equal(t, "Jane", result.Doc.Name)
		assert.Equal(t, []string{"admin", "owner"}, result.Doc.Roles)
		assert.Equal(t, "John", user.Name)
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

	t.Run("binary codec rejects second-order predicates", func(t *testing.T) {
		t.Parallel()

		codec := binary.New()
		_, err := codec.Encode([]jsonpatch.Op{op.NewNot(op.NewTest([]string{"active"}, true))})
		require.Error(t, err)
		assert.True(t, errors.Is(err, binary.ErrUnsupportedOp))
	})
}
