package op

import "github.com/kaptinlin/jsonpatch/internal"

// NotOperation represents a logical NOT operation that negates predicates.
type NotOperation struct {
	BaseOp
	Operand PredicateOp `json:"-"`
}

// NewNot creates a new NOT operation.
func NewNot(path []string, operand PredicateOp) *NotOperation {
	return &NotOperation{
		BaseOp:  NewBaseOp(path),
		Operand: operand,
	}
}

// Op returns the operation type.
func (n *NotOperation) Op() internal.OpType {
	return internal.OpNotType
}

// Test evaluates the NOT predicate condition.
func (n *NotOperation) Test(doc any) (bool, error) {
	operand, err := n.operand()
	if err != nil {
		return false, err
	}
	result, err := operand.Test(doc)
	if err != nil {
		return false, err
	}
	return !result, nil
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
		return internal.OpResult[any]{}, predicateFailure(ErrNotTestFailed)
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// Ops returns the operand operations.
func (n *NotOperation) Ops() []internal.PredicateOp {
	return []internal.PredicateOp{n.Operand}
}

// Validate validates the NOT operation.
func (n *NotOperation) Validate() error {
	operand, err := n.operand()
	if err != nil {
		return err
	}
	return operand.Validate()
}

func (n *NotOperation) operand() (internal.PredicateOp, error) {
	if isNilPredicate(n.Operand) {
		return nil, ErrInvalidPredicateInNot
	}
	return n.Operand, nil
}
