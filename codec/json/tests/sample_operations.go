// Package tests provides test data and automated codec tests for the JSON codec.
package tests

// SampleOperations contains all the test operations from TypeScript version
var SampleOperations = map[string]map[string]any{
	// JSON Patch core operations
	"add1": {
		"op":    "add",
		"path":  "",
		"value": nil,
	},
	"add2": {
		"op":    "add",
		"path":  "/foo",
		"value": 123,
	},
	"add3": {
		"op":    "add",
		"path":  "",
		"value": false,
	},
	"add4": {
		"op":   "add",
		"path": "/mentions/2",
		"value": map[string]any{
			"id":   "1234",
			"name": "Joe Jones",
		},
	},
	"add5": {
		"op":   "add",
		"path": "/mentions/2/-",
		"value": map[string]any{
			"id":   "1234",
			"name": "Joe Jones",
		},
	},
	"add6": {
		"op":    "add",
		"path":  "",
		"value": map[string]any{},
	},
	"add7": {
		"op":    "add",
		"path":  "/5/a",
		"value": []any{1, nil, "3"},
	},

	"remove1": {
		"op":   "remove",
		"path": "",
	},
	"remove2": {
		"op":   "remove",
		"path": "/a/b/c/d/e",
	},
	"remove3": {
		"op":       "remove",
		"path":     "/a/b/c/d/e/-",
		"oldValue": nil,
	},
	"remove4": {
		"op":   "remove",
		"path": "/user/123/name",
		"oldValue": map[string]any{
			"firstName": "John",
			"lastName":  "Notjohn",
		},
	},

	"replace1": {
		"op":    "replace",
		"path":  "/a",
		"value": "asdf",
	},
	"replace2": {
		"op":    "replace",
		"path":  "/a/b/c",
		"value": 123,
	},
	"replace3": {
		"op":   "replace",
		"path": "/a/1/-",
		"value": map[string]any{
			"foo": "qux",
		},
		"oldValue": map[string]any{
			"foo": "bar",
		},
	},

	"copy1": {
		"op":   "copy",
		"path": "/foo/bar",
		"from": "/foo/quz",
	},
	"copy2": {
		"op":   "copy",
		"path": "",
		"from": "",
	},
	"copy3": {
		"op":   "copy",
		"path": "/",
		"from": "",
	},
	"copy4": {
		"op":   "copy",
		"path": "/",
		"from": "/a/b/1/2/-",
	},

	"move1": {
		"op":   "move",
		"path": "",
		"from": "",
	},
	"move2": {
		"op":   "move",
		"path": "/a/b/",
		"from": "/c/d/0",
	},

	"test1": {
		"op":    "test",
		"path":  "",
		"value": nil,
	},
	"test2": {
		"op":    "test",
		"path":  "/ha/hi",
		"value": map[string]any{"foo": "bar"},
	},
	"test3": {
		"op":    "test",
		"path":  "/ha/1/2",
		"value": []any{1, map[string]any{"a": false}, "null"},
	},

	// JSON Predicate operations
	"defined1": {
		"op":   "defined",
		"path": "",
	},
	"defined2": {
		"op":   "defined",
		"path": "/1",
	},
	"defined3": {
		"op":   "defined",
		"path": "/foo",
	},

	"undefined1": {
		"op":   "undefined",
		"path": "/a/abc",
	},
	"undefined2": {
		"op":   "undefined",
		"path": "/",
	},
	"undefined3": {
		"op":   "undefined",
		"path": "/1",
	},

	"test_type1": {
		"op":   "test_type",
		"path": "",
		"type": []string{"array"},
	},
	"test_type2": {
		"op":   "test_type",
		"path": "/a/1",
		"type": []string{"integer", "boolean"},
	},
	"test_type3": {
		"op":   "test_type",
		"path": "/a/b/c",
		"type": []string{"array", "boolean", "integer", "null", "number", "object", "string"},
	},

	"test_string1": {
		"op":   "test_string",
		"path": "/a/b/c",
		"pos":  0.0,
		"str":  "asdf",
	},
	"test_string2": {
		"op":   "test_string",
		"path": "/a/1",
		"pos":  4.0,
		"str":  "",
	},

	"test_string_len1": {
		"op":   "test_string_len",
		"path": "/",
		"len":  123.0,
	},
	"test_string_len2": {
		"op":   "test_string_len",
		"path": "/a/bb/ccc",
		"len":  5.0,
		"not":  true,
	},

	"contains1": {
		"op":    "contains",
		"path":  "/a",
		"value": "",
	},
	"contains2": {
		"op":    "contains",
		"path":  "",
		"value": "asdf",
	},

	"matches1": {
		"op":    "matches",
		"path":  "/gg",
		"value": "a",
	},

	"ends1": {
		"op":    "ends",
		"path":  "/foo",
		"value": "",
	},
	"ends2": {
		"op":    "ends",
		"path":  "/",
		"value": "asdf",
	},
	"ends3": {
		"op":          "ends",
		"path":        "/",
		"value":       "asdf",
		"ignore_case": true,
	},

	"starts1": {
		"op":    "starts",
		"path":  "/foo",
		"value": "",
	},
	"starts2": {
		"op":    "starts",
		"path":  "/foo",
		"value": "aa",
	},
	"starts3": {
		"op":          "starts",
		"path":        "/foo",
		"value":       "aa",
		"ignore_case": true,
	},

	"type1": {
		"op":    "type",
		"path":  "/1/2/3",
		"value": "array",
	},
	"type2": {
		"op":    "type",
		"path":  "/1/2/3",
		"value": "boolean",
	},
	"type3": {
		"op":    "type",
		"path":  "/1/2/3",
		"value": "integer",
	},
	"type4": {
		"op":    "type",
		"path":  "/1/2/3",
		"value": "null",
	},
	"type5": {
		"op":    "type",
		"path":  "/1/2/3",
		"value": "number",
	},
	"type6": {
		"op":    "type",
		"path":  "/1/2/3",
		"value": "object",
	},
	"type7": {
		"op":    "type",
		"path":  "/1/2/3",
		"value": "string",
	},

	"in1": {
		"op":    "in",
		"path":  "/",
		"value": []any{"asdf"},
	},
	"in2": {
		"op":    "in",
		"path":  "/foo/bar",
		"value": []any{"asdf", 132, map[string]any{"a": "b"}, nil},
	},

	"less1": {
		"op":    "less",
		"path":  "/z",
		"value": -0.5,
	},
	"less2": {
		"op":    "less",
		"path":  "",
		"value": 0,
	},

	"more1": {
		"op":    "more",
		"path":  "",
		"value": 1,
	},
	"more2": {
		"op":    "more",
		"path":  "/a",
		"value": -1,
	},

	"and1": {
		"op":   "and",
		"path": "/a",
		"apply": []any{
			map[string]any{"op": "test", "path": "/b", "value": 123},
		},
	},
	"and2": {
		"op":   "and",
		"path": "/",
		"apply": []any{
			map[string]any{"op": "less", "path": "", "value": 0},
			map[string]any{"op": "more", "path": "", "value": 1},
			map[string]any{"op": "in", "path": "/", "value": []any{"asdf"}},
		},
	},
	"and3": {
		"op":   "and",
		"path": "/a/1/.",
		"apply": []any{
			map[string]any{"op": "test", "path": "", "value": nil},
			map[string]any{"op": "test", "path": "/ha/hi", "value": map[string]any{"foo": "bar"}},
			map[string]any{"op": "test", "path": "/ha/1/2", "value": []any{1, map[string]any{"a": false}, "null"}},
		},
	},
	"and4": {
		"op":   "and",
		"path": "/a/1/.",
		"apply": []any{
			map[string]any{"op": "test", "path": "", "value": nil},
			map[string]any{
				"op":   "and",
				"path": "/gg/bet",
				"apply": []any{
					map[string]any{"op": "test", "path": "", "value": nil},
					map[string]any{"op": "test", "path": "/ha/hi", "value": map[string]any{"foo": "bar"}},
				},
			},
			map[string]any{"op": "test", "path": "/ha/hi", "value": map[string]any{"foo": "bar"}},
		},
	},

	"not1": {
		"op":   "not",
		"path": "/",
		"apply": []any{
			map[string]any{"op": "less", "path": "", "value": 0},
			map[string]any{"op": "more", "path": "", "value": 1},
			map[string]any{"op": "in", "path": "/", "value": []any{"asdf"}},
		},
	},
	"not2": {
		"op":   "not",
		"path": "/a/1/.",
		"apply": []any{
			map[string]any{"op": "test", "path": "", "value": nil},
			map[string]any{"op": "test", "path": "/ha/hi", "value": map[string]any{"foo": "bar"}},
			map[string]any{"op": "test", "path": "/ha/1/2", "value": []any{1, map[string]any{"a": false}, "null"}},
		},
	},

	"or1": {
		"op":   "or",
		"path": "/",
		"apply": []any{
			map[string]any{"op": "less", "path": "", "value": 0},
			map[string]any{"op": "more", "path": "", "value": 1},
			map[string]any{"op": "in", "path": "/", "value": []any{"asdf"}},
		},
	},
	"or2": {
		"op":   "or",
		"path": "/a/1/.",
		"apply": []any{
			map[string]any{"op": "test", "path": "", "value": nil},
			map[string]any{"op": "test", "path": "/ha/hi", "value": map[string]any{"foo": "bar"}},
			map[string]any{"op": "test", "path": "/ha/1/2", "value": []any{1, map[string]any{"a": false}, "null"}},
		},
	},

	// JSON Patch Extended operations
	"str_ins1": {
		"op":   "str_ins",
		"path": "/ads",
		"pos":  0.0,
		"str":  "",
	},
	"str_ins2": {
		"op":   "str_ins",
		"path": "/a/b/lkasjdfoiasjdfoiasjdflaksjdflkasjfljasdflkjasdlfjkasdf",
		"pos":  823848493.0,
		"str":  "Component model",
	},

	"str_del1": {
		"op":   "str_del",
		"path": "",
		"pos":  0.0,
		"len":  0.0,
	},
	"str_del2": {
		"op":   "str_del",
		"path": "/asdfasdfasdfasdfasdfasdfasdfpalsdf902039joij2130j9e2093k2309k203f0sjdf0s9djf0skdfs0dfk092j0239j0mospdkf",
		"pos":  92303948.0,
		"len":  84487.0,
	},
	"str_del3": {
		"op":   "str_del",
		"path": "/asdf/reg/asdf/asdf/wer/sdaf234/asf/23/asdf2/asdf2",
		"pos":  92303948.0,
		"str":  "asdfasdfasdfasdflkasdjflakjsdf",
	},

	"flip1": {
		"op":   "flip",
		"path": "",
	},
	"flip2": {
		"op":   "flip",
		"path": "/asdf/df/dfa/dfasfd/",
	},

	"inc1": {
		"op":   "inc",
		"path": "/",
		"inc":  1.0,
	},
	"inc2": {
		"op":   "inc",
		"path": "/asdf/sd/d/f",
		"inc":  -123.0,
	},

	"split1": {
		"op":   "split",
		"path": "/i",
		"pos":  0.0,
	},
	"split2": {
		"op":   "split",
		"path": "/i/asdf/sdf/d",
		"pos":  123.0,
	},
	"split3": {
		"op":   "split",
		"path": "/i/asdf/sdf/d",
		"pos":  123.0,
		"props": map[string]any{
			"foo": "bar",
			"a":   123,
		},
	},

	"merge1": {
		"op":   "merge",
		"path": "",
		"pos":  0.0,
	},
	"merge2": {
		"op":   "merge",
		"path": "/a/b/c",
		"pos":  123412341234.0,
		"props": map[string]any{
			"foo": nil,
			"bar": 23,
			"baz": "asdf",
			"quz": true,
			"qux": []any{1, "2", 3, true, false, nil},
		},
	},

	"extend1": {
		"op":    "extend",
		"path":  "/",
		"props": map[string]any{},
	},
	"extend2": {
		"op":   "extend",
		"path": "/asdf/asdf/asdf",
		"props": map[string]any{
			"foo": "bar",
		},
	},
	"extend3": {
		"op":   "extend",
		"path": "/asdf/asdf/asdf",
		"props": map[string]any{
			"foo": "bar",
			"a":   nil,
			"b":   true,
		},
		"deleteNull": true,
	},
}
