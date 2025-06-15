package jsonpatch_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// OPERATIONS TESTS
// =============================================================================

func TestValidateOperations(t *testing.T) {
	t.Run("throws on not an array", func(t *testing.T) {
		err := jsonpatch.ValidateOperations(nil, false)
		assert.EqualError(t, err, "not an array")
	})

	t.Run("throws on empty array", func(t *testing.T) {
		err := jsonpatch.ValidateOperations([]jsonpatch.Operation{}, false)
		assert.EqualError(t, err, "empty operation patch")
	})

	t.Run("throws on invalid operation type", func(t *testing.T) {
		ops := []jsonpatch.Operation{
			{"op": 123, "path": "/test"},
		}
		err := jsonpatch.ValidateOperations(ops, false)
		assert.Contains(t, err.Error(), "error in operation [index = 0]")
		assert.Contains(t, err.Error(), "field 'op' must be a string")
	})

	t.Run("throws on no operation path", func(t *testing.T) {
		ops := []jsonpatch.Operation{{}}
		err := jsonpatch.ValidateOperations(ops, false)
		assert.Contains(t, err.Error(), "error in operation [index = 0]")
		assert.Contains(t, err.Error(), "missing required field 'path'")
	})

	t.Run("throws on no operation code", func(t *testing.T) {
		ops := []jsonpatch.Operation{{"path": ""}}
		err := jsonpatch.ValidateOperations(ops, false)
		assert.Contains(t, err.Error(), "error in operation [index = 0]")
		assert.Contains(t, err.Error(), "missing required field 'op'")
	})

	t.Run("throws on invalid operation code", func(t *testing.T) {
		ops := []jsonpatch.Operation{{"path": "", "op": "123"}}
		err := jsonpatch.ValidateOperations(ops, false)
		assert.Contains(t, err.Error(), "error in operation [index = 0]")
		assert.Contains(t, err.Error(), "unknown operation '123'")
	})

	t.Run("succeeds on valid operation", func(t *testing.T) {
		ops := []jsonpatch.Operation{{"op": "add", "path": "/test", "value": 123}}
		err := jsonpatch.ValidateOperations(ops, false)
		assert.NoError(t, err)
	})

	t.Run("throws on second invalid operation", func(t *testing.T) {
		ops := []jsonpatch.Operation{
			{"op": "add", "path": "/test", "value": 123},
			{"op": "test", "path": "/test"},
		}
		err := jsonpatch.ValidateOperations(ops, false)
		assert.Contains(t, err.Error(), "error in operation [index = 1]")
		assert.Contains(t, err.Error(), "missing required field 'value'")
	})

	t.Run("throws if JSON pointer does not start with forward slash", func(t *testing.T) {
		ops := []jsonpatch.Operation{
			{"op": "add", "path": "/test", "value": 123},
			{"op": "test", "path": "test", "value": 1},
		}
		err := jsonpatch.ValidateOperations(ops, false)
		assert.Contains(t, err.Error(), "error in operation [index = 1]")
		assert.Contains(t, err.Error(), "invalid JSON pointer")
	})
}

// =============================================================================
// ADD OPERATION TESTS
// =============================================================================

func TestValidateAdd(t *testing.T) {
	t.Run("throws with no path", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "add"}, false)
		assert.EqualError(t, err, "missing required field 'path'")
	})

	t.Run("throws with invalid path", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "add", "path": 123}, false)
		assert.EqualError(t, err, "field 'path' must be a string")
	})

	t.Run("throws with missing value", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "add", "path": ""}, false)
		assert.EqualError(t, err, "missing required field 'value'")
	})

	t.Run("succeeds with null value", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "add", "path": "", "value": nil}, false)
		assert.NoError(t, err)
	})

	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "add", "path": "", "value": 123}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// REMOVE OPERATION TESTS
// =============================================================================

