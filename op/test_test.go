package op

import (
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestTest_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	testOp := NewTest([]string{"foo"}, "bar")

	ok, err := testOp.Test(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Test() unexpected error: %v", err))
	}
	if !ok {
		assert.Fail(t, "Test() = false, want true for equal values")
	}

	testOp = NewTest([]string{"foo"}, "qux")
	ok, err = testOp.Test(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Test() unexpected error: %v", err))
	}
	if ok {
		assert.Fail(t, "Test() = true, want false for different values")
	}
}

func TestTest_UsesExactNumericEquality(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  any
		right any
		want  bool
	}{
		{name: "adjacent integers above float precision", left: int64(1 << 53), right: int64(1<<53 + 1), want: false},
		{name: "adjacent maximum unsigned integers", left: uint64(math.MaxUint64), right: uint64(math.MaxUint64 - 1), want: false},
		{name: "maximum signed versus rounded float", left: int64(math.MaxInt64), right: float64(math.MaxInt64), want: false},
		{name: "negative signed versus unsigned", left: int64(-1), right: uint64(math.MaxUint64), want: false},
		{name: "exact signed and unsigned", left: int64(1), right: uint64(1), want: true},
		{name: "exact integer and float", left: int64(1), right: float64(1), want: true},
		{name: "fractional floats", left: float32(1.5), right: float64(1.5), want: true},
		{name: "positive infinity", left: math.Inf(1), right: math.Inf(1), want: true},
		{name: "not a number", left: math.NaN(), right: math.NaN(), want: false},
		{
			name:  "nested large integers",
			left:  map[string]any{"values": []any{int64(1 << 53)}},
			right: map[string]any{"values": []any{int64(1<<53 + 1)}},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			matched, err := NewTest(nil, tt.right).Test(tt.left)
			require.NoError(t, err)
			assert.Equal(t, tt.want, matched)
		})
	}
}

func TestIn_UsesExactNumericEquality(t *testing.T) {
	t.Parallel()

	matched, err := NewIn(nil, []any{int64(1<<53 + 1)}).Test(int64(1 << 53))
	require.NoError(t, err)
	assert.False(t, matched)

	matched, err = NewIn(nil, []any{float64(1)}).Test(uint64(1))
	require.NoError(t, err)
	assert.True(t, matched)
}

func TestTest_Apply(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	result, err := testOp.Apply(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() unexpected error: %v", err))
	}
	assert.Equal(t, doc, result.Doc)

	testOp = NewTest([]string{"foo"}, "qux")
	_, err = testOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "Apply() expected error for non-matching values")
	}
}

func TestTest_ToJSON(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")

	got, err := testOp.ToJSON()
	if err != nil {
		require.FailNow(t, fmt.Sprintf("ToJSON() unexpected error: %v", err))
	}

	assert.Equal(t, "test", got.Op, "ToJSON().Op")
	assert.Equal(t, "/foo", got.Path, "ToJSON().Path")
	assert.Equal(t, "bar", got.Value, "ToJSON().Value")
}

func TestTest_ToCompact(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")

	compact, err := testOp.ToCompact()
	if err != nil {
		require.FailNow(t, fmt.Sprintf("ToCompact() unexpected error: %v", err))
	}
	if len(compact) != 3 {
		require.FailNow(t, fmt.Sprintf("len(ToCompact()) = %d, want %d", len(compact), 3))
	}
	assert.Equal(t, internal.OpTestCode, compact[0], "compact[0]")
	if diff := cmp.Diff([]string{"foo"}, compact[1]); diff != "" {
		t.Errorf("compact[1] mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, "bar", compact[2], "compact[2]")
}

func TestTest_NegatedContract(t *testing.T) {
	t.Parallel()

	doc := map[string]any{"foo": "bar"}
	testOp := NewTestWithNot([]string{"foo"}, "bar", true)
	assert.True(t, testOp.Not())

	matched, err := testOp.Test(doc)
	assert.NoError(t, err)
	assert.False(t, matched)

	matched, err = testOp.Test(map[string]any{"other": "bar"})
	assert.NoError(t, err)
	assert.True(t, matched)

	result, err := testOp.Apply(map[string]any{"other": "bar"})
	assert.NoError(t, err)
	if diff := cmp.Diff(map[string]any{"other": "bar"}, result.Doc); diff != "" {
		t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
	}

	_, err = testOp.Apply(doc)
	assert.ErrorIs(t, err, ErrTestOperationFailed)

	jsonOp, err := testOp.ToJSON()
	assert.NoError(t, err)
	wantJSON := internal.Operation{Op: "test", Path: "/foo", Value: "bar", Not: true}
	if diff := cmp.Diff(wantJSON, jsonOp); diff != "" {
		t.Errorf("ToJSON() mismatch (-want +got):\n%s", diff)
	}

	compactOp, err := testOp.ToCompact()
	assert.NoError(t, err)
	wantCompact := internal.CompactOperation{internal.OpTestCode, []string{"foo"}, "bar", 1}
	if diff := cmp.Diff(wantCompact, compactOp); diff != "" {
		t.Errorf("ToCompact() mismatch (-want +got):\n%s", diff)
	}
}

func TestTest_Validate(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")
	if err := testOp.Validate(); err != nil {
		assert.Fail(t, fmt.Sprintf("Validate() unexpected error: %v", err))
	}

	// Empty path is valid root-document testing.
	testOp = NewTest([]string{}, "bar")
	if err := testOp.Validate(); err != nil {
		assert.Fail(t, fmt.Sprintf("Validate() unexpected error for empty path: %v", err))
	}
}

func TestTest_InterfaceMethods(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")

	if got := testOp.Op(); got != internal.OpTestType {
		assert.Equal(t, internal.OpTestType, got, "Op()")
	}
	if got := testOp.Code(); got != internal.OpTestCode {
		assert.Equal(t, internal.OpTestCode, got, "Code()")
	}
	assert.Equal(t, []string{"foo"}, testOp.Path(), "Path()")
}

func FuzzDeepEqualIntegers(f *testing.F) {
	f.Add(int64(0), int64(0), uint64(0))
	f.Add(int64(1<<53), int64(1<<53+1), uint64(1<<53))
	f.Add(int64(math.MaxInt64), int64(math.MaxInt64-1), uint64(math.MaxUint64))
	f.Add(int64(-1), int64(0), uint64(math.MaxUint64))

	f.Fuzz(func(t *testing.T, left, right int64, unsigned uint64) {
		assert.Equal(t, left == right, deepEqual(left, right))
		assert.Equal(t, left >= 0 && uint64(left) == unsigned, deepEqual(left, unsigned))
		assert.Equal(t, right >= 0 && uint64(right) == unsigned, deepEqual(unsigned, right))
	})
}
