package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/pkg/slate"
)

// OpMergeOperation represents an array merge operation.
// path: target path
// pos: merge position (array index)
// props: properties to apply after merge (can be nil)
// Only supports array type fields.
type OpMergeOperation struct {
	BaseOp
	Pos   int                    `json:"pos"`   // Merge position
	Props map[string]interface{} `json:"props"` // Properties to apply after merge
}

// NewOpMergeOperation creates a new array merge operation.
func NewOpMergeOperation(path []string, pos int, props map[string]interface{}) *OpMergeOperation {
	return &OpMergeOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Props:  props,
	}
}

// Op returns the operation type.
func (op *OpMergeOperation) Op() internal.OpType {
	return internal.OpMergeType
}

// Code returns the operation code.
func (op *OpMergeOperation) Code() int {
	return internal.OpMergeCode
}

// Path returns the operation path.
func (op *OpMergeOperation) Path() []string {
	return op.path
}

// Apply applies the array merge operation.
func (op *OpMergeOperation) Apply(doc any) (internal.OpResult, error) {
	// Handle root level merge specially
	if len(op.Path()) == 0 {
		// Target is the root document
		targetArr, ok := doc.([]interface{})
		if !ok {
			return internal.OpResult{}, ErrNotAnArray
		}

		// Check if we have at least 2 elements and pos is valid
		if len(targetArr) < 2 {
			return internal.OpResult{}, ErrArrayTooSmall
		}
		if op.Pos <= 0 || op.Pos >= len(targetArr) {
			return internal.OpResult{}, ErrPositionOutOfBounds
		}

		// Clone the target array
		clonedTarget, err := DeepClone(targetArr)
		if err != nil {
			return internal.OpResult{}, err
		}

		// Merge array elements
		mergedArr, oldElements := op.mergeArrayElements(clonedTarget.([]interface{}), op.Pos)

		return internal.OpResult{Doc: mergedArr, Old: oldElements}, nil
	}

	// Get the target value for non-root paths
	target, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult{}, err
	}

	// Check if target is an array
	targetArr, ok := target.([]interface{})
	if !ok {
		return internal.OpResult{}, ErrNotAnArray
	}

	// Check if we have at least 2 elements and pos is valid
	if len(targetArr) < 2 {
		return internal.OpResult{}, ErrArrayTooSmall
	}
	if op.Pos <= 0 || op.Pos >= len(targetArr) {
		return internal.OpResult{}, ErrPositionOutOfBounds
	}

	// Clone the target array
	clonedTarget, err := DeepClone(targetArr)
	if err != nil {
		return internal.OpResult{}, err
	}

	// Merge array elements
	mergedArr, oldElements := op.mergeArrayElements(clonedTarget.([]interface{}), op.Pos)

	// Set the merged array back
	err = setValueAtPath(doc, op.Path(), mergedArr)
	if err != nil {
		return internal.OpResult{}, err
	}

	return internal.OpResult{Doc: doc, Old: oldElements}, nil
}

// mergeArrayElements merges array elements at position pos-1 and pos (like TypeScript version).
// Returns the new array and the old elements that were merged.
func (op *OpMergeOperation) mergeArrayElements(arr []interface{}, pos int) ([]interface{}, []interface{}) {
	if pos <= 0 || pos >= len(arr) {
		return arr, nil
	}

	// Get elements to merge (pos-1 and pos, like TypeScript version)
	one := arr[pos-1]
	two := arr[pos]
	oldElements := []interface{}{one, two}

	// Merge the elements
	merged := op.mergeElements(one, two)

	// Create new array with merged element
	result := make([]interface{}, len(arr)-1)
	copy(result[:pos-1], arr[:pos-1])
	result[pos-1] = merged
	copy(result[pos:], arr[pos+1:])

	return result, oldElements
}

// mergeElements merges two elements based on their internal.
func (op *OpMergeOperation) mergeElements(one, two interface{}) interface{} {
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
func (op *OpMergeOperation) ToJSON() (internal.Operation, error) {
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
func (op *OpMergeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMergeCode, op.Path(), op.Props}, nil
}

// Validate validates the merge operation.
func (op *OpMergeOperation) Validate() error {
	if op.Pos < 0 {
		return ErrPositionNegative
	}
	return nil
}
