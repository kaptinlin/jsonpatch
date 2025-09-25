package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/pkg/slate"
)

// MergeOperation represents an array merge operation.
// path: target path
// pos: merge position (array index)
// props: properties to apply after merge (can be nil)
// Only supports array type fields.
type MergeOperation struct {
	BaseOp
	Pos   float64                `json:"pos"`   // Merge position
	Props map[string]interface{} `json:"props"` // Properties to apply after merge
}

// NewOpMergeOperation creates a new array merge operation.
func NewOpMergeOperation(path []string, pos float64, props map[string]interface{}) *MergeOperation {
	return &MergeOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Props:  props,
	}
}

// Op returns the operation type.
func (op *MergeOperation) Op() internal.OpType {
	return internal.OpMergeType
}

// Code returns the operation code.
func (op *MergeOperation) Code() int {
	return internal.OpMergeCode
}

// Path returns the operation path.
func (op *MergeOperation) Path() []string {
	return op.path
}

// Apply applies the merge operation following TypeScript reference.
func (op *MergeOperation) Apply(doc any) (internal.OpResult[any], error) {
	// TypeScript reference: merge works on arrays directly using pos parameter
	var targetArray []interface{}

	if len(op.Path()) == 0 {
		// Root level array
		slice, ok := doc.([]interface{})
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidTarget
		}
		targetArray = slice
	} else {
		// Get array at path
		target, err := getValue(doc, op.Path())
		if err != nil {
			return internal.OpResult[any]{}, err
		}
		slice, ok := target.([]interface{})
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidTarget
		}
		targetArray = slice
	}

	pos := int(op.Pos)
	// TypeScript: if (ref.key <= 0) throw new Error('INVALID_KEY');
	if pos <= 0 || pos >= len(targetArray) {
		return internal.OpResult[any]{}, ErrInvalidIndex
	}

	// Get elements to merge (pos-1 and pos)
	one := targetArray[pos-1]
	two := targetArray[pos]
	merged := op.mergeElements(one, two)

	// Create new array with merged result
	newSlice := make([]interface{}, len(targetArray)-1)
	copy(newSlice[:pos-1], targetArray[:pos-1])
	newSlice[pos-1] = merged
	copy(newSlice[pos:], targetArray[pos+1:])

	// Update the document
	if len(op.Path()) == 0 {
		// Root array
		return internal.OpResult[any]{Doc: newSlice, Old: []interface{}{one, two}}, nil
	}

	err := setValueAtPath(doc, op.Path(), newSlice)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: []interface{}{one, two}}, nil
}

// mergeElements merges two elements based on their internal.
func (op *MergeOperation) mergeElements(one, two interface{}) interface{} {
	// String concatenation
	if strOne, ok := one.(string); ok {
		if strTwo, ok := two.(string); ok {
			return strOne + strTwo
		}
	}

	// Number addition
	if numOne, ok := one.(float64); ok {
		if numTwo, ok := two.(float64); ok {
			return numOne + numTwo
		}
	}
	if numOne, ok := one.(int); ok {
		if numTwo, ok := two.(int); ok {
			return float64(numOne + numTwo)
		}
	}
	// Handle mixed int/float64 cases
	if numOne, ok := one.(int); ok {
		if numTwo, ok := two.(float64); ok {
			return float64(numOne) + numTwo
		}
	}
	if numOne, ok := one.(float64); ok {
		if numTwo, ok := two.(int); ok {
			return numOne + float64(numTwo)
		}
	}

	// Slate-like text node merging
	if slate.IsTextNode(one) && slate.IsTextNode(two) {
		merged := slate.MergeTextNodesFromMaps(one.(map[string]interface{}), two.(map[string]interface{}))
		// Apply props if specified
		if op.Props != nil {
			for k, v := range op.Props {
				merged[k] = v
			}
		}
		return merged
	}

	// Slate-like element node merging
	if slate.IsElementNode(one) && slate.IsElementNode(two) {
		merged := slate.MergeElementNodesFromMaps(one.(map[string]interface{}), two.(map[string]interface{}))
		// Apply props if specified
		if op.Props != nil {
			for k, v := range op.Props {
				merged[k] = v
			}
		}
		return merged
	}

	// Default: return array of both elements
	return []interface{}{one, two}
}

// Old methods removed - now using pkg/slate functions

// ToJSON serializes the operation to JSON format.
func (op *MergeOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpMergeType),
		Path: formatPath(op.Path()),
		Pos:  int(op.Pos),
	}
	if len(op.Props) > 0 {
		result.Props = op.Props
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *MergeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMergeCode, op.Path(), op.Props}, nil
}

// Validate validates the merge operation.
func (op *MergeOperation) Validate() error {
	if op.Pos < 0 {
		return ErrPositionNegative
	}
	return nil
}

// Short aliases for common use
var (
	// NewMerge creates a new merge operation
	NewMerge = NewOpMergeOperation
)
