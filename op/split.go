package op

import (
	"errors"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/pkg/slate"
)

// OpSplitOperation represents a string split operation.
// path: target path
// pos: split position (rune index)
// props: properties to apply after split (can be nil)
// Only supports string type fields.
type OpSplitOperation struct {
	BaseOp
	Pos   int         `json:"pos"`   // Split position
	Props interface{} `json:"props"` // Properties to apply after split
}

// NewOpSplitOperation creates a new string split operation.
func NewOpSplitOperation(path []string, pos int, props interface{}) *OpSplitOperation {
	return &OpSplitOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Props:  props,
	}
}

// Op returns the operation type.
func (op *OpSplitOperation) Op() internal.OpType {
	return internal.OpSplitType
}

// Code returns the operation code.
func (op *OpSplitOperation) Code() int {
	return internal.OpSplitCode
}

// Path returns the operation path.
func (op *OpSplitOperation) Path() []string {
	return op.path
}

// Apply applies the string split operation.
func (op *OpSplitOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Handle root level split specially
	if len(op.Path()) == 0 {
		// Target is the root document
		parts := op.splitValue(doc)
		return internal.OpResult[any]{Doc: parts, Old: doc}, nil
	}

	// Get the target value for non-root paths
	target, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Split the value
	parts := op.splitValue(target)

	// Check if we're splitting an array element
	if op.isArrayElementPath(doc, op.Path()) {
		// Special handling for array elements - insert the second part as a new element
		err = op.handleArraySplit(doc, parts.([]interface{}))
		if err != nil {
			// Check if this is a root array modification error
			if errors.Is(err, ErrCannotModifyRootArray) {
				// Handle root array splitting specially
				parentPath := op.Path()[:len(op.Path())-1]
				if len(parentPath) == 0 {
					// This is a root array element split
					elementKey := op.Path()[len(op.Path())-1]
					index, parseErr := parseArrayIndex(elementKey)
					if parseErr != nil {
						return internal.OpResult[any]{}, parseErr
					}

					rootArray, ok := doc.([]interface{})
					if !ok {
						return internal.OpResult[any]{}, ErrNotAnArray
					}

					if index < 0 || index >= len(rootArray) {
						return internal.OpResult[any]{}, ErrArrayIndexOutOfBounds
					}

					// Create new array with the split result
					newArray := make([]interface{}, len(rootArray)+1)
					copy(newArray[:index], rootArray[:index])
					newArray[index] = parts.([]interface{})[0]
					newArray[index+1] = parts.([]interface{})[1]
					copy(newArray[index+2:], rootArray[index+1:])

					return internal.OpResult[any]{Doc: newArray, Old: target}, nil
				}
			}
			return internal.OpResult[any]{}, err
		}
	} else {
		// For objects and other cases, replace the entire value
		err = setValueAtPath(doc, op.Path(), parts)
		if err != nil {
			return internal.OpResult[any]{}, err
		}
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// isArrayElementPath checks if the path points to an array element
func (op *OpSplitOperation) isArrayElementPath(doc interface{}, path []string) bool {
	if len(path) == 0 {
		return false
	}

	// Get the parent path and check if it's an array
	parentPath := path[:len(path)-1]
	var parent interface{}
	var err error

	if len(parentPath) == 0 {
		// Parent is the root document
		parent = doc
	} else {
		parent, err = getValue(doc, parentPath)
		if err != nil {
			return false
		}
	}

	_, isArray := parent.([]interface{})
	return isArray
}

// handleArraySplit handles splitting when the target is an array element
func (op *OpSplitOperation) handleArraySplit(doc interface{}, parts []interface{}) error {
	if len(parts) != 2 {
		return ErrArrayTooSmall
	}

	// Get parent array and element index
	parentPath := op.Path()[:len(op.Path())-1]
	elementKey := op.Path()[len(op.Path())-1]

	var parentArray []interface{}
	var ok bool

	if len(parentPath) == 0 {
		// Parent is the root document
		parentArray, ok = doc.([]interface{})
		if !ok {
			return ErrNotAnArray
		}
	} else {
		parent, err := getValue(doc, parentPath)
		if err != nil {
			return err
		}

		parentArray, ok = parent.([]interface{})
		if !ok {
			return ErrNotAnArray
		}
	}

	// Parse the array index
	index, err := parseArrayIndex(elementKey)
	if err != nil {
		return err
	}

	if index < 0 || index >= len(parentArray) {
		return ErrArrayIndexOutOfBounds
	}

	// Create new array with the split result
	newArray := make([]interface{}, len(parentArray)+1)
	copy(newArray[:index], parentArray[:index])
	newArray[index] = parts[0]
	newArray[index+1] = parts[1]
	copy(newArray[index+2:], parentArray[index+1:])

	// Set the new array back
	if len(parentPath) == 0 {
		// We can't modify the root document directly, but we can return an error
		// indicating that the caller should handle this case
		return ErrCannotModifyRootArray
	} else {
		return setValueAtPath(doc, parentPath, newArray)
	}
}

// splitValue splits a value based on its type
func (op *OpSplitOperation) splitValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return op.splitString(v)
	case float64:
		return op.splitNumber(v)
	case int:
		return op.splitNumber(float64(v))
	case bool:
		// For boolean, return a tuple of the same value
		return []interface{}{v, v}
	case map[string]interface{}:
		// Check if it's a Slate-like text node
		if slate.IsTextNode(v) {
			propsMap, _ := op.Props.(map[string]interface{})
			results := slate.SplitTextNodeFromMap(v, op.Pos, propsMap)
			if results != nil {
				return []interface{}{results[0], results[1]}
			}
		}
		// Check if it's a Slate-like element node with children
		if slate.IsElementNode(v) {
			propsMap, _ := op.Props.(map[string]interface{})
			results := slate.SplitElementNodeFromMap(v, op.Pos, propsMap)
			if results != nil {
				return []interface{}{results[0], results[1]}
			}
		}
		// For other objects, return a tuple of the same value
		return []interface{}{v, v}
	default:
		// For unknown types, return a tuple of the same value
		return []interface{}{value, value}
	}
}

