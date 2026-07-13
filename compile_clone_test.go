package jsonpatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

var errCloneFailed = errors.New("clone failed")

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

type cloneFailingOperation struct{}

func (cloneFailingOperation) Op() internal.OpType { return internal.OpAddType }
func (cloneFailingOperation) Path() []string      { return []string{"name"} }
func (cloneFailingOperation) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (cloneFailingOperation) Validate() error             { return nil }
func (cloneFailingOperation) Clone() (internal.Op, error) { return nil, errCloneFailed }