func TestValidateRemove(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "remove", "path": ""}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid path", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "remove", "path": "asdf"}, false)
		assert.Contains(t, err.Error(), "invalid JSON pointer")
	})

	t.Run("throws on invalid path - 2", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "remove", "path": 123}, false)
		assert.EqualError(t, err, "field 'path' must be a string")
	})
}

// =============================================================================
// REPLACE OPERATION TESTS
// =============================================================================

func TestValidateReplace(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "replace", "path": "", "value": "test", "oldValue": "test"}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// COPY OPERATION TESTS
// =============================================================================

func TestValidateCopy(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "copy", "from": "", "path": ""}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// MOVE OPERATION TESTS
// =============================================================================

func TestValidateMove(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "move", "from": "/", "path": "/foo/bar"}, false)
		assert.NoError(t, err)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "move", "from": "/foo/bar", "path": "/foo"}, false)
		assert.NoError(t, err)
	})

	t.Run("cannot move into its own children", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "move", "from": "/foo", "path": "/foo/bar"}, false)
		assert.EqualError(t, err, "cannot move into own children")
	})
}

// =============================================================================
// TEST OPERATION TESTS
// =============================================================================

func TestValidateTest(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "test", "path": "/foo/bar", "value": nil}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// TEST_EXISTS (DEFINED) OPERATION TESTS
// =============================================================================

func TestValidateTestExists(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{"op": "defined", "path": ""},
			{"op": "defined", "path": "/"},
			{"op": "defined", "path": "/foo/bar"},
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
			{"op": "test_type", "path": "/foo", "type": []interface{}{"number"}},
			{"op": "test_type", "path": "/foo", "type": []interface{}{"number", "array"}},
			{"op": "test_type", "path": "/foo", "type": []interface{}{"number", "string"}},
			{"op": "test_type", "path": "/foo", "type": []interface{}{"number", "boolean"}},
			{"op": "test_type", "path": "/foo", "type": []interface{}{"number", "integer"}},
			{"op": "test_type", "path": "/foo", "type": []interface{}{"number", "null"}},
			{"op": "test_type", "path": "/foo", "type": []interface{}{"number", "object"}},
		}

		for _, test := range tests {
			err := jsonpatch.ValidateOperation(test, false)
			assert.NoError(t, err)
		}
	})

	t.Run("throws on no types provided", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "test_type", "path": "/foo", "type": []interface{}{}}, false)
		assert.EqualError(t, err, "empty type list")
	})

	t.Run("throws on invalid type", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "test_type", "path": "/foo", "type": []interface{}{"monkey"}}, false)
		assert.EqualError(t, err, "invalid type")

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "test_type", "path": "/foo", "type": []interface{}{"number", "monkey"}}, false)
		assert.EqualError(t, err, "invalid type")
	})
}

// =============================================================================
// TEST_STRING OPERATION TESTS
// =============================================================================

