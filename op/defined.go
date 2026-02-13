package op

import "github.com/kaptinlin/jsonpatch/internal"

// DefinedOperation represents a test operation that checks if a path is defined.
type DefinedOperation struct {
	BaseOp
}

// NewDefined creates a new defined operation.
func NewDefined(path []string) *DefinedOperation {
	return &DefinedOperation{
		BaseOp: NewBaseOp(path),
	}
}

// Op returns the operation type.
func (d *DefinedOperation) Op() internal.OpType {
	return internal.OpDefinedType
}

// Code returns the operation code.
func (d *DefinedOperation) Code() int {
	return internal.OpDefinedCode
}

// Test performs the defined operation.
func (d *DefinedOperation) Test(doc any) (bool, error) {
	return pathExists(doc, d.path), nil
}

// Apply applies the defined operation.
func (d *DefinedOperation) Apply(doc any) (internal.OpResult[any], error) {
	if !pathExists(doc, d.path) {
		return internal.OpResult[any]{}, ErrDefinedTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (d *DefinedOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpDefinedType),
		Path: formatPath(d.Path()),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (d *DefinedOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpDefinedCode, d.Path()}, nil
}

// Validate validates the defined operation.
func (d *DefinedOperation) Validate() error {
	// Empty path (root) is valid for defined operation
	return nil
}
