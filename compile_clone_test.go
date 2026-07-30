package jsonpatch

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

var (
	errCloneFailed  = errors.New("clone failed")
	errCloneInvalid = errors.New("clone invalid")
)

func TestCompileReportsCloneFailure(t *testing.T) {
	t.Parallel()

	patch, err := Compile(cloneFailingOperation{})
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, ErrPayloadInvalid)
	assert.ErrorIs(t, err, errCloneFailed)

	var patchErr *Error
	require.ErrorAs(t, err, &patchErr)
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "add", patchErr.Op())
	assert.Equal(t, "/name", patchErr.Path())
}

func TestCompileRejectsNilClone(t *testing.T) {
	t.Parallel()

	patch, err := Compile(nilCloningOperation{})
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, ErrPayloadInvalid)

	var patchErr *Error
	require.ErrorAs(t, err, &patchErr)
	assert.Equal(t, 0, patchErr.Index())
	assert.Equal(t, "add", patchErr.Op())
	assert.Equal(t, "/name", patchErr.Path())
}

func TestCompileRejectsTypedNilOperation(t *testing.T) {
	t.Parallel()

	var operation *typedNilOperation
	patch, err := Compile(operation)
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, ErrPayloadInvalid)
}

func TestCompileRejectsTypedNilClone(t *testing.T) {
	t.Parallel()

	patch, err := Compile(typedNilCloningOperation{})
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, ErrPayloadInvalid)
}

func TestCompileValidatesClone(t *testing.T) {
	t.Parallel()

	patch, err := Compile(cloneBecomesInvalidOperation{})
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, ErrPayloadInvalid)
	assert.ErrorIs(t, err, errCloneInvalid)
}

func TestCompileDoesNotValidateSourceAfterClone(t *testing.T) {
	t.Parallel()

	patch, err := Compile(cloneBecomesValidOperation{})
	require.NoError(t, err)
	require.NotNil(t, patch)
	assert.Equal(t, 1, patch.Len())
}

func TestCompileChecksCloneCapability(t *testing.T) {
	t.Parallel()

	patch, err := Compile(cloneBecomesUnsupportedOperation{})
	require.Error(t, err)
	assert.Nil(t, patch)
	assert.ErrorIs(t, err, ErrUnsupportedCapability)

	var patchErr *Error
	require.ErrorAs(t, err, &patchErr)
	assert.Equal(t, "inc", patchErr.Op())
	assert.Equal(t, "/count", patchErr.Path())
}

func TestCompileJSONPreservesLargeCoordinates(t *testing.T) {
	t.Parallel()
	if strconv.IntSize < 64 {
		t.Skip("requires a 64-bit int")
	}

	patch, err := CompileJSON(
		[]byte(`[{"op":"str_ins","path":"/text","pos":9007199254740993,"str":"x"},{"op":"and","path":"/profile","apply":[{"op":"test_string","path":"/name","pos":9007199254740993,"str":"A"}]}]`),
		WithCapabilities(Extended, Predicate),
	)
	require.NoError(t, err)
	require.Len(t, patch.ops, 2)

	large := int64(9007199254740993)
	assert.Equal(t, int(large), patch.ops[0].(*op.StrInsOperation).Pos)
	and := patch.ops[1].(*op.AndOperation)
	require.Len(t, and.Operations, 1)
	assert.Equal(t, int(large), and.Operations[0].(*op.TestStringOperation).Pos)
}

type cloneFailingOperation struct{}

func (cloneFailingOperation) Op() internal.OpType { return internal.OpAddType }
func (cloneFailingOperation) Path() []string      { return []string{"name"} }
func (cloneFailingOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (cloneFailingOperation) Validate() error             { return nil }
func (cloneFailingOperation) Clone() (internal.Op, error) { return nil, errCloneFailed }

type nilCloningOperation struct{}

func (nilCloningOperation) Op() internal.OpType { return internal.OpAddType }
func (nilCloningOperation) Path() []string      { return []string{"name"} }
func (nilCloningOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (nilCloningOperation) Validate() error             { return nil }
func (nilCloningOperation) Clone() (internal.Op, error) { return nil, nil }

type typedNilOperation struct {
	path []string
}

func (*typedNilOperation) Op() internal.OpType { return internal.OpAddType }
func (o *typedNilOperation) Path() []string    { return o.path }
func (*typedNilOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (*typedNilOperation) Validate() error { return nil }
func (o *typedNilOperation) Clone() (internal.Op, error) {
	return &typedNilOperation{path: o.path}, nil
}

type typedNilCloningOperation struct{}

func (typedNilCloningOperation) Op() internal.OpType { return internal.OpAddType }
func (typedNilCloningOperation) Path() []string      { return []string{"name"} }
func (typedNilCloningOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (typedNilCloningOperation) Validate() error { return nil }
func (typedNilCloningOperation) Clone() (internal.Op, error) {
	var cloned *typedNilOperation
	return cloned, nil
}

type cloneBecomesInvalidOperation struct {
	invalid bool
}

func (cloneBecomesInvalidOperation) Op() internal.OpType { return internal.OpAddType }
func (cloneBecomesInvalidOperation) Path() []string      { return []string{"name"} }
func (cloneBecomesInvalidOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (o cloneBecomesInvalidOperation) Validate() error {
	if o.invalid {
		return errCloneInvalid
	}
	return nil
}
func (cloneBecomesInvalidOperation) Clone() (internal.Op, error) {
	return cloneBecomesInvalidOperation{invalid: true}, nil
}

type cloneBecomesValidOperation struct{}

func (cloneBecomesValidOperation) Op() internal.OpType { return internal.OpAddType }
func (cloneBecomesValidOperation) Path() []string      { return []string{"name"} }
func (cloneBecomesValidOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (cloneBecomesValidOperation) Validate() error { return errCloneInvalid }
func (cloneBecomesValidOperation) Clone() (internal.Op, error) {
	return validClonedOperation{}, nil
}

type validClonedOperation struct{}

func (validClonedOperation) Op() internal.OpType { return internal.OpAddType }
func (validClonedOperation) Path() []string      { return []string{"name"} }
func (validClonedOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (validClonedOperation) Validate() error { return nil }
func (validClonedOperation) Clone() (internal.Op, error) {
	return validClonedOperation{}, nil
}

type cloneBecomesUnsupportedOperation struct{}

func (cloneBecomesUnsupportedOperation) Op() internal.OpType { return internal.OpAddType }
func (cloneBecomesUnsupportedOperation) Path() []string      { return []string{"count"} }
func (cloneBecomesUnsupportedOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (cloneBecomesUnsupportedOperation) Validate() error { return nil }
func (cloneBecomesUnsupportedOperation) Clone() (internal.Op, error) {
	return unsupportedClonedOperation{}, nil
}

type unsupportedClonedOperation struct{}

func (unsupportedClonedOperation) Op() internal.OpType { return internal.OpIncType }
func (unsupportedClonedOperation) Path() []string      { return []string{"count"} }
func (unsupportedClonedOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (unsupportedClonedOperation) Validate() error { return nil }
func (unsupportedClonedOperation) Clone() (internal.Op, error) {
	return unsupportedClonedOperation{}, nil
}
