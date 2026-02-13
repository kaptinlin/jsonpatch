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

// Code returns the operation code.
func (u *UndefinedOperation) Code() int {
	return internal.OpUndefinedCode
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

// ToJSON serializes the operation to JSON format.
func (u *UndefinedOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpUndefinedType),
		Path: formatPath(u.path),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (u *UndefinedOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpUndefinedCode, u.path}, nil
}

// Validate validates the undefined operation.
func (u *UndefinedOperation) Validate() error {
	if len(u.path) == 0 {
		return ErrPathEmpty
	}
	return nil
}
