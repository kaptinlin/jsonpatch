package op

import "github.com/kaptinlin/jsonpatch/internal"

// NotOperation represents a logical NOT operation that negates predicates.
type NotOperation struct {
	BaseOp
	Operations []any `json:"apply"` // Array of operations to apply (then negate)
}

// NewNot creates a new NOT operation.
func NewNot(operand internal.PredicateOp) *NotOperation {
	var path []string
	if operand != nil {
		path = operand.Path()
	}
	return &NotOperation{
		BaseOp:     NewBaseOp(path),
		Operations: []any{operand},
	}
}

// NewNotMultiple creates a new NOT operation with multiple operands.
func NewNotMultiple(path []string, ops []any) *NotOperation {
	return &NotOperation{
		BaseOp:     NewBaseOp(path),
		Operations: ops,
	}
}

// Op returns the operation type.
func (n *NotOperation) Op() internal.OpType {
	return internal.OpNotType
}

// Code returns the operation code.
func (n *NotOperation) Code() int {
	return internal.OpNotCode
}

// Test evaluates the NOT predicate condition.
func (n *NotOperation) Test(doc any) (bool, error) {
	// NOT operation: returns true if ALL operands are false
	for _, op := range n.Operations {
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
func (n *NotOperation) Not() bool {
	return true
}

// Apply applies the NOT operation.
func (n *NotOperation) Apply(doc any) (internal.OpResult[any], error) {
	ok, err := n.Test(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	if !ok {
		return internal.OpResult[any]{}, ErrNotTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (n *NotOperation) ToJSON() (internal.Operation, error) {
	opsJSON, err := predicateOpsToJSON(n.Operations, ErrInvalidPredicateInNot)
	if err != nil {
		return internal.Operation{}, err
	}

	return internal.Operation{
		Op:    string(internal.OpNotType),
		Path:  formatPath(n.Path()),
		Apply: opsJSON,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (n *NotOperation) ToCompact() (internal.CompactOperation, error) {
	opsCompact, err := predicateOpsToCompact(n.Operations, ErrInvalidPredicateInNot)
	if err != nil {
		return nil, err
	}

	return internal.CompactOperation{internal.OpNotCode, n.Path(), opsCompact}, nil
}

// Ops returns the operand operations.
func (n *NotOperation) Ops() []internal.PredicateOp {
	return extractPredicateOps(n.Operations)
}

// Validate validates the NOT operation.
func (n *NotOperation) Validate() error {
	if len(n.Operations) == 0 {
		return ErrNotNoOperands
	}
	return validatePredicateOps(n.Operations, ErrInvalidPredicateInNot)
}
