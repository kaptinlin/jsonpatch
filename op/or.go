package op

import (
	"slices"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OrOperation represents an OR operation that combines multiple predicate operations.
type OrOperation struct {
	BaseOp
	Operations []PredicateOp `json:"apply"`
}

// NewOr creates a new OR operation.
func NewOr(path []string, predicates ...PredicateOp) *OrOperation {
	return &OrOperation{
		BaseOp:     NewBaseOp(path),
		Operations: slices.Clone(predicates),
	}
}

// Op returns the operation type.
func (oo *OrOperation) Op() internal.OpType {
	return internal.OpOrType
}

// Ops returns the predicate operations.
func (oo *OrOperation) Ops() []internal.PredicateOp {
	return slices.Clone(oo.Operations)
}

// Test performs the OR operation.
func (oo *OrOperation) Test(doc any) (bool, error) {
	// If no operations, return false (empty OR is false)
	if len(oo.Operations) == 0 {
		return false, nil
	}

	for _, predicate := range oo.Operations {
		if isNilPredicate(predicate) {
			return false, ErrInvalidPredicateInOr
		}
		ok, err := predicate.Test(doc)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

// Apply applies the OR operation to the document.
func (oo *OrOperation) Apply(doc any) (internal.OpResult[any], error) {
	ok, err := oo.Test(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	if !ok {
		return internal.OpResult[any]{}, predicateFailure(ErrOrTestFailed)
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// Validate validates the OR operation.
func (oo *OrOperation) Validate() error {
	return validatePredicateOps(oo.Operations, ErrInvalidPredicateInOr)
}
