package jsonpatch_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// OPERATIONS TESTS
// =============================================================================

func TestValidateOperations(t *testing.T) {
	t.Run("throws on not an array", func(t *testing.T) {
		err := jsonpatch.ValidateOperations(nil, false)
		assert.ErrorIs(t, err, jsonpatch.ErrNotArray)
	})

	t.Run("throws on empty array", func(t *testing.T) {
		err := jsonpatch.ValidateOperations([]jsonpatch.Operation{}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrEmptyPatch)
	})

	t.Run("throws on no operation path", func(t *testing.T) {
		ops := []jsonpatch.Operation{{}}
		err := jsonpatch.ValidateOperations(ops, false)
		require.Error(t, err)
		assert.ErrorIs(t, err, jsonpatch.ErrMissingOp)
	})

	t.Run("throws on no operation code", func(t *testing.T) {
		ops := []jsonpatch.Operation{{Path: "/"}}
		err := jsonpatch.ValidateOperations(ops, false)
		require.Error(t, err)
		assert.ErrorIs(t, err, jsonpatch.ErrMissingOp)
	})

	t.Run("throws on invalid operation code", func(t *testing.T) {
		ops := []jsonpatch.Operation{{Path: "/", Op: "123"}}
		err := jsonpatch.ValidateOperations(ops, false)
		require.Error(t, err)
		assert.ErrorIs(t, err, jsonpatch.ErrInvalidOperation)
	})

	t.Run("succeeds on valid operation", func(t *testing.T) {
		ops := []jsonpatch.Operation{{Op: "add", Path: "/test", Value: 123}}
		err := jsonpatch.ValidateOperations(ops, false)
		assert.NoError(t, err)
	})

	t.Run("throws on second invalid operation", func(t *testing.T) {
		ops := []jsonpatch.Operation{
			{Op: "add", Path: "/test", Value: 123},
			{Op: "test", Path: "/test"},
		}
		err := jsonpatch.ValidateOperations(ops, false)
		require.Error(t, err)
		assert.ErrorIs(t, err, jsonpatch.ErrMissingValue)
	})

	t.Run("throws if JSON pointer does not start with forward slash", func(t *testing.T) {
		ops := []jsonpatch.Operation{
			{Op: "add", Path: "/test", Value: 123},
			{Op: "test", Path: "test", Value: 1},
		}
		err := jsonpatch.ValidateOperations(ops, false)
		require.Error(t, err)
		assert.ErrorIs(t, err, jsonpatch.ErrInvalidJSONPointer)
	})
}

// =============================================================================
// ADD OPERATION TESTS
// =============================================================================

func TestValidateAdd(t *testing.T) {
	t.Run("throws with no path", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "add"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrMissingPath)
	})

	t.Run("throws with missing value", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "add", Path: "/"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrMissingValue)
	})

	t.Run("throws with null value", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "add", Path: "/", Value: nil}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrMissingValue)
	})

	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "add", Path: "/", Value: 123}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// REMOVE OPERATION TESTS
// =============================================================================

func TestValidateRemove(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "remove", Path: "/"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid path", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "remove", Path: "asdf"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrInvalidJSONPointer)
	})
}

// =============================================================================
// REPLACE OPERATION TESTS
// =============================================================================

func TestValidateReplace(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "replace", Path: "/", Value: "test", OldValue: "test"}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// COPY OPERATION TESTS
// =============================================================================

func TestValidateCopy(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "copy", From: "/", Path: "/"}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// MOVE OPERATION TESTS
// =============================================================================

func TestValidateMove(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "move", From: "/", Path: "/foo/bar"}, false)
		assert.NoError(t, err)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "move", From: "/foo/bar", Path: "/foo"}, false)
		assert.NoError(t, err)
	})

	t.Run("cannot move into its own children", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "move", From: "/foo", Path: "/foo/bar"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrCannotMoveToChildren)
	})
}

// =============================================================================
// TEST OPERATION TESTS
// =============================================================================

func TestValidateTest(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "test", Path: "/foo/bar", Value: "test"}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// TEST_EXISTS (DEFINED) OPERATION TESTS
// =============================================================================

func TestValidateTestExists(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "defined", Path: "/"},
			{Op: "defined", Path: "/"},
			{Op: "defined", Path: "/foo/bar"},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// TEST_TYPE OPERATION TESTS
