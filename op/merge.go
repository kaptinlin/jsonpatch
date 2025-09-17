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
	// Following TypeScript reference logic
	parent, key, err := navigateToParent(doc, op.Path())
	if err != nil {
		if len(op.Path()) == 0 {
			// Root level - check if it's an array
			if slice, ok := doc.([]interface{}); ok {
				return op.mergeRootArray(slice)
			}
			return internal.OpResult[any]{}, ErrNotAnArray
		}
		return internal.OpResult[any]{}, err
	}

	// Must be array reference (like TypeScript)
	slice, ok := parent.([]interface{})
	if !ok {
		return internal.OpResult[any]{}, ErrNotAnArray
	}

	index, ok := key.(int)
	if !ok {
		return internal.OpResult[any]{}, ErrInvalidIndex
	}

	// TypeScript: if (ref.key <= 0) throw new Error('INVALID_KEY');
	if index <= 0 {
		return internal.OpResult[any]{}, ErrInvalidIndex
	}

	if index >= len(slice) {
		return internal.OpResult[any]{}, ErrIndexOutOfRange
	}

	// Get elements to merge (pos-1 and pos)
	one := slice[index-1]
	two := slice[index]
	merged := op.mergeElements(one, two)

	// Create new array with merged result
	newSlice := make([]interface{}, len(slice)-1)
	copy(newSlice[:index-1], slice[:index-1])
	newSlice[index-1] = merged
	copy(newSlice[index:], slice[index+1:])

	// Update parent
	parentPath := op.Path()[:len(op.Path())-1]
	if len(parentPath) == 0 {
		// Root array
		return internal.OpResult[any]{Doc: newSlice, Old: []interface{}{one, two}}, nil
	}

	err = setValueAtPath(doc, parentPath, newSlice)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: []interface{}{one, two}}, nil
}

func (op *MergeOperation) mergeRootArray(slice []interface{}) (internal.OpResult[any], error) {
	pos := int(op.Pos)
	if pos <= 0 || pos >= len(slice) {
		return internal.OpResult[any]{}, ErrInvalidIndex
	}

	one := slice[pos-1]
	two := slice[pos]
	merged := op.mergeElements(one, two)

	// Create new array
	newSlice := make([]interface{}, len(slice)-1)
	copy(newSlice[:pos-1], slice[:pos-1])
	newSlice[pos-1] = merged
	copy(newSlice[pos:], slice[pos+1:])

	return internal.OpResult[any]{Doc: newSlice, Old: []interface{}{one, two}}, nil
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
		"op":   string(internal.OpMergeType),
		"path": formatPath(op.Path()),
		"pos":  op.Pos,
	}
	if len(op.Props) > 0 {
		result["props"] = op.Props
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