func TestValidateTestString(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		tests := []jsonpatch.Operation{
			{"op": "test_string", "path": "/foo", "pos": 0, "str": "test", "not": true},
			{"op": "test_string", "path": "/foo", "pos": 0, "str": "test", "not": false},
			{"op": "test_string", "path": "/foo", "pos": 0, "str": "test"},
			{"op": "test_string", "path": "/foo", "pos": 0, "str": "", "not": true},
			{"op": "test_string", "path": "/foo", "pos": 0, "str": "", "not": false},
			{"op": "test_string", "path": "/foo", "pos": 0, "str": ""},
			{"op": "test_string", "path": "/foo", "pos": 123, "str": "test", "not": true},
			{"op": "test_string", "path": "/foo", "pos": 123, "str": "test", "not": false},
			{"op": "test_string", "path": "/foo", "pos": 123, "str": "test"},
			{"op": "test_string", "path": "/foo", "pos": 123, "str": "", "not": true},
			{"op": "test_string", "path": "/foo", "pos": 123, "str": "", "not": false},
			{"op": "test_string", "path": "/foo", "pos": 123, "str": ""},
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
			{"op": "test_string_len", "path": "/foo", "len": 1, "not": false},
			{"op": "test_string_len", "path": "/foo", "len": 0, "not": false},
			{"op": "test_string_len", "path": "/foo", "len": 1, "not": true},
			{"op": "test_string_len", "path": "/foo", "len": 0, "not": true},
			{"op": "test_string_len", "path": "/foo", "len": 1},
			{"op": "test_string_len", "path": "/foo", "len": 0},
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
			{"op": "flip", "path": ""},
			{"op": "flip", "path": "/"},
			{"op": "flip", "path": "/foo"},
			{"op": "flip", "path": "/foo/bar"},
			{"op": "flip", "path": "/foo/123/bar"},
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
			{"op": "inc", "path": "/foo/bar", "inc": 0.0},
			{"op": "inc", "path": "/foo/bar", "inc": 0},
			{"op": "inc", "path": "/foo/bar", "inc": 1},
			{"op": "inc", "path": "/foo/bar", "inc": 1.5},
			{"op": "inc", "path": "/foo/bar", "inc": -1},
			{"op": "inc", "path": "/foo/bar", "inc": -1.5},
			{"op": "inc", "path": "", "inc": 0},
			{"op": "inc", "path": "", "inc": 1},
			{"op": "inc", "path": "", "inc": 1.5},
			{"op": "inc", "path": "", "inc": -1},
			{"op": "inc", "path": "", "inc": -1.5},
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
			{"op": "str_ins", "path": "/foo/bar", "pos": 0, "str": ""},
			{"op": "str_ins", "path": "/foo/bar", "pos": 0, "str": "test"},
			{"op": "str_ins", "path": "/foo/bar", "pos": 1, "str": ""},
			{"op": "str_ins", "path": "/foo/bar", "pos": 1, "str": "test"},
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
			{"op": "str_del", "path": "/foo/bar", "pos": 0, "str": ""},
			{"op": "str_del", "path": "/foo/bar", "pos": 0, "str": "test"},
			{"op": "str_del", "path": "/foo/bar", "pos": 0, "len": 0},
			{"op": "str_del", "path": "/foo/bar", "pos": 0, "len": 4},
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
			{"op": "extend", "path": "/foo/bar", "props": map[string]interface{}{}, "deleteNull": true},
			{"op": "extend", "path": "/foo/bar", "props": map[string]interface{}{}, "deleteNull": false},
			{"op": "extend", "path": "/foo/bar", "props": map[string]interface{}{}},
			{"op": "extend", "path": "/foo/bar", "props": map[string]interface{}{"foo": "bar"}, "deleteNull": true},
			{"op": "extend", "path": "/foo/bar", "props": map[string]interface{}{"foo": "bar"}, "deleteNull": false},
			{"op": "extend", "path": "/foo/bar", "props": map[string]interface{}{"foo": "bar"}},
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
			{"op": "merge", "path": "/foo/bar", "pos": 1, "props": map[string]interface{}{}},
			{"op": "merge", "path": "/foo/bar", "pos": 2, "props": map[string]interface{}{}},
			{"op": "merge", "path": "/foo/bar", "pos": 1, "props": map[string]interface{}{"foo": "bar"}},
			{"op": "merge", "path": "/foo/bar", "pos": 2, "props": map[string]interface{}{"foo": "bar"}},
			{"op": "merge", "path": "/foo/bar", "pos": 1, "props": map[string]interface{}{"foo": nil}},
			{"op": "merge", "path": "/foo/bar", "pos": 2, "props": map[string]interface{}{"foo": nil}},
			{"op": "merge", "path": "/foo/bar", "pos": 1},
			{"op": "merge", "path": "/foo/bar", "pos": 2},
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
			{"op": "split", "path": "/foo/bar", "pos": 0},
			{"op": "split", "path": "/foo/bar", "pos": 2},
			{"op": "split", "path": "/foo/bar", "pos": 0, "props": map[string]interface{}{}},
			{"op": "split", "path": "/foo/bar", "pos": 2, "props": map[string]interface{}{}},
			{"op": "split", "path": "/foo/bar", "pos": 0, "props": map[string]interface{}{"foo": "bar"}},
			{"op": "split", "path": "/foo/bar", "pos": 2, "props": map[string]interface{}{"foo": "bar"}},
			{"op": "split", "path": "/foo/bar", "pos": 2, "props": map[string]interface{}{"foo": nil}},
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
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "contains", "path": "/foo/bar", "value": "test"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "contains", "path": "/foo/bar", "value": 123}, false)
		assert.Contains(t, err.Error(), "expected value to be string")

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "contains", "path": "/foo/bar", "value": "test", "ignore_case": 1}, false)
		assert.Contains(t, err.Error(), "expected ignore_case to be boolean")
	})
}