// splitString splits a string at the specified position
func (op *OpSplitOperation) splitString(s string) []interface{} {
	runes := []rune(s)
	pos := op.Pos

	// Handle negative positions (count from end)
	if pos < 0 {
		pos = len(runes) + pos
		if pos < 0 {
			pos = 0
		}
	}

	// Handle positions beyond string length
	if pos > len(runes) {
		pos = len(runes)
	}

	before := string(runes[:pos])
	after := string(runes[pos:])

	// If props are specified, wrap in text nodes
	if op.Props != nil {
		if propsMap, ok := op.Props.(map[string]interface{}); ok {
			beforeNode := map[string]interface{}{"text": before}
			afterNode := map[string]interface{}{"text": after}

			// Copy props to both nodes
			for k, v := range propsMap {
				beforeNode[k] = v
				afterNode[k] = v
			}

			return []interface{}{beforeNode, afterNode}
		}
	}

	return []interface{}{before, after}
}

// splitNumber splits a number at the specified position
func (op *OpSplitOperation) splitNumber(n float64) []interface{} {
	pos := float64(op.Pos)
	if pos > n {
		pos = n
	}
	if pos < 0 {
		pos = 0
	}
	return []interface{}{pos, n - pos}
}

// Old Slate-specific split methods removed - now using pkg/slate functions

// ToJSON serializes the operation to JSON format.
func (op *OpSplitOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":   string(internal.OpSplitType),
		"path": formatPath(op.Path()),
		"pos":  op.Pos,
	}
	if op.Props != nil {
		result["props"] = op.Props
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpSplitOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpSplitCode, op.Path(), op.Pos, op.Props}, nil
}

// Validate validates the split operation.
func (op *OpSplitOperation) Validate() error {
	// Empty path is valid for split operation (root level)
	// Position bounds are checked in Apply method
	return nil
}