// =============================================================================

func TestValidateTestType(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "test_type", Path: "/foo", Type: "number"},
			{Op: "test_type", Path: "/foo", Type: "array"},
			{Op: "test_type", Path: "/foo", Type: "string"},
			{Op: "test_type", Path: "/foo", Type: "boolean"},
			{Op: "test_type", Path: "/foo", Type: "integer"},
			{Op: "test_type", Path: "/foo", Type: "null"},
			{Op: "test_type", Path: "/foo", Type: "object"},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})

	t.Run("throws on no types provided", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "test_type", Path: "/foo", Type: ""}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrInvalidTypeField)
	})

	t.Run("throws on invalid type", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "test_type", Path: "/foo", Type: "monkey"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrInvalidType)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "test_type", Path: "/foo", Type: "invalid"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrInvalidType)
	})
}

// =============================================================================
// TEST_STRING OPERATION TESTS
// =============================================================================

func TestValidateTestString(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "test_string", Path: "/foo", Pos: 0, Str: "test", Not: true},
			{Op: "test_string", Path: "/foo", Pos: 0, Str: "test", Not: false},
			{Op: "test_string", Path: "/foo", Pos: 0, Str: "test"},
			{Op: "test_string", Path: "/foo", Pos: 0, Str: "", Not: true},
			{Op: "test_string", Path: "/foo", Pos: 0, Str: "", Not: false},
			{Op: "test_string", Path: "/foo", Pos: 0, Str: ""},
			{Op: "test_string", Path: "/foo", Pos: 123, Str: "test", Not: true},
			{Op: "test_string", Path: "/foo", Pos: 123, Str: "test", Not: false},
			{Op: "test_string", Path: "/foo", Pos: 123, Str: "test"},
			{Op: "test_string", Path: "/foo", Pos: 123, Str: "", Not: true},
			{Op: "test_string", Path: "/foo", Pos: 123, Str: "", Not: false},
			{Op: "test_string", Path: "/foo", Pos: 123, Str: ""},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// TEST_STRING_LEN OPERATION TESTS
// =============================================================================

func TestValidateTestStringLen(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "test_string_len", Path: "/foo", Len: 1, Not: false},
			{Op: "test_string_len", Path: "/foo", Len: 0, Not: false},
			{Op: "test_string_len", Path: "/foo", Len: 1, Not: true},
			{Op: "test_string_len", Path: "/foo", Len: 0, Not: true},
			{Op: "test_string_len", Path: "/foo", Len: 1},
			{Op: "test_string_len", Path: "/foo", Len: 0},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// FLIP OPERATION TESTS
// =============================================================================

func TestValidateFlip(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "flip", Path: "/"},
			{Op: "flip", Path: "/"},
			{Op: "flip", Path: "/foo"},
			{Op: "flip", Path: "/foo/bar"},
			{Op: "flip", Path: "/foo/123/bar"},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// INC OPERATION TESTS
// =============================================================================

func TestValidateInc(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "inc", Path: "/foo/bar", Inc: 0.0},
			{Op: "inc", Path: "/foo/bar", Inc: 0},
			{Op: "inc", Path: "/foo/bar", Inc: 1},
			{Op: "inc", Path: "/foo/bar", Inc: 1.5},
			{Op: "inc", Path: "/foo/bar", Inc: -1},
			{Op: "inc", Path: "/foo/bar", Inc: -1.5},
			{Op: "inc", Path: "/", Inc: 0},
			{Op: "inc", Path: "/", Inc: 1},
			{Op: "inc", Path: "/", Inc: 1.5},
			{Op: "inc", Path: "/", Inc: -1},
			{Op: "inc", Path: "/", Inc: -1.5},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// STR_INS OPERATION TESTS
// =============================================================================

func TestValidateStrIns(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "str_ins", Path: "/foo/bar", Pos: 0, Str: ""},
			{Op: "str_ins", Path: "/foo/bar", Pos: 0, Str: "test"},
			{Op: "str_ins", Path: "/foo/bar", Pos: 1, Str: ""},
			{Op: "str_ins", Path: "/foo/bar", Pos: 1, Str: "test"},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// STR_DEL OPERATION TESTS
// =============================================================================

