package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// SplitOperation represents a string split operation.
// path: target path
// pos: split position (rune index)
// props: properties to apply after split (can be nil)
// Only supports string type fields.
type SplitOperation struct {
	BaseOp
	Pos   float64     `json:"pos"`   // Split position
	Props interface{} `json:"props"` // Properties to apply after split
}

// NewOpSplitOperation creates a new string split operation.
func NewOpSplitOperation(path []string, pos float64, props interface{}) *SplitOperation {
	return &SplitOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Props:  props,
	}
}

// Op returns the operation type.
func (op *SplitOperation) Op() internal.OpType {
	return internal.OpSplitType
}

// Code returns the operation code.
func (op *SplitOperation) Code() int {
	return internal.OpSplitCode
}

// Path returns the operation path.
func (op *SplitOperation) Path() []string {
	return op.path
}

// Apply applies the split operation following TypeScript reference.
func (op *SplitOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get the target value
	var target interface{}
	var err error

	if len(op.Path()) == 0 {
		target = doc
	} else {
		target, err = getValue(doc, op.Path())
		if err != nil {
			return internal.OpResult[any]{}, err
		}
	}

	// Split the value following TypeScript logic
	parts := op.splitValue(target)

	// Following TypeScript reference behavior
	if len(op.Path()) == 0 {
		// Root level split - return the split result as new document
		return internal.OpResult[any]{Doc: parts, Old: target}, nil
	}

	// For array elements, follow TypeScript pattern: replace element and insert new one
	parent, key, err := navigateToParent(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if slice, ok := parent.([]interface{}); ok {
		if index, ok := key.(int); ok {
			// TypeScript: ref.obj[ref.key] = tuple[0]; ref.obj.splice(ref.key + 1, 0, tuple[1]);
			splitResult := parts.([]interface{})

			// Create new array with split
			newSlice := make([]interface{}, len(slice)+1)
			copy(newSlice[:index], slice[:index])
			newSlice[index] = splitResult[0]
			newSlice[index+1] = splitResult[1]
			copy(newSlice[index+2:], slice[index+1:])

			// Update parent
			parentPath := op.Path()[:len(op.Path())-1]
			if len(parentPath) == 0 {
				// Root array - return new array
				return internal.OpResult[any]{Doc: newSlice, Old: target}, nil
			}
			err = setValueAtPath(doc, parentPath, newSlice)
			if err != nil {
				return internal.OpResult[any]{}, err
			}
		}
	} else {
		// For objects, replace the value with split result
		err = setValueAtPath(doc, op.Path(), parts)
		if err != nil {
			return internal.OpResult[any]{}, err
		}
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// splitValue splits a value based on its type
func (op *SplitOperation) splitValue(value interface{}) interface{} {
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
		if isSlateTextNode(v) {
			propsMap, _ := op.Props.(map[string]interface{})
			results := splitSlateTextNode(v, int(op.Pos), propsMap)
			if results != nil {
				return []interface{}{results[0], results[1]}
			}
		}
		// Check if it's a Slate-like element node with children
		if isSlateElementNode(v) {
			propsMap, _ := op.Props.(map[string]interface{})
			results := splitSlateElementNode(v, int(op.Pos), propsMap)
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
func (op *SplitOperation) splitString(s string) []interface{} {
	runes := []rune(s)
	// High-performance type conversion (single, boundary conversion)
	pos := int(op.Pos) // Already validated as safe integer

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
func (op *SplitOperation) splitNumber(n float64) []interface{} {
	pos := op.Pos // Already validated as safe number
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
func (op *SplitOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpSplitType),
		Path: formatPath(op.Path()),
		Pos:  int(op.Pos),
	}
	if op.Props != nil {
		if props, ok := op.Props.(map[string]any); ok {
			result.Props = props
		}
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *SplitOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpSplitCode, op.Path(), op.Pos, op.Props}, nil
}

// Validate validates the split operation.
func (op *SplitOperation) Validate() error {
	// Empty path is valid for split operation (root level)
	// Position bounds are checked in Apply method
	return nil
}

// Short aliases for common use
var (
	// NewSplit creates a new split operation
	NewSplit = NewOpSplitOperation
)

// splitSlateTextNode splits a Slate text node at the specified position
func splitSlateTextNode(nodeMap map[string]interface{}, pos int, props map[string]interface{}) []map[string]interface{} {
	text, ok := nodeMap["text"].(string)
	if !ok {
		return nil
	}

	runes := []rune(text)

	// Clamp position to valid bounds
	if pos < 0 {
		pos = 0
	}
	if pos > len(runes) {
		pos = len(runes)
	}

	before := string(runes[:pos])
	after := string(runes[pos:])

	// Create two new nodes with inherited properties
	beforeNode := make(map[string]interface{})
	afterNode := make(map[string]interface{})

	// Copy properties from original node
	for k, v := range nodeMap {
		if k != "text" {
			beforeNode[k] = v
			afterNode[k] = v
		}
	}

	beforeNode["text"] = before
	afterNode["text"] = after

	// Apply extra properties if specified
	for k, v := range props {
		beforeNode[k] = v
		afterNode[k] = v
	}

	return []map[string]interface{}{beforeNode, afterNode}
}

// splitSlateElementNode splits a Slate element node at the specified position in its children
func splitSlateElementNode(nodeMap map[string]interface{}, pos int, props map[string]interface{}) []map[string]interface{} {
	children, ok := nodeMap["children"].([]interface{})
	if !ok {
		return nil
	}

	// Clamp position to valid bounds
	if pos < 0 {
		pos = 0
	}
	if pos > len(children) {
		pos = len(children)
	}

	beforeChildren := children[:pos]
	afterChildren := children[pos:]

	// Create two new nodes with inherited properties
	beforeNode := make(map[string]interface{})
	afterNode := make(map[string]interface{})

	// Copy properties from original node
	for k, v := range nodeMap {
		if k != "children" {
			beforeNode[k] = v
			afterNode[k] = v
		}
	}

	// Copy children slices to avoid mutation
	beforeChildrenCopy := make([]interface{}, len(beforeChildren))
	copy(beforeChildrenCopy, beforeChildren)
	afterChildrenCopy := make([]interface{}, len(afterChildren))
	copy(afterChildrenCopy, afterChildren)

	beforeNode["children"] = beforeChildrenCopy
	afterNode["children"] = afterChildrenCopy

	// Apply extra properties if specified
	for k, v := range props {
		beforeNode[k] = v
		afterNode[k] = v
	}

	return []map[string]interface{}{beforeNode, afterNode}
}
