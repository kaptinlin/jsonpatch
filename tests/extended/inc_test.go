package ops_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestIncOp(t *testing.T) {
	t.Run("casts values and then increments them", func(t *testing.T) {
		doc := map[string]interface{}{
			"val1": true,
			"val2": false,
			"val3": 1,
			"val4": 0,
		}
		operations := []internal.Operation{
			{Op: "inc", Path: "/val1", Inc: 1},
			{Op: "inc", Path: "/val2", Inc: 1},
			{Op: "inc", Path: "/val3", Inc: 1},
			{Op: "inc", Path: "/val4", Inc: 1},
		}
		result := testutils.ApplyInternalOps(t, doc, operations)
		expected := map[string]interface{}{
			"val1": float64(2),
			"val2": float64(1),
			"val3": float64(2),
			"val4": float64(1),
		}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("can use arbitrary increment value and can decrement", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": 1,
		}
		operations := []internal.Operation{
			{Op: "inc", Path: "/foo", Inc: 10},
			{Op: "inc", Path: "/foo", Inc: -3},
		}
		result := testutils.ApplyInternalOps(t, doc, operations)
		expected := map[string]interface{}{
			"foo": float64(8),
		}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("increment can be a floating point number", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": 1,
		}
		operations := []internal.Operation{
			{Op: "inc", Path: "/foo", Inc: 0.1},
		}
		result := testutils.ApplyInternalOps(t, doc, operations)
		expected := map[string]interface{}{
			"foo": 1.1,
		}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("root", func(t *testing.T) {
		t.Run("increments from 0 to 5", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "inc",
				Path: "",
				Inc:  5,
			}
			result := testutils.ApplyInternalOps(t, 0, []internal.Operation{operation})
			if result != float64(5) {
				t.Errorf("result = %v, want %v", result, float64(5))
			}
		})

		t.Run("increments from -0 to 5", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "inc",
				Path: "",
				Inc:  5,
			}
			result := testutils.ApplyInternalOps(t, -0, []internal.Operation{operation})
			if result != float64(5) {
				t.Errorf("result = %v, want %v", result, float64(5))
			}
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("increments from 0 to 5", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "inc",
				Path: "/lala",
				Inc:  5,
			}
			result := testutils.ApplyInternalOps(t, map[string]interface{}{"lala": 0}, []internal.Operation{operation})
			expected := map[string]interface{}{"lala": float64(5)}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("increments from -0 to 5", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "inc",
				Path: "/lala",
				Inc:  5,
			}
			result := testutils.ApplyInternalOps(t, map[string]interface{}{"lala": -0}, []internal.Operation{operation})
			expected := map[string]interface{}{"lala": float64(5)}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("casts string to number", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "inc",
				Path: "/lala",
				Inc:  5,
			}
			result := testutils.ApplyInternalOps(t, map[string]interface{}{"lala": "4"}, []internal.Operation{operation})
			expected := map[string]interface{}{"lala": float64(9)}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("can increment twice", func(t *testing.T) {
			operations := []internal.Operation{
				{
					Op:   "inc",
					Path: "/lala",
					Inc:  1,
				},
				{
					Op:   "inc",
					Path: "/lala",
					Inc:  2,
				},
			}
			result := testutils.ApplyInternalOps(t, map[string]interface{}{"lala": 0}, operations)
			expected := map[string]interface{}{"lala": float64(3)}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("creates value when path doesn't exist", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "inc",
				Path: "/newfield",
				Inc:  5,
			}
			result := testutils.ApplyInternalOps(t, map[string]interface{}{}, []internal.Operation{operation})
			expected := map[string]interface{}{"newfield": float64(5)}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("increments from 0 to -3", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "inc",
				Path: "/0",
				Inc:  -3,
			}
			result := testutils.ApplyInternalOps(t, []interface{}{0}, []internal.Operation{operation})
			expected := []interface{}{float64(-3)}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("increments from -0 to -3", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "inc",
				Path: "/0",
				Inc:  -3,
			}
			result := testutils.ApplyInternalOps(t, []interface{}{-0}, []internal.Operation{operation})
			expected := []interface{}{float64(-3)}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	})
}
