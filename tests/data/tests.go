package data

// TestCase is a single vector from the suite.
type TestCase struct {
	Comment  string           // Free‑text description
	Doc      any              // Original document (object / array / scalar)
	Patch    []map[string]any // JSON‑Patch as parsed maps
	Expected any              // Resulting document (nil when an error is expected)
	Error    string           // Sub‑string expected in the error (if any)
	Disabled bool             // Skip this test if true (present for completeness)
}

// TestCases is the full suite.
var TestCases = []TestCase{
	{
		Comment:  "empty list, empty docs",
		Doc:      map[string]any{},
		Patch:    []map[string]any{},
		Expected: map[string]any{},
	},
	{
		Comment:  "empty patch list",
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{},
		Expected: map[string]any{"foo": 1},
	},
	{
		Comment:  "rearrangements OK?",
		Doc:      map[string]any{"foo": 1, "bar": 2},
		Patch:    []map[string]any{},
		Expected: map[string]any{"bar": 2, "foo": 1},
	},
	{
		Comment: "rearrangements OK?  How about one level down ... array",
		Doc: []any{
			map[string]any{"foo": 1, "bar": 2},
		},
		Patch: []map[string]any{},
		Expected: []any{
			map[string]any{"bar": 2, "foo": 1},
		},
	},
	{
		Comment: "rearrangements OK?  How about one level down...",
		Doc: map[string]any{
			"foo": map[string]any{"foo": 1, "bar": 2},
		},
		Patch: []map[string]any{},
		Expected: map[string]any{
			"foo": map[string]any{"bar": 2, "foo": 1},
		},
	},
	{
		Comment:  "add replaces any existing field",
		Doc:      map[string]any{"foo": nil},
		Patch:    []map[string]any{{"op": "add", "path": "/foo", "value": 1}},
		Expected: map[string]any{"foo": 1},
	},
	{
		Comment:  "toplevel array",
		Doc:      []any{},
		Patch:    []map[string]any{{"op": "add", "path": "/0", "value": "foo"}},
		Expected: []any{"foo"},
	},
	{
		Comment:  "toplevel array, no change",
		Doc:      []any{"foo"},
		Patch:    []map[string]any{},
		Expected: []any{"foo"},
	},
	{
		Comment:  "toplevel object, numeric string",
		Doc:      map[string]any{},
		Patch:    []map[string]any{{"op": "add", "path": "/foo", "value": "1"}},
		Expected: map[string]any{"foo": "1"},
	},
	{
		Comment:  "toplevel object, integer",
		Doc:      map[string]any{},
		Patch:    []map[string]any{{"op": "add", "path": "/foo", "value": 1}},
		Expected: map[string]any{"foo": 1},
	},
	{
		Comment:  "Toplevel scalar values OK?",
		Doc:      "foo",
		Patch:    []map[string]any{{"op": "replace", "path": "", "value": "bar"}},
		Expected: "bar",
		Disabled: true,
	},
	{
		Comment:  "Add, / target",
		Doc:      map[string]any{},
		Patch:    []map[string]any{{"op": "add", "path": "/", "value": 1}},
		Expected: map[string]any{"": 1},
	},
	{
		Comment:  "Add composite value at top level",
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{{"op": "add", "path": "/bar", "value": []any{1, 2}}},
		Expected: map[string]any{"foo": 1, "bar": []any{1, 2}},
	},
	{
		Comment: "Add into composite value",
		Doc: map[string]any{
			"foo": 1,
			"baz": []any{map[string]any{"qux": "hello"}},
		},
		Patch: []map[string]any{{"op": "add", "path": "/baz/0/foo", "value": "world"}},
		Expected: map[string]any{
			"foo": 1,
			"baz": []any{map[string]any{"qux": "hello", "foo": "world"}},
		},
	},
	{
		Doc:   map[string]any{"bar": []any{1, 2}},
		Patch: []map[string]any{{"op": "add", "path": "/bar/8", "value": "5"}},
		Error: "Out of bounds (upper)",
	},
	{
		Doc:   map[string]any{"bar": []any{1, 2}},
		Patch: []map[string]any{{"op": "add", "path": "/bar/-1", "value": "5"}},
		Error: "Out of bounds (lower)",
	},
	{
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{{"op": "add", "path": "/bar", "value": true}},
		Expected: map[string]any{"foo": 1, "bar": true},
	},
	{
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{{"op": "add", "path": "/bar", "value": false}},
		Expected: map[string]any{"foo": 1, "bar": false},
	},
	{
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{{"op": "add", "path": "/bar", "value": nil}},
		Expected: map[string]any{"foo": 1, "bar": nil},
	},
	{
		Comment:  "0 can be an array index or object element name",
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{{"op": "add", "path": "/0", "value": "bar"}},
		Expected: map[string]any{"foo": 1, "0": "bar"},
	},
	{
		Doc:      []any{"foo"},
		Patch:    []map[string]any{{"op": "add", "path": "/1", "value": "bar"}},
		Expected: []any{"foo", "bar"},
	},
	{
		Doc:      []any{"foo", "sil"},
		Patch:    []map[string]any{{"op": "add", "path": "/1", "value": "bar"}},
		Expected: []any{"foo", "bar", "sil"},
	},
	{
		Doc:      []any{"foo", "sil"},
		Patch:    []map[string]any{{"op": "add", "path": "/0", "value": "bar"}},
		Expected: []any{"bar", "foo", "sil"},
	},
	{
		Doc:      []any{"foo", "sil"},
		Patch:    []map[string]any{{"op": "add", "path": "/2", "value": "bar"}},
		Expected: []any{"foo", "sil", "bar"},
	},
	{
		Comment:  "test against implementation-specific numeric parsing",
		Doc:      map[string]any{"1e0": "foo"},
		Patch:    []map[string]any{{"op": "test", "path": "/1e0", "value": "foo"}},
		Expected: map[string]any{"1e0": "foo"},
	},
	{
		Comment: "test with bad number should fail",
		Doc:     []any{"foo", "bar"},
		Patch:   []map[string]any{{"op": "test", "path": "/1e0", "value": "bar"}},
		Error:   "test op shouldn't get array element 1",
	},
	{
		Doc:   []any{"foo", "sil"},
		Patch: []map[string]any{{"op": "add", "path": "/bar", "value": 42}},
		Error: "Object operation on array target",
	},
	{
		Doc:      []any{"foo", "sil"},
		Patch:    []map[string]any{{"op": "add", "path": "/1", "value": []any{"bar", "baz"}}},
		Expected: []any{"foo", []any{"bar", "baz"}, "sil"},
		Comment:  "value in array add not flattened",
	},
	{
		Doc:      map[string]any{"foo": 1, "bar": []any{1, 2, 3, 4}},
		Patch:    []map[string]any{{"op": "remove", "path": "/bar"}},
		Expected: map[string]any{"foo": 1},
	},
	{
		Doc:      map[string]any{"foo": 1, "baz": []any{map[string]any{"qux": "hello"}}},
		Patch:    []map[string]any{{"op": "remove", "path": "/baz/0/qux"}},
		Expected: map[string]any{"foo": 1, "baz": []any{map[string]any{}}},
	},
	{
		Doc:      map[string]any{"foo": 1, "baz": []any{map[string]any{"qux": "hello"}}},
		Patch:    []map[string]any{{"op": "replace", "path": "/foo", "value": []any{1, 2, 3, 4}}},
		Expected: map[string]any{"foo": []any{1, 2, 3, 4}, "baz": []any{map[string]any{"qux": "hello"}}},
	},
	{
		Doc:      map[string]any{"foo": []any{1, 2, 3, 4}, "baz": []any{map[string]any{"qux": "hello"}}},
		Patch:    []map[string]any{{"op": "replace", "path": "/baz/0/qux", "value": "world"}},
		Expected: map[string]any{"foo": []any{1, 2, 3, 4}, "baz": []any{map[string]any{"qux": "world"}}},
	},
	{
		Doc:      []any{"foo"},
		Patch:    []map[string]any{{"op": "replace", "path": "/0", "value": "bar"}},
		Expected: []any{"bar"},
	},
	{
		Doc:      []any{""},
		Patch:    []map[string]any{{"op": "replace", "path": "/0", "value": 0}},
		Expected: []any{0},
	},
	{
		Doc:      []any{""},
		Patch:    []map[string]any{{"op": "replace", "path": "/0", "value": true}},
		Expected: []any{true},
	},
	{
		Doc:      []any{""},
		Patch:    []map[string]any{{"op": "replace", "path": "/0", "value": false}},
		Expected: []any{false},
	},
	{
		Doc:      []any{""},
		Patch:    []map[string]any{{"op": "replace", "path": "/0", "value": nil}},
		Expected: []any{nil},
	},
	{
		Doc:      []any{"foo", "sil"},
		Patch:    []map[string]any{{"op": "replace", "path": "/1", "value": []any{"bar", "baz"}}},
		Expected: []any{"foo", []any{"bar", "baz"}},
		Comment:  "value in array replace not flattened",
	},
	{
		Comment:  "replace whole document",
		Doc:      map[string]any{"foo": "bar"},
		Patch:    []map[string]any{{"op": "replace", "path": "", "value": map[string]any{"baz": "qux"}}},
		Expected: map[string]any{"baz": "qux"},
	},
	{
		Comment:  "spurious patch properties",
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{{"op": "test", "path": "/foo", "value": 1, "spurious": 1}},
		Expected: map[string]any{"foo": 1},
	},
	{
		Doc:      map[string]any{"foo": nil},
		Patch:    []map[string]any{{"op": "test", "path": "/foo", "value": nil}},
		Comment:  "null value should be valid obj property",
		Expected: map[string]any{"foo": nil},
	},
	{
		Doc:      map[string]any{"foo": nil},
		Patch:    []map[string]any{{"op": "replace", "path": "/foo", "value": "truthy"}},
		Expected: map[string]any{"foo": "truthy"},
		Comment:  "null value should be valid obj property to be replaced with something truthy",
	},
	{
		Doc:      map[string]any{"foo": nil},
		Patch:    []map[string]any{{"op": "move", "from": "/foo", "path": "/bar"}},
		Expected: map[string]any{"bar": nil},
		Comment:  "null value should be valid obj property to be moved",
	},
	{
		Doc:      map[string]any{"foo": nil},
		Patch:    []map[string]any{{"op": "copy", "from": "/foo", "path": "/bar"}},
		Expected: map[string]any{"foo": nil, "bar": nil},
		Comment:  "null value should be valid obj property to be copied",
	},
	{
		Doc:      map[string]any{"foo": nil},
		Patch:    []map[string]any{{"op": "remove", "path": "/foo"}},
		Expected: map[string]any{},
		Comment:  "null value should be valid obj property to be removed",
	},
	{
		Doc:      map[string]any{"foo": "bar"},
		Patch:    []map[string]any{{"op": "replace", "path": "/foo", "value": nil}},
		Expected: map[string]any{"foo": nil},
		Comment:  "null value should still be valid obj property replace other value",
	},
	{
		Doc:      map[string]any{"foo": map[string]any{"foo": 1, "bar": 2}},
		Patch:    []map[string]any{{"op": "test", "path": "/foo", "value": map[string]any{"bar": 2, "foo": 1}}},
		Comment:  "test should pass despite rearrangement",
		Expected: map[string]any{"foo": map[string]any{"foo": 1, "bar": 2}},
	},
	{
		Doc:      map[string]any{"foo": []any{map[string]any{"foo": 1, "bar": 2}}},
		Patch:    []map[string]any{{"op": "test", "path": "/foo", "value": []any{map[string]any{"bar": 2, "foo": 1}}}},
		Comment:  "test should pass despite (nested) rearrangement",
		Expected: map[string]any{"foo": []any{map[string]any{"foo": 1, "bar": 2}}},
	},
	{
		Doc:      map[string]any{"foo": map[string]any{"bar": []any{1, 2, 5, 4}}},
		Patch:    []map[string]any{{"op": "test", "path": "/foo", "value": map[string]any{"bar": []any{1, 2, 5, 4}}}},
		Comment:  "test should pass - no error",
		Expected: map[string]any{"foo": map[string]any{"bar": []any{1, 2, 5, 4}}},
	},
	{
		Doc:   map[string]any{"foo": map[string]any{"bar": []any{1, 2, 5, 4}}},
		Patch: []map[string]any{{"op": "test", "path": "/foo", "value": []any{1, 2}}},
		Error: "test op should fail",
	},
	{
		Comment:  "Whole document",
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{{"op": "test", "path": "", "value": map[string]any{"foo": 1}}},
		Disabled: true,
	},
	{
		Comment:  "Empty-string element",
		Doc:      map[string]any{"": 1},
		Expected: map[string]any{"": 1},
		Patch:    []map[string]any{{"op": "test", "path": "/", "value": 1}},
	},
	{
		Doc: map[string]any{
			"foo":  []any{"bar", "baz"},
			"":     0,
			"a/b":  1,
			"c%d":  2,
			"e^f":  3,
			"g|h":  4,
			"i\\j": 5,
			"k\"l": 6,
			" ":    7,
			"m~n":  8,
		},
		Expected: map[string]any{
			"foo":  []any{"bar", "baz"},
			"":     0,
			"a/b":  1,
			"c%d":  2,
			"e^f":  3,
			"g|h":  4,
			"i\\j": 5,
			"k\"l": 6,
			" ":    7,
			"m~n":  8,
		},
		Patch: []map[string]any{
			{"op": "test", "path": "/foo", "value": []any{"bar", "baz"}},
			{"op": "test", "path": "/foo/0", "value": "bar"},
			{"op": "test", "path": "/", "value": 0},
			{"op": "test", "path": "/a~1b", "value": 1},
			{"op": "test", "path": "/c%d", "value": 2},
			{"op": "test", "path": "/e^f", "value": 3},
			{"op": "test", "path": "/g|h", "value": 4},
			{"op": "test", "path": "/i\\j", "value": 5},
			{"op": "test", "path": "/k\"l", "value": 6},
			{"op": "test", "path": "/ ", "value": 7},
			{"op": "test", "path": "/m~0n", "value": 8},
		},
	},
	{
		Comment:  "Move to same location has no effect",
		Doc:      map[string]any{"foo": 1},
		Patch:    []map[string]any{{"op": "move", "from": "/foo", "path": "/foo"}},
		Expected: map[string]any{"foo": 1},
	},
	{
		Doc:      map[string]any{"foo": 1, "baz": []any{map[string]any{"qux": "hello"}}},
		Patch:    []map[string]any{{"op": "move", "from": "/foo", "path": "/bar"}},
		Expected: map[string]any{"baz": []any{map[string]any{"qux": "hello"}}, "bar": 1},
	},
	{
		Doc:      map[string]any{"baz": []any{map[string]any{"qux": "hello"}}, "bar": 1},
		Patch:    []map[string]any{{"op": "move", "from": "/baz/0/qux", "path": "/baz/1"}},
		Expected: map[string]any{"baz": []any{map[string]any{}, "hello"}, "bar": 1},
	},
	{
		Doc:      map[string]any{"baz": []any{map[string]any{"qux": "hello"}}, "bar": 1},
		Patch:    []map[string]any{{"op": "copy", "from": "/baz/0", "path": "/boo"}},
		Expected: map[string]any{"baz": []any{map[string]any{"qux": "hello"}}, "bar": 1, "boo": map[string]any{"qux": "hello"}},
	},
	{
		Comment:  "replacing the root of the document is possible with add",
		Doc:      map[string]any{"foo": "bar"},
		Patch:    []map[string]any{{"op": "add", "path": "", "value": map[string]any{"baz": "qux"}}},
		Expected: map[string]any{"baz": "qux"},
	},
	{
		Comment:  "Adding to \"/-\" adds to the end of the array",
		Doc:      []any{1, 2},
		Patch:    []map[string]any{{"op": "add", "path": "/-", "value": map[string]any{"foo": []any{"bar", "baz"}}}},
		Expected: []any{1, 2, map[string]any{"foo": []any{"bar", "baz"}}},
	},
	{
		Comment:  "Adding to \"/-\" adds to the end of the array, even n levels down",
		Doc:      []any{1, 2, []any{3, []any{4, 5}}},
		Patch:    []map[string]any{{"op": "add", "path": "/2/1/-", "value": map[string]any{"foo": []any{"bar", "baz"}}}},
		Expected: []any{1, 2, []any{3, []any{4, 5, map[string]any{"foo": []any{"bar", "baz"}}}}},
	},
	{
		Comment: "test remove with bad number should fail",
		Doc:     map[string]any{"foo": 1, "baz": []any{map[string]any{"qux": "hello"}}},
		Patch:   []map[string]any{{"op": "remove", "path": "/baz/1e0/qux"}},
		Error:   "remove op shouldn't remove from array with bad number",
	},
	{
		Comment:  "test remove on array",
		Doc:      []any{1, 2, 3, 4},
		Patch:    []map[string]any{{"op": "remove", "path": "/0"}},
		Expected: []any{2, 3, 4},
	},
	{
		Comment:  "test repeated removes",
		Doc:      []any{1, 2, 3, 4},
		Patch:    []map[string]any{{"op": "remove", "path": "/1"}, {"op": "remove", "path": "/2"}},
		Expected: []any{1, 3},
	},
	{
		Comment: "test remove with bad index should fail",
		Doc:     []any{1, 2, 3, 4},
		Patch:   []map[string]any{{"op": "remove", "path": "/1e0"}},
		Error:   "remove op shouldn't remove from array with bad number",
	},
	{
		Comment: "test replace with bad number should fail",
		Doc:     []any{""},
		Patch:   []map[string]any{{"op": "replace", "path": "/1e0", "value": false}},
		Error:   "replace op shouldn't replace in array with bad number",
	},
	{
		Comment: "test copy with bad number should fail",
		Doc:     map[string]any{"baz": []any{1, 2, 3}, "bar": 1},
		Patch:   []map[string]any{{"op": "copy", "from": "/baz/1e0", "path": "/boo"}},
		Error:   "copy op shouldn't work with bad number",
	},
	{
		Comment: "test move with bad number should fail",
		Doc:     map[string]any{"foo": 1, "baz": []any{1, 2, 3, 4}},
		Patch:   []map[string]any{{"op": "move", "from": "/baz/1e0", "path": "/foo"}},
		Error:   "move op shouldn't work with bad number",
	},
	{
		Comment: "test add with bad number should fail",
		Doc:     []any{"foo", "sil"},
		Patch:   []map[string]any{{"op": "add", "path": "/1e0", "value": "bar"}},
		Error:   "add op shouldn't add to array with bad number",
	},
	{
		Comment: "missing 'value' parameter to test",
		Doc:     []any{nil},
		Patch:   []map[string]any{{"op": "test", "path": "/0"}},
		Error:   "missing 'value' parameter",
	},
	{
		Comment: "missing value parameter to test - where undef is falsy",
		Doc:     []any{false},
		Patch:   []map[string]any{{"op": "test", "path": "/0"}},
		Error:   "missing 'value' parameter",
	},
	{
		Comment: "missing from parameter to copy",
		Doc:     []any{1},
		Patch:   []map[string]any{{"op": "copy", "path": "/-"}},
		Error:   "missing 'from' parameter",
	},
	{
		Comment: "unrecognized op should fail",
		Doc:     map[string]any{"foo": 1},
		Patch:   []map[string]any{{"op": "spam", "path": "/foo", "value": 1}},
		Error:   "Unrecognized op 'spam'",
	},
}
