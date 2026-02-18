package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
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
		t.Fatalf("Test() unexpected error: %v", err)
	}
	if !ok {
		assert.Fail(t, "Test() = false, want true for equal values")
	}

	testOp = NewTest([]string{"foo"}, "qux")
	ok, err = testOp.Test(doc)
	if err != nil {
		t.Fatalf("Test() unexpected error: %v", err)
	}
	if ok {
		assert.Fail(t, "Test() = true, want false for different values")
	}
}

func TestTest_Apply(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	result, err := testOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
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
		t.Fatalf("ToJSON() unexpected error: %v", err)
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
		t.Fatalf("ToCompact() unexpected error: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want %d", len(compact), 3)
	}
	assert.Equal(t, internal.OpTestCode, compact[0], "compact[0]")
	assert.Equal(t, []string{"foo"}, compact[1])
	assert.Equal(t, "bar", compact[2], "compact[2]")
}

func TestTest_Validate(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")
	if err := testOp.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	testOp = NewTest([]string{}, "bar")
	if err := testOp.Validate(); err == nil {
		assert.Fail(t, "Validate() expected error for empty path")
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
