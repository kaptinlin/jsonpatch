package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestAnd_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{}, []any{test1, test2})

	ok, err := andOp.Test(doc)
	if err != nil {
		t.Fatalf("And.Test() failed: %v", err)
	}
	if !ok {
		assert.Fail(t, "And.Test() = false, want true when all operations pass")
	}
}

func TestAnd_OneFails(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 456)

	andOp := NewAnd([]string{}, []any{test1, test2})

	ok, err := andOp.Test(doc)
	if err != nil {
		t.Fatalf("And.Test() failed: %v", err)
	}
	if ok {
		assert.Fail(t, "And.Test() = true, want false when any operation fails")
	}
}

func TestAnd_Empty(t *testing.T) {
	t.Parallel()
	andOp := NewAnd([]string{}, []any{})

	doc := map[string]any{"foo": "bar"}
	ok, err := andOp.Test(doc)
	if err != nil {
		t.Fatalf("And.Test() failed: %v", err)
	}
	if !ok {
		assert.Fail(t, "And.Test() = false, want true for empty AND (vacuous truth)")
	}
}

func TestAnd_Apply(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{}, []any{test1, test2})

	result, err := andOp.Apply(doc)
	if err != nil {
		t.Fatalf("And.Apply() failed: %v", err)
	}
	assert.Equal(t, doc, result.Doc)
}

func TestAnd_Apply_Fails(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 456)

	andOp := NewAnd([]string{}, []any{test1, test2})

	_, err := andOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "And.Apply() succeeded, want error when any operation fails")
	}
	if !errors.Is(err, ErrAndTestFailed) {
		assert.Equal(t, ErrAndTestFailed, err, "And.Apply() error")
	}
}

func TestAnd_InterfaceMethods(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	if got := andOp.Op(); got != internal.OpAndType {
		assert.Equal(t, internal.OpAndType, got, "Op()")
	}
	if got := andOp.Code(); got != internal.OpAndCode {
		assert.Equal(t, internal.OpAndCode, got, "Code()")
	}
	if diff := cmp.Diff([]string{"test"}, andOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}

	ops := andOp.Ops()
	if len(ops) != 2 {
		t.Fatalf("len(Ops()) = %d, want 2", len(ops))
	}
	if ops[0] != test1 {
		assert.Fail(t, "Ops()[0] does not match first operation")
	}
	if ops[1] != test2 {
		assert.Fail(t, "Ops()[1] does not match second operation")
	}
}

func TestAnd_ToJSON(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	got, err := andOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}
	if got.Op != "and" {
		assert.Equal(t, "and", got.Op, "ToJSON().Op")
	}
	if got.Path != "/test" {
		assert.Equal(t, "/test", got.Path, "ToJSON().Path")
	}
	if len(got.Apply) != 2 {
		t.Errorf("len(ToJSON().Apply) = %d, want 2", len(got.Apply))
	}
}

func TestAnd_ToCompact(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	compact, err := andOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() failed: %v", err)
	}
	assert.Equal(t, internal.OpAndCode, compact[0], "ToCompact()[0]")
	assert.Equal(t, []string{"test"}, compact[1])
	if _, ok := compact[2].([]any); !ok {
		t.Errorf("ToCompact()[2] type = %T, want []any", compact[2])
	}
}

func TestAnd_Validate(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})
	if err := andOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for valid operation", err)
	}

	// Vacuous truth allows empty operations
	andOp = NewAnd([]string{"test"}, []any{})
	if err := andOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for empty operations (vacuous truth)", err)
	}
}
