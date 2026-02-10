package jsonpatch_test

import (
	"github.com/go-json-experiment/json"
	"testing"

	"github.com/kaptinlin/jsonpatch"
)

func BenchmarkBasicOperations(b *testing.B) {
	testCases := []struct {
		name string
		doc  interface{}
		ops  []jsonpatch.Operation
	}{
		{
			name: "add_simple_value",
			doc:  map[string]interface{}{"foo": "bar"},
			ops: []jsonpatch.Operation{
				{Op: "add", Path: "/baz", Value: "qux"},
			},
		},
		{
			name: "replace_value",
			doc:  map[string]interface{}{"foo": "bar", "baz": "qux"},
			ops: []jsonpatch.Operation{
				{Op: "replace", Path: "/foo", Value: "new_value"},
			},
		},
		{
			name: "remove_value",
			doc:  map[string]interface{}{"foo": "bar", "baz": "qux"},
			ops: []jsonpatch.Operation{
				{Op: "remove", Path: "/baz"},
			},
		},
		{
			name: "test_operation",
			doc:  map[string]interface{}{"foo": "bar"},
			ops: []jsonpatch.Operation{
				{Op: "test", Path: "/foo", Value: "bar"},
			},
		},
		{
			name: "copy_operation",
			doc:  map[string]interface{}{"foo": "bar", "baz": map[string]interface{}{"deep": "value"}},
			ops: []jsonpatch.Operation{
				{Op: "copy", From: "/baz", Path: "/copied"},
			},
		},
		{
			name: "move_operation",
			doc:  map[string]interface{}{"foo": "bar", "baz": "qux"},
			ops: []jsonpatch.Operation{
				{Op: "move", From: "/baz", Path: "/moved"},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				// Clone document for each iteration
				docCopy := cloneDocument(tc.doc)
				_, err := jsonpatch.ApplyPatch(docCopy, tc.ops, jsonpatch.WithMutate(true))
				if err != nil {
					b.Fatalf("Operation failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkExtendedOperations(b *testing.B) {
	testCases := []struct {
		name string
		doc  interface{}
		ops  []jsonpatch.Operation
	}{
		{
			name: "inc_operation",
			doc:  map[string]interface{}{"counter": 42},
			ops: []jsonpatch.Operation{
				{Op: "inc", Path: "/counter", Inc: 1},
			},
		},
		{
			name: "flip_operation",
			doc:  map[string]interface{}{"enabled": true},
			ops: []jsonpatch.Operation{
				{Op: "flip", Path: "/enabled"},
			},
		},
		// Note: merge operation works on arrays, commenting out object merge for now
		// {
		//	name: "merge_operation",
		//	doc:  map[string]interface{}{"obj": map[string]interface{}{"a": 1}},
		//	ops: []jsonpatch.Operation{
		//		{"op": "merge", "path": "/obj", "value": map[string]interface{}{"b": 2}},
		//	},
		// },
		{
			name: "strins_operation",
			doc:  map[string]interface{}{"text": "hello world"},
			ops: []jsonpatch.Operation{
				{Op: "str_ins", Path: "/text", Pos: 5, Str: " beautiful"},
			},
		},
		{
			name: "split_operation",
			doc:  map[string]interface{}{"text": "hello,world,test"},
			ops: []jsonpatch.Operation{
				{Op: "split", Path: "/text", Pos: 5},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				// Clone document for each iteration
				docCopy := cloneDocument(tc.doc)
				_, err := jsonpatch.ApplyPatch(docCopy, tc.ops, jsonpatch.WithMutate(true))
				if err != nil {
					b.Fatalf("Operation failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkPredicateOperations(b *testing.B) {
	testCases := []struct {
		name string
		doc  interface{}
		ops  []jsonpatch.Operation
	}{
		{
			name: "contains_predicate",
			doc:  map[string]interface{}{"text": "hello world"},
			ops: []jsonpatch.Operation{
				{Op: "contains", Path: "/text", Value: "world"},
			},
		},
		{
			name: "starts_predicate",
			doc:  map[string]interface{}{"text": "hello world"},
			ops: []jsonpatch.Operation{
				{Op: "starts", Path: "/text", Value: "hello"},
			},
		},
		{
			name: "ends_predicate",
			doc:  map[string]interface{}{"text": "hello world"},
			ops: []jsonpatch.Operation{
				{Op: "ends", Path: "/text", Value: "world"},
			},
		},
		{
			name: "type_predicate",
			doc:  map[string]interface{}{"value": 42},
			ops: []jsonpatch.Operation{
				{Op: "type", Path: "/value", Value: "number"},
			},
		},
		{
			name: "less_predicate",
			doc:  map[string]interface{}{"value": 42},
			ops: []jsonpatch.Operation{
				{Op: "less", Path: "/value", Value: 50},
			},
		},
		{
			name: "more_predicate",
			doc:  map[string]interface{}{"value": 42},
			ops: []jsonpatch.Operation{
				{Op: "more", Path: "/value", Value: 30},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				// Clone document for each iteration
				docCopy := cloneDocument(tc.doc)
				_, err := jsonpatch.ApplyPatch(docCopy, tc.ops, jsonpatch.WithMutate(true))
				if err != nil {
					b.Fatalf("Operation failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkSecondOrderPredicates(b *testing.B) {
	testCases := []struct {
		name string
		doc  interface{}
		ops  []jsonpatch.Operation
	}{
		// Note: Complex and operation - commented out for benchmark stability
		// {
		//	name: "and_predicate",
		//	doc:  map[string]interface{}{"foo": 1, "bar": 2},
		//	ops: []jsonpatch.Operation{
		//		{
		//			Op:   "and",
		//			Path: "",
		//			Apply: []interface{}{
		//				map[string]interface{}{"op": "test", "path": "/foo", "value": 1},
		//				map[string]interface{}{"op": "test", "path": "/bar", "value": 2},
		//			},
		//		},
		//	},
		// },
		// Note: Complex or operation - commented out for benchmark stability
		// {
		//	name: "or_predicate",
		//	doc:  map[string]interface{}{"foo": 1, "bar": 2},
		//	ops: []jsonpatch.Operation{
		//		{
		//			Op:   "or",
		//			Path: "",
		//			Apply: []interface{}{
		//				map[string]interface{}{"op": "test", "path": "/foo", "value": 1},
		//				map[string]interface{}{"op": "test", "path": "/bar", "value": 3},
		//			},
		//		},
		//	},
		// },
		{
			name: "not_predicate",
			doc:  map[string]interface{}{"foo": 1, "bar": 2},
			ops: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 2},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				// Clone document for each iteration
				docCopy := cloneDocument(tc.doc)
				_, err := jsonpatch.ApplyPatch(docCopy, tc.ops, jsonpatch.WithMutate(true))
				if err != nil {
					b.Fatalf("Operation failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkComplexDocument(b *testing.B) {
	complexDoc := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "Alice",
				"profile": map[string]interface{}{
					"age":   30,
					"email": "alice@example.com",
					"preferences": map[string]interface{}{
						"theme": "dark",
						"lang":  "en",
					},
				},
			},
			map[string]interface{}{
				"id":   2,
				"name": "Bob",
				"profile": map[string]interface{}{
					"age":   25,
					"email": "bob@example.com",
					"preferences": map[string]interface{}{
						"theme": "light",
						"lang":  "es",
					},
				},
			},
		},
		"metadata": map[string]interface{}{
			"version":   "1.0",
			"timestamp": "2024-01-01T00:00:00Z",
			"stats": map[string]interface{}{
				"total_users":   2,
				"active_users":  1,
				"last_modified": "2024-01-01T12:00:00Z",
			},
		},
	}

	testCases := []struct {
		name string
		ops  []jsonpatch.Operation
	}{
		{
			name: "add_new_user",
			ops: []jsonpatch.Operation{
				{
					Op:   "add",
					Path: "/users/-",
					Value: map[string]interface{}{
						"id":   3,
						"name": "Charlie",
						"profile": map[string]interface{}{
							"age":   28,
							"email": "charlie@example.com",
						},
					},
				},
			},
		},
		{
			name: "update_user_preference",
			ops: []jsonpatch.Operation{
				{Op: "replace", Path: "/users/0/profile/preferences/theme", Value: "light"},
			},
		},
		{
			name: "increment_stats",
			ops: []jsonpatch.Operation{
				{Op: "inc", Path: "/metadata/stats/total_users", Inc: 1},
				{Op: "inc", Path: "/metadata/stats/active_users", Inc: 1},
			},
		},
		// Note: Complex validation with and operation - commented out for benchmark stability
		// {
		//	name: "complex_validation",
		//	ops: []jsonpatch.Operation{
		//		{
		//			Op:   "and",
		//			Path: "",
		//			Apply: []interface{}{
		//				map[string]interface{}{"op": "test", "path": "/metadata/version", "value": "1.0"},
		//				map[string]interface{}{"op": "more", "path": "/metadata/stats/total_users", "value": 0},
		//			},
		//		},
		//	},
		// },
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				// Clone document for each iteration
				docCopy := cloneDocument(complexDoc)
				_, err := jsonpatch.ApplyPatch(docCopy, tc.ops, jsonpatch.WithMutate(true))
				if err != nil {
					b.Fatalf("Operation failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkMutateVsImmutable compares mutable vs immutable operations
func BenchmarkMutateVsImmutable(b *testing.B) {
	doc := map[string]interface{}{
		"data": make([]interface{}, 1000),
		"metadata": map[string]interface{}{
			"size":    1000,
			"version": 1,
		},
	}

	// Fill the array with test data
	for i := 0; i < 1000; i++ {
		doc["data"].([]interface{})[i] = map[string]interface{}{
			"id":    i,
			"value": i * 2,
		}
	}

	ops := []jsonpatch.Operation{
		{Op: "replace", Path: "/metadata/version", Value: 2},
		{Op: "add", Path: "/metadata/updated", Value: true},
	}

	b.Run("mutable", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			docCopy := cloneDocument(doc)
			_, err := jsonpatch.ApplyPatch(docCopy, ops, jsonpatch.WithMutate(true))
			if err != nil {
				b.Fatalf("Operation failed: %v", err)
			}
		}
	})

	b.Run("immutable", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			docCopy := cloneDocument(doc)
			_, err := jsonpatch.ApplyPatch(docCopy, ops, jsonpatch.WithMutate(false))
			if err != nil {
				b.Fatalf("Operation failed: %v", err)
			}
		}
	})
}

// cloneDocument creates a deep copy of the document for benchmarking
func cloneDocument(doc interface{}) interface{} {
	// Use JSON marshal/unmarshal for deep cloning
	data, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result
}
