package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// FlipOperation represents a flip operation that inverts boolean values
// or converts other types to boolean and then inverts them.
type FlipOperation struct {
	BaseOp
}

// NewFlip creates a new flip operation.
func NewFlip(path []string) *FlipOperation {
	return &FlipOperation{
		BaseOp: NewBaseOp(path),
	}
}

// Op returns the operation type.
func (f *FlipOperation) Op() internal.OpType {
	return internal.OpFlipType
}

// Apply applies the flip operation to the document.
func (f *FlipOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(f.Path()) == 0 {
		flipped := flipValue(doc)
		return internal.OpResult[any]{Doc: flipped, Old: doc}, nil
	}

	value, err := value(doc, f.Path())
	if err != nil {
		value = nil
	}
	oldValue := value
	flipped := flipValue(value)

	err = setValueAtPath(doc, f.Path(), flipped)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// Uses JavaScript truthiness rules for non-boolean types.
func flipValue(value any) any {
	switch v := value.(type) {
	case nil:
		return true
	case bool:
		return !v
	case float64:
		return v == 0
	case int:
		return v == 0
	case string:
		return v == ""
	default:
		return false
	}
}

// Validate validates the flip operation.
func (f *FlipOperation) Validate() error {
	return nil
}
