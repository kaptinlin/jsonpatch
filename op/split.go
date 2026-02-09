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
	Pos   float64 `json:"pos"`   // Split position
	Props any     `json:"props"` // Properties to apply after split
}

// NewSplit creates a new split operation.
func NewSplit(path []string, pos float64, props any) *SplitOperation {
	return &SplitOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Props:  props,
	}
}

// Op returns the operation type.
func (o *SplitOperation) Op() internal.OpType {
	return internal.OpSplitType
}

// Code returns the operation code.
func (o *SplitOperation) Code() int {
	return internal.OpSplitCode
}

// Apply applies the split operation following TypeScript reference.
func (o *SplitOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get the target value
	var target any
	var err error

	if len(o.Path()) == 0 {
		target = doc
	} else {
		target, err = getValue(doc, o.Path())
		if err != nil {
			return internal.OpResult[any]{}, err
		}
	}

	// Split the value following TypeScript logic
	parts := o.splitValue(target)

	// Following TypeScript reference behavior
	if len(o.Path()) == 0 {
		// Root level split - return the split result as new document
		return internal.OpResult[any]{Doc: parts, Old: target}, nil
	}

	// For array elements, follow TypeScript pattern: replace element and insert new one
	parent, key, err := navigateToParent(doc, o.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if slice, ok := parent.([]any); ok {
		if index, ok := key.(int); ok {
			// TypeScript: ref.obj[ref.key] = tuple[0]; ref.obj.splice(ref.key + 1, 0, tuple[1]);
			splitResult := parts.([]any)

			// Create new array with split
			newSlice := make([]any, len(slice)+1)
			copy(newSlice[:index], slice[:index])
			newSlice[index] = splitResult[0]
			newSlice[index+1] = splitResult[1]
			copy(newSlice[index+2:], slice[index+1:])

			// Update parent
			parentPath := o.Path()[:len(o.Path())-1]
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
		err = setValueAtPath(doc, o.Path(), parts)
		if err != nil {
			return internal.OpResult[any]{}, err
		}
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// splitValue splits a value based on its type
func (o *SplitOperation) splitValue(value any) any {
	switch v := value.(type) {
	case string:
		return o.splitString(v)
	case float64:
		return o.splitNumber(v)
	case int:
		return o.splitNumber(float64(v))
	case bool:
		// For boolean, return a tuple of the same value
		return []any{v, v}
	case map[string]any:
		// Check if it's a Slate-like text node
		if isSlateTextNode(v) {
			propsMap, _ := o.Props.(map[string]any)
			results := splitSlateTextNode(v, int(o.Pos), propsMap)
			if results != nil {
				return []any{results[0], results[1]}
			}
		}
		// Check if it's a Slate-like element node with children
		if isSlateElementNode(v) {
			propsMap, _ := o.Props.(map[string]any)
			results := splitSlateElementNode(v, int(o.Pos), propsMap)
			if results != nil {
				return []any{results[0], results[1]}
			}
		}
		// For other objects, return a tuple of the same value
		return []any{v, v}
	default:
		// For unknown types, return a tuple of the same value
		return []any{value, value}
	}
}

// splitString splits a string at the specified position
func (o *SplitOperation) splitString(s string) []any {
	runes := []rune(s)
	// High-performance type conversion (single, boundary conversion)
	pos := int(o.Pos) // Already validated as safe integer

	// Handle negative positions (count from end)
	if pos < 0 {
		pos = max(len(runes)+pos, 0)
	}

	// Handle positions beyond string length
	if pos > len(runes) {
		pos = len(runes)
	}

	before := string(runes[:pos])
	after := string(runes[pos:])

	// If props are specified, wrap in text nodes
	if o.Props != nil {
		if propsMap, ok := o.Props.(map[string]any); ok {
			beforeNode := map[string]any{"text": before}
			afterNode := map[string]any{"text": after}

			// Copy props to both nodes
			for k, v := range propsMap {
				beforeNode[k] = v
				afterNode[k] = v
			}

			return []any{beforeNode, afterNode}
		}
	}

	return []any{before, after}
}

// splitNumber splits a number at the specified position
func (o *SplitOperation) splitNumber(n float64) []any {
	pos := o.Pos // Already validated as safe number
	if pos > n {
		pos = n
	}
	if pos < 0 {
		pos = 0
	}
	return []any{pos, n - pos}
}

// Old Slate-specific split methods removed - now using pkg/slate functions

// ToJSON serializes the operation to JSON format.
func (o *SplitOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpSplitType),
		Path: formatPath(o.Path()),
		Pos:  int(o.Pos),
	}
	if o.Props != nil {
		if props, ok := o.Props.(map[string]any); ok {
			result.Props = props
		}
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *SplitOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpSplitCode, o.Path(), o.Pos, o.Props}, nil
}

// Validate validates the split operation.
func (o *SplitOperation) Validate() error {
	// Empty path is valid for split operation (root level)
	// Position bounds are checked in Apply method
	return nil
}

// splitSlateTextNode splits a Slate text node at the specified position
func splitSlateTextNode(nodeMap map[string]any, pos int, props map[string]any) []map[string]any {
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
	beforeNode := make(map[string]any)
	afterNode := make(map[string]any)

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

	return []map[string]any{beforeNode, afterNode}
}

// splitSlateElementNode splits a Slate element node at the specified position in its children
func splitSlateElementNode(nodeMap map[string]any, pos int, props map[string]any) []map[string]any {
	children, ok := nodeMap["children"].([]any)
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
	beforeNode := make(map[string]any)
	afterNode := make(map[string]any)

	// Copy properties from original node
	for k, v := range nodeMap {
		if k != "children" {
			beforeNode[k] = v
			afterNode[k] = v
		}
	}

	// Copy children slices to avoid mutation
	beforeChildrenCopy := make([]any, len(beforeChildren))
	copy(beforeChildrenCopy, beforeChildren)
	afterChildrenCopy := make([]any, len(afterChildren))
	copy(afterChildrenCopy, afterChildren)

	beforeNode["children"] = beforeChildrenCopy
	afterNode["children"] = afterChildrenCopy

	// Apply extra properties if specified
	for k, v := range props {
		beforeNode[k] = v
		afterNode[k] = v
	}

	return []map[string]any{beforeNode, afterNode}
}
