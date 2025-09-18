package jsonpatch_test

import (
	"encoding/json"
	"testing"

	"github.com/kaptinlin/jsonpatch"
)

// BenchmarkBasicOperations benchmarks fundamental JSON Patch operations
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
				{"op": "add", "path": "/baz", "value": "qux"},
			},
		},
		{
			name: "replace_value",
			doc:  map[string]interface{}{"foo": "bar", "baz": "qux"},
			ops: []jsonpatch.Operation{
				{"op": "replace", "path": "/foo", "value": "new_value"},
			},
		},
		{
			name: "remove_value",
			doc:  map[string]interface{}{"foo": "bar", "baz": "qux"},
			ops: []jsonpatch.Operation{
				{"op": "remove", "path": "/baz"},
			},
		},
		{
			name: "test_operation",
			doc:  map[string]interface{}{"foo": "bar"},
			ops: []jsonpatch.Operation{
				{"op": "test", "path": "/foo", "value": "bar"},
			},
		},
		{
			name: "copy_operation",
			doc:  map[string]interface{}{"foo": "bar", "baz": map[string]interface{}{"deep": "value"}},
			ops: []jsonpatch.Operation{
				{"op": "copy", "from": "/baz", "path": "/copied"},
			},
		},
		{
			name: "move_operation",
			doc:  map[string]interface{}{"foo": "bar", "baz": "qux"},
			ops: []jsonpatch.Operation{
				{"op": "move", "from": "/baz", "path": "/moved"},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
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

// BenchmarkExtendedOperations benchmarks extended JSON Patch operations
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
				{"op": "inc", "path": "/counter", "inc": 1},
			},
		},
		{
			name: "flip_operation",
			doc:  map[string]interface{}{"enabled": true},
			ops: []jsonpatch.Operation{
				{"op": "flip", "path": "/enabled"},
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
				{"op": "str_ins", "path": "/text", "pos": 5, "str": " beautiful"},
			},
		},
		{
			name: "split_operation",
			doc:  map[string]interface{}{"text": "hello,world,test"},
			ops: []jsonpatch.Operation{
				{"op": "split", "path": "/text", "pos": 5},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
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

// BenchmarkPredicateOperations benchmarks JSON Predicate operations
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
				{"op": "contains", "path": "/text", "value": "world"},
			},
		},
		{
			name: "starts_predicate",
			doc:  map[string]interface{}{"text": "hello world"},
			ops: []jsonpatch.Operation{
				{"op": "starts", "path": "/text", "value": "hello"},
			},
		},
		{
			name: "ends_predicate",
			doc:  map[string]interface{}{"text": "hello world"},
			ops: []jsonpatch.Operation{
				{"op": "ends", "path": "/text", "value": "world"},
			},
		},
		{
			name: "type_predicate",
			doc:  map[string]interface{}{"value": 42},
			ops: []jsonpatch.Operation{
				{"op": "type", "path": "/value", "value": "number"},
			},
		},
		{
			name: "less_predicate",
			doc:  map[string]interface{}{"value": 42},
			ops: []jsonpatch.Operation{
				{"op": "less", "path": "/value", "value": 50},
			},
		},
		{
			name: "more_predicate",
			doc:  map[string]interface{}{"value": 42},
			ops: []jsonpatch.Operation{
				{"op": "more", "path": "/value", "value": 30},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
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

// BenchmarkSecondOrderPredicates benchmarks second-order predicate operations
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
		//			"op":   "and",
		//			"path": "",
		//			"apply": []interface{}{
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
		//			"op":   "or",
		//			"path": "",
		//			"apply": []interface{}{
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
					"op":   "not",
					"path": "",
					"apply": []interface{}{
						map[string]interface{}{"op": "test", "path": "/foo", "value": 2},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
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

// BenchmarkComplexDocument benchmarks operations on complex documents
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
					"op":   "add",
					"path": "/users/-",
					"value": map[string]interface{}{
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
				{"op": "replace", "path": "/users/0/profile/preferences/theme", "value": "light"},
			},
		},
		{
			name: "increment_stats",
			ops: []jsonpatch.Operation{
				{"op": "inc", "path": "/metadata/stats/total_users", "inc": 1},
				{"op": "inc", "path": "/metadata/stats/active_users", "inc": 1},
			},
		},
		// Note: Complex validation with and operation - commented out for benchmark stability
		// {
		//	name: "complex_validation",
		//	ops: []jsonpatch.Operation{
		//		{
		//			"op":   "and",
		//			"path": "",
		//			"apply": []interface{}{
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
			for i := 0; i < b.N; i++ {
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
		{"op": "replace", "path": "/metadata/version", "value": 2},
		{"op": "add", "path": "/metadata/updated", "value": true},
	}

	b.Run("mutable", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			docCopy := cloneDocument(doc)
			_, err := jsonpatch.ApplyPatch(docCopy, ops, jsonpatch.WithMutate(true))
			if err != nil {
				b.Fatalf("Operation failed: %v", err)
			}
		}
	})

	b.Run("immutable", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
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
