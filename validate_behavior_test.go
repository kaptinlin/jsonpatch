package jsonpatch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestCompileOperationsRejectsInvalidPayloads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation jsoncodec.Operation
		data      []byte
		opts      []jsonpatch.CompileOption
		wantKind  error
		wantCause error
		wantOp    string
		wantPath  string
	}{
		{
			name:      "missing operation name",
			operation: jsoncodec.Operation{Path: "/name"},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: jsoncodec.ErrCodecOpUnknown,
			wantPath:  "/name",
		},
		{
			name:      "invalid path pointer",
			operation: jsoncodec.Operation{Op: "test", Path: "name", Value: "Ada"},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: jsoncodec.ErrInvalidPointer,
			wantOp:    "test",
			wantPath:  "name",
		},
		{
			name:      "copy invalid from pointer",
			operation: jsoncodec.Operation{Op: "copy", Path: "/name", From: "name"},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: jsoncodec.ErrInvalidPointer,
			wantOp:    "copy",
			wantPath:  "/name",
		},
		{
			name:      "move into own child",
			operation: jsoncodec.Operation{Op: "move", Path: "/profile/name", From: "/profile"},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: op.ErrCannotMoveIntoChildren,
			wantOp:    "move",
			wantPath:  "/profile/name",
		},
		{
			name:      "str_del negative length",
			operation: jsoncodec.Operation{Op: "str_del", Path: "/name", Pos: 0, Len: -1},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: op.ErrLengthNegative,
			wantOp:    "str_del",
			wantPath:  "/name",
		},
		{
			name:      "matches requires regex capability",
			operation: jsoncodec.Operation{Op: "matches", Path: "/name", Value: "^A"},
			wantKind:  jsonpatch.ErrUnsupportedCapability,
			wantOp:    "matches",
			wantPath:  "/name",
		},
		{
			name:      "test_type missing type",
			operation: jsoncodec.Operation{Op: "test_type", Path: "/name"},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: jsoncodec.ErrTestTypeOpMissingType,
			wantOp:    "test_type",
			wantPath:  "/name",
		},
		{
			name:      "test_type invalid string",
			operation: jsoncodec.Operation{Op: "test_type", Path: "/name", Type: "invalid"},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: jsoncodec.ErrInvalidType,
			wantOp:    "test_type",
			wantPath:  "/name",
		},
		{
			name:      "in requires array value",
			operation: jsoncodec.Operation{Op: "in", Path: "/role", Value: "admin"},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: jsoncodec.ErrInOpValueMustBeArray,
			wantOp:    "in",
			wantPath:  "/role",
		},
		{
			name:      "type invalid value",
			operation: jsoncodec.Operation{Op: "type", Path: "/name", Value: "invalid"},
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: op.ErrInvalidType,
			wantOp:    "type",
			wantPath:  "/name",
		},
		{
			name:      "merge position must be positive",
			data:      []byte(`[{"op":"merge","path":"/nodes/1","pos":0}]`),
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: op.ErrPositionOutOfBounds,
			wantOp:    "merge",
			wantPath:  "/nodes/1",
		},
		{
			name:      "and requires predicate operands",
			data:      []byte(`[{"op":"and","path":"/profile","apply":[]}]`),
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: op.ErrInvalidPredicateInAnd,
			wantOp:    "and",
			wantPath:  "/profile",
		},
		{
			name:      "composite rejects mutation child",
			data:      []byte(`[{"op":"and","path":"/profile","apply":[{"op":"add","path":"/profile/name","value":"Ada"}]}]`),
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: jsoncodec.ErrInvalidPredicateOperand,
			wantOp:    "and",
			wantPath:  "/profile",
		},
		{
			name:      "not requires one predicate",
			data:      []byte(`[{"op":"not","path":"/profile","apply":[{"op":"defined","path":"/profile/name"},{"op":"defined","path":"/profile/role"}]}]`),
			opts:      allCapabilities(),
			wantKind:  jsonpatch.ErrPayloadInvalid,
			wantCause: jsoncodec.ErrNotOpRequiresSingleOperand,
			wantOp:    "not",
			wantPath:  "/profile",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := compileOne(&tc.operation, tc.data, tc.opts...)
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.wantKind)
			if tc.wantCause != nil {
				assert.ErrorIs(t, err, tc.wantCause)
			}

			var patchErr *jsonpatch.Error
			require.ErrorAs(t, err, &patchErr)
			assert.Equal(t, 0, patchErr.Index())
			assert.Equal(t, tc.wantOp, patchErr.Op())
			assert.Equal(t, tc.wantPath, patchErr.Path())
			assert.Equal(t, "json", patchErr.Codec())
		})
	}
}

func TestCompileOperationsAcceptsEnabledFamilies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation jsoncodec.Operation
	}{
		{name: "root add with null", operation: jsoncodec.Operation{Op: "add", Path: "", Value: nil}},
		{name: "root replace with null", operation: jsoncodec.Operation{Op: "replace", Path: "", Value: nil}},
		{name: "root remove", operation: jsoncodec.Operation{Op: "remove", Path: ""}},
		{name: "root test with null", operation: jsoncodec.Operation{Op: "test", Path: "", Value: nil}},
		{name: "test_type string", operation: jsoncodec.Operation{Op: "test_type", Path: "/name", Type: "string"}},
		{name: "test_type array", operation: jsoncodec.Operation{Op: "test_type", Path: "/name", Type: []any{"string", "null"}}},
		{name: "test_type string array", operation: jsoncodec.Operation{Op: "test_type", Path: "/name", Type: []string{"string", "null"}}},
		{name: "matches", operation: jsoncodec.Operation{Op: "matches", Path: "/name", Value: "^A"}},
		{name: "contains root", operation: jsoncodec.Operation{Op: "contains", Path: "", Value: "A"}},
		{name: "ends", operation: jsoncodec.Operation{Op: "ends", Path: "/name", Value: "a"}},
		{name: "starts", operation: jsoncodec.Operation{Op: "starts", Path: "/name", Value: "A"}},
		{name: "in", operation: jsoncodec.Operation{Op: "in", Path: "/role", Value: []any{"admin", "editor"}}},
		{name: "more", operation: jsoncodec.Operation{Op: "more", Path: "/score", Value: 5}},
		{name: "less", operation: jsoncodec.Operation{Op: "less", Path: "/score", Value: 10}},
		{name: "type", operation: jsoncodec.Operation{Op: "type", Path: "/name", Value: "string"}},
		{name: "defined root", operation: jsoncodec.Operation{Op: "defined", Path: ""}},
		{name: "undefined", operation: jsoncodec.Operation{Op: "undefined", Path: "/deleted"}},
		{name: "composite", operation: jsoncodec.Operation{Op: "and", Path: "/profile", Apply: []jsoncodec.Operation{{Op: "defined", Path: "/profile/name"}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := compileOne(&tc.operation, nil, allCapabilities()...)
			require.NoError(t, err)
		})
	}

	patch, err := jsonpatch.CompileOperations(nil, allCapabilities()...)
	require.NoError(t, err)
	assert.Equal(t, 0, patch.Len())
}

func compileOne(operation *jsoncodec.Operation, data []byte, opts ...jsonpatch.CompileOption) error {
	if data != nil {
		_, err := jsonpatch.CompileJSON(data, opts...)
		return err
	}
	_, err := jsonpatch.CompileOperations([]jsoncodec.Operation{*operation}, opts...)
	return err
}

func allCapabilities() []jsonpatch.CompileOption {
	return []jsonpatch.CompileOption{jsonpatch.WithCapabilities(jsonpatch.AllCapabilities)}
}
