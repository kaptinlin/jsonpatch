package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// AndOperation represents an and operation that combines multiple predicate operations.
type AndOperation struct {
	BaseOp
	Operations []interface{} `json:"apply"` // Array of operations to apply
}

type OpAndOperation = AndOperation //nolint:revive // Backward compatibility alias

// NewOpAndOperation creates a new AndOperation operation.
func NewOpAndOperation(path []string, ops []interface{}) *AndOperation {
	return &AndOperation{
		BaseOp:     NewBaseOp(path),
		Operations: ops,
	}
}

// Op returns the operation type.
func (o *AndOperation) Op() internal.OpType {
	return internal.OpAndType
}

// Code returns the operation code.
func (o *AndOperation) Code() int {
	return internal.OpAndCode
}

// Ops returns the predicate operations.
func (o *AndOperation) Ops() []internal.PredicateOp {
	ops := make([]internal.PredicateOp, 0, len(o.Operations))
	for _, op := range o.Operations {
		if predicateOp, ok := op.(internal.PredicateOp); ok {
			ops = append(ops, predicateOp)
		}
	}
	return ops
}

// Test performs the AND operation.
func (o *AndOperation) Test(doc interface{}) (bool, error) {
	// If no operations, return false (empty AND is false)
	if len(o.Operations) == 0 {
		result := false
		return result, nil
	}

	// Test all operations
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return false, ErrInvalidPredicateInAnd
		}
		ok, err := predicateOp.Test(doc)
		if err != nil {
			return false, err
		}
		if !ok {
			// If any operation fails, AND fails
			return false, nil
		}
	}

	// All operations passed
	return true, nil
}

// Apply applies the AND operation.
func (o *AndOperation) Apply(doc any) (internal.OpResult[any], error) {
	ok, err := o.Test(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	if !ok {
		return internal.OpResult[any]{}, ErrAndTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *AndOperation) ToJSON() (internal.Operation, error) {
	opsJSON := make([]interface{}, 0, len(o.Operations))
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, ErrInvalidPredicateInAnd
		}
		jsonVal, err := predicateOp.ToJSON()
		if err != nil {
			return nil, err
		}
		opsJSON = append(opsJSON, jsonVal)
	}
	return internal.Operation{
		"op":    string(internal.OpAndType),
		"path":  formatPath(o.path),
		"apply": opsJSON,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *AndOperation) ToCompact() (internal.CompactOperation, error) {
	opsCompact := make([]interface{}, 0, len(o.Operations))
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, ErrInvalidPredicateInAnd
		}
		compact, err := predicateOp.ToCompact()
		if err != nil {
			return nil, err
		}
		opsCompact = append(opsCompact, compact)
	}
	return internal.CompactOperation{internal.OpAndCode, o.path, opsCompact}, nil
}

// Validate validates the AND operation.
func (o *AndOperation) Validate() error {
	if len(o.Operations) == 0 {
		return ErrAndNoOperands
	}
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return ErrInvalidPredicateInAnd
		}
		if err := predicateOp.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the path for the AND operation.
func (o *AndOperation) Path() []string {
	return o.path
}

// Not returns false since this is not a NOT operation.
func (o *AndOperation) Not() bool {
	return false
}

// Short aliases for common use
var (
	// NewAnd creates a new and operation
	NewAnd = NewOpAndOperation
)