func TestValidateStrDel(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "str_del", Path: "/foo/bar", Pos: 0, Str: ""},
			{Op: "str_del", Path: "/foo/bar", Pos: 0, Str: "test"},
			{Op: "str_del", Path: "/foo/bar", Pos: 0, Len: 0},
			{Op: "str_del", Path: "/foo/bar", Pos: 0, Len: 4},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// EXTEND OPERATION TESTS
// =============================================================================

func TestValidateExtend(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "extend", Path: "/foo/bar", Props: map[string]interface{}{}, DeleteNull: true},
			{Op: "extend", Path: "/foo/bar", Props: map[string]interface{}{}, DeleteNull: false},
			{Op: "extend", Path: "/foo/bar", Props: map[string]interface{}{}},
			{Op: "extend", Path: "/foo/bar", Props: map[string]interface{}{"foo": "bar"}, DeleteNull: true},
			{Op: "extend", Path: "/foo/bar", Props: map[string]interface{}{"foo": "bar"}, DeleteNull: false},
			{Op: "extend", Path: "/foo/bar", Props: map[string]interface{}{"foo": "bar"}},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// MERGE OPERATION TESTS
// =============================================================================

func TestValidateMerge(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "merge", Path: "/foo/bar", Pos: 1, Props: map[string]interface{}{}},
			{Op: "merge", Path: "/foo/bar", Pos: 2, Props: map[string]interface{}{}},
			{Op: "merge", Path: "/foo/bar", Pos: 1, Props: map[string]interface{}{"foo": "bar"}},
			{Op: "merge", Path: "/foo/bar", Pos: 2, Props: map[string]interface{}{"foo": "bar"}},
			{Op: "merge", Path: "/foo/bar", Pos: 1, Props: map[string]interface{}{"foo": nil}},
			{Op: "merge", Path: "/foo/bar", Pos: 2, Props: map[string]interface{}{"foo": nil}},
			{Op: "merge", Path: "/foo/bar", Pos: 1},
			{Op: "merge", Path: "/foo/bar", Pos: 2},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// SPLIT OPERATION TESTS
// =============================================================================

func TestValidateSplit(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{Op: "split", Path: "/foo/bar", Pos: 0},
			{Op: "split", Path: "/foo/bar", Pos: 2},
			{Op: "split", Path: "/foo/bar", Pos: 0, Props: map[string]interface{}{}},
			{Op: "split", Path: "/foo/bar", Pos: 2, Props: map[string]interface{}{}},
			{Op: "split", Path: "/foo/bar", Pos: 0, Props: map[string]interface{}{"foo": "bar"}},
			{Op: "split", Path: "/foo/bar", Pos: 2, Props: map[string]interface{}{"foo": "bar"}},
			{Op: "split", Path: "/foo/bar", Pos: 2, Props: map[string]interface{}{"foo": nil}},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})
}

// =============================================================================
// CONTAINS OPERATION TESTS
// =============================================================================

func TestValidateContains(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "contains", Path: "/foo/bar", Value: "test"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "contains", Path: "/foo/bar", Value: 123}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrExpectedValueToBeString)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "contains", Path: "/foo/bar", Value: "test", IgnoreCase: true}, false)
		assert.NoError(t, err) // IgnoreCase is a boolean field, should be ok
	})
}

// =============================================================================
// ENDS OPERATION TESTS
// =============================================================================

func TestValidateEnds(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "ends", Path: "/foo/bar", Value: "test"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "ends", Path: "/foo/bar", Value: 123}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrExpectedValueToBeString)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "ends", Path: "/foo/bar", Value: "test", IgnoreCase: true}, false)
		assert.NoError(t, err) // IgnoreCase is a boolean field, should be ok
	})
}

// =============================================================================
// STARTS OPERATION TESTS
// =============================================================================

func TestValidateStarts(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "starts", Path: "/foo/bar", Value: "test"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "starts", Path: "/foo/bar", Value: 123}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrExpectedValueToBeString)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "starts", Path: "/foo/bar", Value: "test", IgnoreCase: true}, false)
		assert.NoError(t, err) // IgnoreCase is a boolean field, should be ok
	})
}

// =============================================================================
// MATCHES OPERATION TESTS
// =============================================================================

