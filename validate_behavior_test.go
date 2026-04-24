package jsonpatch_test

import (
	"testing"

	"github.com/kaptinlin/jsonpointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
)

func TestValidateOperationRejectsInvalidPayloads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation jsonpatch.Operation
		allow     bool
		wantErr   error
	}{
		{name: "invalid pointer", operation: jsonpatch.Operation{Op: "test", Path: "name", Value: "Ada"}, wantErr: jsonpatch.ErrInvalidJSONPointer},
		{name: "replace missing value", operation: jsonpatch.Operation{Op: "replace", Path: "/name"}, wantErr: jsonpatch.ErrMissingValue},
		{name: "copy missing from", operation: jsonpatch.Operation{Op: "copy", Path: "/name"}, wantErr: jsonpatch.ErrMissingFrom},
		{name: "copy invalid from", operation: jsonpatch.Operation{Op: "copy", Path: "/name", From: "name"}, wantErr: jsonpointer.ErrPointerInvalid},
		{name: "move invalid from", operation: jsonpatch.Operation{Op: "move", Path: "/name", From: "name"}, wantErr: jsonpointer.ErrPointerInvalid},
		{name: "move into own child", operation: jsonpatch.Operation{Op: "move", Path: "/profile/name", From: "/profile"}, wantErr: jsonpatch.ErrCannotMoveToChildren},
		{name: "test missing value", operation: jsonpatch.Operation{Op: "test", Path: "/name"}, wantErr: jsonpatch.ErrMissingValue},
		{name: "str_ins negative position", operation: jsonpatch.Operation{Op: "str_ins", Path: "/name", Pos: -1}, wantErr: jsonpatch.ErrNegativeNumber},
		{name: "str_del negative position", operation: jsonpatch.Operation{Op: "str_del", Path: "/name", Pos: -1}, wantErr: jsonpatch.ErrNegativeNumber},
		{name: "str_del negative length", operation: jsonpatch.Operation{Op: "str_del", Path: "/name", Len: -1}, wantErr: jsonpatch.ErrNegativeNumber},
		{name: "merge position must be positive", operation: jsonpatch.Operation{Op: "merge", Path: "/nodes/1"}, wantErr: jsonpatch.ErrPosGreaterThanZero},
		{name: "test_type missing type", operation: jsonpatch.Operation{Op: "test_type", Path: "/name"}, wantErr: jsonpatch.ErrInvalidTypeField},
		{name: "test_type empty string", operation: jsonpatch.Operation{Op: "test_type", Path: "/name", Type: ""}, wantErr: jsonpatch.ErrInvalidTypeField},
		{name: "test_type invalid string", operation: jsonpatch.Operation{Op: "test_type", Path: "/name", Type: "invalid"}, wantErr: jsonpatch.ErrInvalidType},
		{name: "test_type empty array", operation: jsonpatch.Operation{Op: "test_type", Path: "/name", Type: []any{}}, wantErr: jsonpatch.ErrInvalidTypeField},
		{name: "test_type non-string item", operation: jsonpatch.Operation{Op: "test_type", Path: "/name", Type: []any{"string", 1}}, wantErr: jsonpatch.ErrInvalidType},
		{name: "test_type invalid field type", operation: jsonpatch.Operation{Op: "test_type", Path: "/name", Type: 1}, wantErr: jsonpatch.ErrInvalidType},
		{name: "test_string negative position", operation: jsonpatch.Operation{Op: "test_string", Path: "/name", Pos: -1}, wantErr: jsonpatch.ErrNegativeNumber},
		{name: "test_string_len negative length", operation: jsonpatch.Operation{Op: "test_string_len", Path: "/name", Len: -1}, wantErr: jsonpatch.ErrNegativeNumber},
		{name: "matches disallowed", operation: jsonpatch.Operation{Op: "matches", Path: "/name", Value: "^A"}, wantErr: jsonpatch.ErrMatchesNotAllowed},
		{name: "matches value missing", operation: jsonpatch.Operation{Op: "matches", Path: "/name"}, allow: true, wantErr: jsonpatch.ErrExpectedValueToBeString},
		{name: "contains value missing", operation: jsonpatch.Operation{Op: "contains", Path: "/name"}, wantErr: jsonpatch.ErrExpectedValueToBeString},
		{name: "ends value not string", operation: jsonpatch.Operation{Op: "ends", Path: "/name", Value: 1}, wantErr: jsonpatch.ErrExpectedValueToBeString},
		{name: "starts value not string", operation: jsonpatch.Operation{Op: "starts", Path: "/name", Value: 1}, wantErr: jsonpatch.ErrExpectedValueToBeString},
		{name: "in value missing", operation: jsonpatch.Operation{Op: "in", Path: "/role"}, wantErr: jsonpatch.ErrInOperationValueMustBeArray},
		{name: "in value not array", operation: jsonpatch.Operation{Op: "in", Path: "/role", Value: "admin"}, wantErr: jsonpatch.ErrInOperationValueMustBeArray},
		{name: "more value missing", operation: jsonpatch.Operation{Op: "more", Path: "/score"}, wantErr: jsonpatch.ErrValueMustBeNumber},
		{name: "more value not number", operation: jsonpatch.Operation{Op: "more", Path: "/score", Value: "1"}, wantErr: jsonpatch.ErrValueMustBeNumber},
		{name: "less value missing", operation: jsonpatch.Operation{Op: "less", Path: "/score"}, wantErr: jsonpatch.ErrValueMustBeNumber},
		{name: "type value missing", operation: jsonpatch.Operation{Op: "type", Path: "/name"}, wantErr: jsonpatch.ErrExpectedValueToBeString},
		{name: "type invalid value", operation: jsonpatch.Operation{Op: "type", Path: "/name", Value: "invalid"}, wantErr: jsonpatch.ErrInvalidType},
		{name: "and empty apply", operation: jsonpatch.Operation{Op: "and", Path: "/profile"}, wantErr: jsonpatch.ErrEmptyPredicateList},
		{name: "composite invalid child", operation: jsonpatch.Operation{Op: "or", Path: "/profile", Apply: []jsonpatch.Operation{{Op: "test", Path: "/name"}}}, wantErr: jsonpatch.ErrMissingValue},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := jsonpatch.ValidateOperation(tc.operation, tc.allow)
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestValidateOperationAcceptsPredicateFamilies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation jsonpatch.Operation
		allow     bool
	}{
		{name: "test_type string", operation: jsonpatch.Operation{Op: "test_type", Path: "/name", Type: "string"}},
		{name: "test_type array", operation: jsonpatch.Operation{Op: "test_type", Path: "/name", Type: []any{"string", "null"}}},
		{name: "matches allowed", operation: jsonpatch.Operation{Op: "matches", Path: "/name", Value: "^A"}, allow: true},
		{name: "contains", operation: jsonpatch.Operation{Op: "contains", Path: "/name", Value: "A"}},
		{name: "ends", operation: jsonpatch.Operation{Op: "ends", Path: "/name", Value: "a"}},
		{name: "starts", operation: jsonpatch.Operation{Op: "starts", Path: "/name", Value: "A"}},
		{name: "in", operation: jsonpatch.Operation{Op: "in", Path: "/role", Value: []any{"admin", "editor"}}},
		{name: "more", operation: jsonpatch.Operation{Op: "more", Path: "/score", Value: 5}},
		{name: "less", operation: jsonpatch.Operation{Op: "less", Path: "/score", Value: 10}},
		{name: "type", operation: jsonpatch.Operation{Op: "type", Path: "/name", Value: "string"}},
		{name: "defined", operation: jsonpatch.Operation{Op: "defined", Path: "/name"}},
		{name: "undefined", operation: jsonpatch.Operation{Op: "undefined", Path: "/deleted"}},
		{name: "composite", operation: jsonpatch.Operation{Op: "and", Path: "/profile", Apply: []jsonpatch.Operation{{Op: "defined", Path: "/profile/name"}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := jsonpatch.ValidateOperation(tc.operation, tc.allow)
			require.NoError(t, err)
		})
	}
}
