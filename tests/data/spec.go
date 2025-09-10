// Package data contains test case data and specifications for JSON Patch operations.
package data

// SpecTestCases contains test cases for JSON Patch operations.
var SpecTestCases = []TestCase{
	{
		Comment: "4.1. add with missing object",
		Doc:     map[string]any{"q": map[string]any{"bar": 2}},
		Patch:   []map[string]any{{"op": "add", "path": "/a/b", "value": 1}},
		Error:   "path /a does not exist -- missing objects are not created recursively",
	},
	{
		Comment:  "A.1.  Adding an Object Member",
		Doc:      map[string]any{"foo": "bar"},
		Patch:    []map[string]any{{"op": "add", "path": "/baz", "value": "qux"}},
		Expected: map[string]any{"foo": "bar", "baz": "qux"},
	},
	{
		Comment:  "A.2.  Adding an Array Element",
		Doc:      map[string]any{"foo": []any{"bar", "baz"}},
		Patch:    []map[string]any{{"op": "add", "path": "/foo/1", "value": "qux"}},
		Expected: map[string]any{"foo": []any{"bar", "qux", "baz"}},
	},
	{
		Comment:  "A.3.  Removing an Object Member",
		Doc:      map[string]any{"baz": "qux", "foo": "bar"},
		Patch:    []map[string]any{{"op": "remove", "path": "/baz"}},
		Expected: map[string]any{"foo": "bar"},
	},
	{
		Comment:  "A.4.  Removing an Array Element",
		Doc:      map[string]any{"foo": []any{"bar", "qux", "baz"}},
		Patch:    []map[string]any{{"op": "remove", "path": "/foo/1"}},
		Expected: map[string]any{"foo": []any{"bar", "baz"}},
	},
	{
		Comment:  "A.5.  Replacing a Value",
		Doc:      map[string]any{"baz": "qux", "foo": "bar"},
		Patch:    []map[string]any{{"op": "replace", "path": "/baz", "value": "boo"}},
		Expected: map[string]any{"baz": "boo", "foo": "bar"},
	},
	{
		Comment: "A.6.  Moving a Value",
		Doc: map[string]any{
			"foo": map[string]any{
				"bar":   "baz",
				"waldo": "fred",
			},
			"qux": map[string]any{
				"corge": "grault",
			},
		},
		Patch: []map[string]any{{"op": "move", "from": "/foo/waldo", "path": "/qux/thud"}},
		Expected: map[string]any{
			"foo": map[string]any{
				"bar": "baz",
			},
			"qux": map[string]any{
				"corge": "grault",
				"thud":  "fred",
			},
		},
	},
	{
		Comment:  "A.7.  Moving an Array Element",
		Doc:      map[string]any{"foo": []any{"all", "grass", "cows", "eat"}},
		Patch:    []map[string]any{{"op": "move", "from": "/foo/1", "path": "/foo/3"}},
		Expected: map[string]any{"foo": []any{"all", "cows", "eat", "grass"}},
	},
	{
		Comment: "A.8.  Testing a Value: Success",
		Doc:     map[string]any{"baz": "qux", "foo": []any{"a", 2, "c"}},
		Patch: []map[string]any{
			{"op": "test", "path": "/baz", "value": "qux"},
			{"op": "test", "path": "/foo/1", "value": 2},
		},
		Expected: map[string]any{"baz": "qux", "foo": []any{"a", 2, "c"}},
	},
	{
		Comment: "A.9.  Testing a Value: Error",
		Doc:     map[string]any{"baz": "qux"},
		Patch:   []map[string]any{{"op": "test", "path": "/baz", "value": "bar"}},
		Error:   "string not equivalent",
	},
	{
		Comment:  "A.10.  Adding a nested Member Object",
		Doc:      map[string]any{"foo": "bar"},
		Patch:    []map[string]any{{"op": "add", "path": "/child", "value": map[string]any{"grandchild": map[string]any{}}}},
		Expected: map[string]any{"foo": "bar", "child": map[string]any{"grandchild": map[string]any{}}},
	},
	{
		Comment:  "A.11.  Ignoring Unrecognized Elements",
		Doc:      map[string]any{"foo": "bar"},
		Patch:    []map[string]any{{"op": "add", "path": "/baz", "value": "qux", "xyz": 123}},
		Expected: map[string]any{"foo": "bar", "baz": "qux"},
	},
	{
		Comment: "A.12.  Adding to a Non-existent Target",
		Doc:     map[string]any{"foo": "bar"},
		Patch:   []map[string]any{{"op": "add", "path": "/baz/bat", "value": "qux"}},
		Error:   "add to a non-existent target",
	},
	{
		Comment:  "A.13 Invalid JSON Patch Document",
		Doc:      map[string]any{"foo": "bar"},
		Patch:    []map[string]any{{"path": "/baz", "value": "qux", "op": "remove"}},
		Error:    "operation has two 'op' members",
		Disabled: true,
	},
	{
		Comment:  "A.14.  ~ Escape Ordering",
		Doc:      map[string]any{"/": 9, "~1": 10},
		Patch:    []map[string]any{{"op": "test", "path": "/~01", "value": 10}},
		Expected: map[string]any{"/": 9, "~1": 10},
	},
	{
		Comment: "A.15.  Comparing Strings and Numbers",
		Doc:     map[string]any{"/": 9, "~1": 10},
		Patch:   []map[string]any{{"op": "test", "path": "/~01", "value": "10"}},
		Error:   "number is not equal to string",
	},
	{
		Comment:  "A.16.  Adding an Array Value",
		Doc:      map[string]any{"foo": []any{"bar"}},
		Patch:    []map[string]any{{"op": "add", "path": "/foo/-", "value": []any{"abc", "def"}}},
		Expected: map[string]any{"foo": []any{"bar", []any{"abc", "def"}}},
	},
}