// =============================================================================
// ENDS OPERATION TESTS
// =============================================================================

func TestValidateEnds(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "ends", "path": "/foo/bar", "value": "test"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "ends", "path": "/foo/bar", "value": 123}, false)
		assert.Contains(t, err.Error(), "expected value to be string")

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "ends", "path": "/foo/bar", "value": "test", "ignore_case": 1}, false)
		assert.Contains(t, err.Error(), "expected ignore_case to be boolean")
	})
}

// =============================================================================
// STARTS OPERATION TESTS
// =============================================================================

func TestValidateStarts(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "starts", "path": "/foo/bar", "value": "test"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "starts", "path": "/foo/bar", "value": 123}, false)
		assert.Contains(t, err.Error(), "expected value to be string")

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "starts", "path": "/foo/bar", "value": "test", "ignore_case": 1}, false)
		assert.Contains(t, err.Error(), "expected ignore_case to be boolean")
	})
}

// =============================================================================
// MATCHES OPERATION TESTS
// =============================================================================

func TestValidateMatches(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "matches", "path": "/foo/bar", "value": "test"}, true)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "matches", "path": "/foo/bar", "value": 123}, true)
		assert.Contains(t, err.Error(), "expected value to be string")

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "matches", "path": "/foo/bar", "value": "test", "ignore_case": 1}, true)
		assert.Contains(t, err.Error(), "expected ignore_case to be boolean")
	})
}

// =============================================================================
// DEFINED OPERATION TESTS
// =============================================================================

func TestValidateDefined(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "defined", "path": "/foo/bar"}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// UNDEFINED OPERATION TESTS
// =============================================================================

func TestValidateUndefined(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "undefined", "path": "/foo/bar"}, false)
		assert.NoError(t, err)
	})
}

// =============================================================================
// IN OPERATION TESTS
// =============================================================================

func TestValidateIn(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "in", "path": "/foo/bar", "value": []interface{}{"test"}}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "in", "path": "/foo/bar", "value": 123}, false)
		assert.Contains(t, err.Error(), "in operation value must be an array")

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "in", "path": "/foo/bar", "value": "test"}, false)
		assert.Contains(t, err.Error(), "in operation value must be an array")
	})
}

// =============================================================================
// MORE OPERATION TESTS
// =============================================================================

func TestValidateMore(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "more", "path": "/foo/bar", "value": 5}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "more", "path": "/foo/bar", "value": "test"}, false)
		assert.Contains(t, err.Error(), "value must be a number")
	})
}

// =============================================================================
// LESS OPERATION TESTS
// =============================================================================

func TestValidateLess(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "less", "path": "/foo/bar", "value": 5}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "less", "path": "/foo/bar", "value": "test"}, false)
		assert.Contains(t, err.Error(), "value must be a number")
	})
}

// =============================================================================
// TYPE OPERATION TESTS
// =============================================================================

