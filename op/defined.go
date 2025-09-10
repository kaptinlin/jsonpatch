package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpDefinedOperation represents a test operation that checks if a path is defined.
type OpDefinedOperation struct {
	BaseOp
}

// NewOpDefinedOperation creates a new OpDefinedOperation operation.
func NewOpDefinedOperation(path []string) *OpDefinedOperation {
	return &OpDefinedOperation{
		BaseOp: NewBaseOp(path),
	}
}

// Op returns the operation type.
func (o *OpDefinedOperation) Op() internal.OpType {
	return internal.OpDefinedType
}

// Code returns the operation code.
func (o *OpDefinedOperation) Code() int {
	return internal.OpDefinedCode
}

// checkPathExists is a helper function that checks if a path exists
func (o *OpDefinedOperation) checkPathExists(doc interface{}) bool {
	_, err := getValue(doc, o.path)
	return err == nil
}

// Test performs the defined operation.
func (o *OpDefinedOperation) Test(doc interface{}) (bool, error) {
	// Direct return without intermediate variables
	return o.checkPathExists(doc), nil
}

// Not returns false (defined operation doesn't support not modifier).
func (o *OpDefinedOperation) Not() bool {
	return false
}

// Apply applies the defined operation.
func (o *OpDefinedOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Use the same logic as Test but avoid double call
	if !o.checkPathExists(doc) {
		return internal.OpResult[any]{}, ErrDefinedTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *OpDefinedOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":   string(internal.OpDefinedType),
		"path": formatPath(o.Path()),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpDefinedOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpDefinedCode, o.Path()}, nil
}

// Validate validates the defined operation.
func (o *OpDefinedOperation) Validate() error {
	// Empty path (root) is valid for defined operation
	return nil
}

// Path returns the path for the defined operation.
func (o *OpDefinedOperation) Path() []string {
	return o.path
}
