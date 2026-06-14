package op

import "github.com/kaptinlin/jsonpatch/internal"

// UndefinedOperation represents an undefined operation that checks if a path doesn't exist.
type UndefinedOperation struct {
	BaseOp
}

// NewUndefined creates a new undefined operation.
func NewUndefined(path []string) *UndefinedOperation {
	return &UndefinedOperation{
		BaseOp: NewBaseOp(path),
	}
}

// Op returns the operation type.
func (u *UndefinedOperation) Op() internal.OpType {
	return internal.OpUndefinedType
}

// Test performs the undefined operation.
func (u *UndefinedOperation) Test(doc any) (bool, error) {
	return !pathExists(doc, u.path), nil
}

// Apply applies the undefined operation.
func (u *UndefinedOperation) Apply(doc any) (internal.OpResult[any], error) {
	if pathExists(doc, u.path) {
		return internal.OpResult[any]{}, ErrUndefinedTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// Validate validates the undefined operation.
func (u *UndefinedOperation) Validate() error {
	// Empty path (root) is valid for undefined operation, symmetric with defined
	return nil
}