func TestValidateType(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "type", "path": "/foo/bar", "value": "number"}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "type", "path": "/foo/bar", "value": 123}, false)
		assert.Contains(t, err.Error(), "expected value to be string")
	})
}

// =============================================================================
// AND OPERATION TESTS
// =============================================================================

func TestValidateAnd(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{
			"op":   "and",
			"path": "/foo/bar",
			"apply": []interface{}{
				map[string]interface{}{"op": "test", "path": "/foo", "value": 123},
			},
		}, false)
		assert.NoError(t, err)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{
			"op":   "and",
			"path": "/foo/bar",
			"apply": []interface{}{
				map[string]interface{}{
					"op":   "not",
					"path": "",
					"apply": []interface{}{
						map[string]interface{}{"op": "test", "path": "/foo", "value": 123},
					},
				},
			},
		}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "and", "path": "/foo/bar", "apply": []interface{}{}}, false)
		assert.Contains(t, err.Error(), "predicate list is empty")
	})
}

// =============================================================================
// NOT OPERATION TESTS
// =============================================================================

func TestValidateNot(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{
			"op":   "not",
			"path": "/foo/bar",
			"apply": []interface{}{
				map[string]interface{}{"op": "test", "path": "/foo", "value": 123},
			},
		}, false)
		assert.NoError(t, err)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{
			"op":   "not",
			"path": "/foo/bar",
			"apply": []interface{}{
				map[string]interface{}{
					"op":   "not",
					"path": "",
					"apply": []interface{}{
						map[string]interface{}{"op": "test", "path": "/foo", "value": 123},
					},
				},
			},
		}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "not", "path": "/foo/bar", "apply": []interface{}{}}, false)
		assert.Contains(t, err.Error(), "predicate list is empty")
	})
}

// =============================================================================
// OR OPERATION TESTS
// =============================================================================

func TestValidateOr(t *testing.T) {
	t.Run("succeeds on valid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{
			"op":   "or",
			"path": "/foo/bar",
			"apply": []interface{}{
				map[string]interface{}{"op": "test", "path": "/foo", "value": 123},
			},
		}, false)
		assert.NoError(t, err)

		err = jsonpatch.ValidateOperation(jsonpatch.Operation{
			"op":   "or",
			"path": "/foo/bar",
			"apply": []interface{}{
				map[string]interface{}{
					"op":   "not",
					"path": "",
					"apply": []interface{}{
						map[string]interface{}{"op": "test", "path": "/foo", "value": 123},
					},
				},
			},
		}, false)
		assert.NoError(t, err)
	})

	t.Run("throws on invalid operation", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "or", "path": "/foo/bar", "apply": []interface{}{}}, false)
		assert.Contains(t, err.Error(), "predicate list is empty")
	})
}

// =============================================================================
// MATCHES OPERATION NOT ALLOWED TESTS
// =============================================================================

func TestValidateMatchesNotAllowed(t *testing.T) {
	t.Run("throws when matches operation not allowed", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "matches", "path": "/foo/bar", "value": "test"}, false)
		assert.Contains(t, err.Error(), "matches operation not allowed")
	})
}

// =============================================================================
// MERGE OPERATION ERROR TESTS
// =============================================================================

func TestValidateMergeErrors(t *testing.T) {
	t.Run("throws on missing pos", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "merge", "path": "/foo/bar"}, false)
		assert.Contains(t, err.Error(), "expected pos field to be greater than 0")
	})

	t.Run("throws on invalid pos", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "merge", "path": "/foo/bar", "pos": "invalid"}, false)
		assert.Contains(t, err.Error(), "not an integer")
	})

	t.Run("throws on pos less than 1", func(t *testing.T) {
		err := jsonpatch.ValidateOperation(jsonpatch.Operation{"op": "merge", "path": "/foo/bar", "pos": 0}, false)
		assert.Contains(t, err.Error(), "expected pos field to be greater than 0")
	})
}