func TestValidateMatches(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "matches", Path: "/foo/bar", Value: "test"}, true)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "matches", Path: "/foo/bar", Value: 123}, true)
		assert.ErrorIs(t, err, jsonpatch.ErrExpectedValueToBeString)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "matches", Path: "/foo/bar", Value: "test", IgnoreCase: true}, true)
		assert.NoError(t, err) // IgnoreCase is a boolean field, should be ok
	})
}

// =============================================================================
// DEFINED OPERATION TESTS
// =============================================================================

func TestValidateDefined(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "defined", Path: "/foo/bar"}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// UNDEFINED OPERATION TESTS
// =============================================================================

func TestValidateUndefined(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "undefined", Path: "/foo/bar"}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// IN OPERATION TESTS
// =============================================================================

func TestValidateIn(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "in", Path: "/foo/bar", Value: []interface{}{"test"}}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "in", Path: "/foo/bar", Value: 123}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrInOperationValueMustBeArray)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "in", Path: "/foo/bar", Value: "test"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrInOperationValueMustBeArray)
	})
}

// =============================================================================
// MORE OPERATION TESTS
// =============================================================================

func TestValidateMore(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "more", Path: "/foo/bar", Value: 5}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "more", Path: "/foo/bar", Value: "test"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrValueMustBeNumber)
	})
}

// =============================================================================
// LESS OPERATION TESTS
// =============================================================================

func TestValidateLess(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "less", Path: "/foo/bar", Value: 5}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "less", Path: "/foo/bar", Value: "test"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrValueMustBeNumber)
	})
}

// =============================================================================
// TYPE OPERATION TESTS
// =============================================================================

func TestValidateType(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "type", Path: "/foo/bar", Value: "number"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "type", Path: "/foo/bar", Value: 123}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrExpectedValueToBeString)
	})
}

// =============================================================================
// AND OPERATION TESTS
// =============================================================================

func TestValidateAnd(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{
			Op:   "and",
			Path: "/foo/bar",
			Apply: []jsonpatch.Operation{
				{Op: "test", Path: "/foo", Value: 123},
			},
		}, false)
		assert.NoError(t, err)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{
			Op:   "and",
			Path: "/foo/bar",
			Apply: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "/",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 123},
					},
				},
			},
		}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "and", Path: "/foo/bar", Apply: []jsonpatch.Operation{}}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrEmptyPredicateList)
	})
}

// =============================================================================
// NOT OPERATION TESTS
// =============================================================================

func TestValidateNot(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{
			Op:   "not",
			Path: "/foo/bar",
			Apply: []jsonpatch.Operation{
				{Op: "test", Path: "/foo", Value: 123},
			},
		}, false)
		assert.NoError(t, err)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{
			Op:   "not",
			Path: "/foo/bar",
			Apply: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "/",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 123},
					},
				},
			},
		}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "not", Path: "/foo/bar", Apply: []jsonpatch.Operation{}}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrEmptyPredicateList)
	})
}

// =============================================================================
// OR OPERATION TESTS
// =============================================================================

func TestValidateOr(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{
			Op:   "or",
			Path: "/foo/bar",
			Apply: []jsonpatch.Operation{
				{Op: "test", Path: "/foo", Value: 123},
			},
		}, false)
		assert.NoError(t, err)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{
			Op:   "or",
			Path: "/foo/bar",
			Apply: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "/",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 123},
					},
				},
			},
		}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "or", Path: "/foo/bar", Apply: []jsonpatch.Operation{}}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrEmptyPredicateList)
	})
}

// =============================================================================
// MATCHES OPERATION NOT ALLOWED TESTS
// =============================================================================

func TestValidateMatchesNotAllowed(t *testing.T) {
	t.Run("throws when matches operation not allowed", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "matches", Path: "/foo/bar", Value: "test"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrMatchesNotAllowed)
	})
}

// =============================================================================
// MERGE OPERATION ERROR TESTS
// =============================================================================

func TestValidateMergeErrors(t *testing.T) {
	t.Run("throws on missing pos", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "merge", Path: "/foo/bar"}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrPosGreaterThanZero)
	})

	t.Run("throws on pos less than 1", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{Op: "merge", Path: "/foo/bar", Pos: 0}, false)
		assert.ErrorIs(t, err, jsonpatch.ErrPosGreaterThanZero)
	})
}
