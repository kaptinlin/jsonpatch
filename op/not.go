package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpNotOperation represents a logical NOT operation that negates predicates.
type OpNotOperation struct {
	BaseOp
	Operations []interface{} `json:"ops"` // Array of operations to negate
}

// NewOpNotOperation creates a new NOT operation.
func NewOpNotOperation(operand internal.PredicateOp) *OpNotOperation {
	var path []string
	if operand != nil {
		path = operand.Path()
	}
	return &OpNotOperation{
		BaseOp:     NewBaseOp(path),
		Operations: []interface{}{operand},
	}
}

// NewOpNotOperationMultiple creates a new NOT operation with multiple operands.
func NewOpNotOperationMultiple(path []string, ops []interface{}) *OpNotOperation {
	return &OpNotOperation{
		BaseOp:     NewBaseOp(path),
		Operations: ops,
	}
}

// Op returns the operation type.
func (o *OpNotOperation) Op() internal.OpType {
	return internal.OpNotType
}

// Code returns the operation code.
func (o *OpNotOperation) Code() int {
	return internal.OpNotCode
}

// Path returns the operation path.
func (o *OpNotOperation) Path() []string {
	return o.path
}

// Test evaluates the NOT predicate condition.
func (o *OpNotOperation) Test(doc any) (bool, error) {
	// NOT operation: returns true if ALL operands are false
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return false, ErrInvalidPredicateInNot
		}
		result, err := predicateOp.Test(doc)
		if err != nil {
			// For NOT operations, an error in the operand (like path not found)
			// should be treated as the operand returning false, continue to next
			continue
		}
		if result {
			// If any operand is true, NOT returns false
			return false, nil
		}
	}
	// All operands are false, NOT returns true
	return true, nil
}

// Not returns true since this is a NOT operation.
func (o *OpNotOperation) Not() bool {
	return true
}

// Apply applies the NOT operation.
func (o *OpNotOperation) Apply(doc any) (internal.OpResult[any], error) {
	ok, err := o.Test(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	if !ok {
		return internal.OpResult[any]{}, ErrNotTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *OpNotOperation) ToJSON() (internal.Operation, error) {
	opsJSON := make([]interface{}, 0, len(o.Operations))
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, ErrInvalidPredicateInNot
		}
		jsonVal, err := predicateOp.ToJSON()
		if err != nil {
			return nil, err
		}
		opsJSON = append(opsJSON, jsonVal)
	}

	return internal.Operation{
		"op":    string(internal.OpNotType),
		"path":  formatPath(o.Path()),
		"apply": opsJSON,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpNotOperation) ToCompact() (internal.CompactOperation, error) {
	opsCompact := make([]interface{}, 0, len(o.Operations))
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, ErrInvalidPredicateInNot
		}
		compact, err := predicateOp.ToCompact()
		if err != nil {
			return nil, err
		}
		opsCompact = append(opsCompact, compact)
	}

	return internal.CompactOperation{internal.OpNotCode, o.Path(), opsCompact}, nil
}

// Ops returns the operand operations.
func (o *OpNotOperation) Ops() []internal.PredicateOp {
	ops := make([]internal.PredicateOp, 0, len(o.Operations))
	for _, op := range o.Operations {
		if predicateOp, ok := op.(internal.PredicateOp); ok {
			ops = append(ops, predicateOp)
		}
	}
	return ops
}

// Validate validates the NOT operation.
func (o *OpNotOperation) Validate() error {
	if len(o.Operations) == 0 {
		return ErrNotNoOperands
	}
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return ErrInvalidPredicateInNot
		}
		if err := predicateOp.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Short aliases for common use
var (
	// NewNot creates a new not operation
	NewNot = NewOpNotOperation
	// NewNotMultiple creates a new not operation with multiple operands
	NewNotMultiple = NewOpNotOperationMultiple
)
