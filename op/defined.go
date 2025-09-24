package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// DefinedOperation represents a test operation that checks if a path is defined.
type DefinedOperation struct {
	BaseOp
}

type OpDefinedOperation = DefinedOperation //nolint:revive // Backward compatibility alias

// NewOpDefinedOperation creates a new OpDefinedOperation operation.
func NewOpDefinedOperation(path []string) *DefinedOperation {
	return &DefinedOperation{
		BaseOp: NewBaseOp(path),
	}
}

// Op returns the operation type.
func (o *DefinedOperation) Op() internal.OpType {
	return internal.OpDefinedType
}

// Code returns the operation code.
func (o *DefinedOperation) Code() int {
	return internal.OpDefinedCode
}

// checkPathExists is a helper function that checks if a path exists
func (o *DefinedOperation) checkPathExists(doc interface{}) bool {
	_, err := getValue(doc, o.path)
	return err == nil
}

// Test performs the defined operation.
func (o *DefinedOperation) Test(doc interface{}) (bool, error) {
	// Direct return without intermediate variables
	return o.checkPathExists(doc), nil
}

// Not returns false (defined operation doesn't support not modifier).
func (o *DefinedOperation) Not() bool {
	return false
}

// Apply applies the defined operation.
func (o *DefinedOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Use the same logic as Test but avoid double call
	if !o.checkPathExists(doc) {
		return internal.OpResult[any]{}, ErrDefinedTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *DefinedOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpDefinedType),
		Path: formatPath(o.Path()),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *DefinedOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpDefinedCode, o.Path()}, nil
}

// Validate validates the defined operation.
func (o *DefinedOperation) Validate() error {
	// Empty path (root) is valid for defined operation
	return nil
}

// Path returns the path for the defined operation.
func (o *DefinedOperation) Path() []string {
	return o.path
}

// Short aliases for common use
var (
	// NewDefined creates a new defined operation
	NewDefined = NewOpDefinedOperation
)
