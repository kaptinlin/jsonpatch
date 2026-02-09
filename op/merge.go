package op

import (
	"maps"

	"github.com/kaptinlin/jsonpatch/internal"
)

// MergeOperation represents an array merge operation.
// path: target path
// pos: merge position (array index)
// props: properties to apply after merge (can be nil)
// Only supports array type fields.
type MergeOperation struct {
	BaseOp
	Pos   float64        `json:"pos"`   // Merge position
	Props map[string]any `json:"props"` // Properties to apply after merge
}

// NewMerge creates a new merge operation.
func NewMerge(path []string, pos float64, props map[string]any) *MergeOperation {
	return &MergeOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Props:  props,
	}
}

// Op returns the operation type.
func (mg *MergeOperation) Op() internal.OpType {
	return internal.OpMergeType
}

// Code returns the operation code.
func (mg *MergeOperation) Code() int {
	return internal.OpMergeCode
}

// Apply applies the merge operation following TypeScript reference.
func (mg *MergeOperation) Apply(doc any) (internal.OpResult[any], error) {
	// TypeScript reference: merge works on arrays directly using pos parameter
	var targetArray []any

	if len(mg.Path()) == 0 {
		// Root level array
		slice, ok := doc.([]any)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidTarget
		}
		targetArray = slice
	} else {
		// Get array at path
		target, err := getValue(doc, mg.Path())
		if err != nil {
			return internal.OpResult[any]{}, err
		}
		slice, ok := target.([]any)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidTarget
		}
		targetArray = slice
	}

	pos := int(mg.Pos)
	// TypeScript: if (ref.key <= 0) throw new Error('INVALID_KEY');
	if pos <= 0 || pos >= len(targetArray) {
		return internal.OpResult[any]{}, ErrInvalidIndex
	}

	// Get elements to merge (pos-1 and pos)
	one := targetArray[pos-1]
	two := targetArray[pos]
	merged := mg.mergeElements(one, two)

	// Create new array with merged result
	newSlice := make([]any, len(targetArray)-1)
	copy(newSlice[:pos-1], targetArray[:pos-1])
	newSlice[pos-1] = merged
	copy(newSlice[pos:], targetArray[pos+1:])

	// Update the document
	if len(mg.Path()) == 0 {
		// Root array
		return internal.OpResult[any]{Doc: newSlice, Old: []any{one, two}}, nil
	}

	err := setValueAtPath(doc, mg.Path(), newSlice)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: []any{one, two}}, nil
}

// mergeElements merges two elements based on their type.
func (mg *MergeOperation) mergeElements(one, two any) any {
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
	if isSlateTextNode(one) && isSlateTextNode(two) {
		merged := mergeSlateTextNodes(one.(map[string]any), two.(map[string]any))
		// Apply props if specified
		maps.Copy(merged, mg.Props)
		return merged
	}

	// Slate-like element node merging
	if isSlateElementNode(one) && isSlateElementNode(two) {
		merged := mergeSlateElementNodes(one.(map[string]any), two.(map[string]any))
		// Apply props if specified
		maps.Copy(merged, mg.Props)
		return merged
	}

	// Default: return array of both elements
	return []any{one, two}
}

// ToJSON serializes the operation to JSON format.
func (mg *MergeOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpMergeType),
		Path: formatPath(mg.Path()),
		Pos:  int(mg.Pos),
	}
	if len(mg.Props) > 0 {
		result.Props = mg.Props
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (mg *MergeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMergeCode, mg.Path(), mg.Props}, nil
}

// Validate validates the merge operation.
func (mg *MergeOperation) Validate() error {
	if mg.Pos < 0 {
		return ErrPositionNegative
	}
	return nil
}

// Slate node helper functions (inlined from pkg/slate)

// isSlateTextNode checks if a value is a Slate.js text node
func isSlateTextNode(value any) bool {
	if nodeMap, ok := value.(map[string]any); ok {
		_, hasText := nodeMap["text"]
		return hasText
	}
	return false
}

// isSlateElementNode checks if a value is a Slate.js element node
func isSlateElementNode(value any) bool {
	if nodeMap, ok := value.(map[string]any); ok {
		if children, hasChildren := nodeMap["children"]; hasChildren {
			_, isArray := children.([]any)
			return isArray
		}
	}
	return false
}

// mergeSlateTextNodes merges two Slate text nodes by concatenating their text
func mergeSlateTextNodes(one, two map[string]any) map[string]any {
	result := make(map[string]any)

	// Copy properties from first node
	for k, v := range one {
		if k != "text" {
			result[k] = v
		}
	}
	// Copy properties from second node (overwrites first)
	for k, v := range two {
		if k != "text" {
			result[k] = v
		}
	}

	// Concatenate text
	textOne, _ := one["text"].(string)
	textTwo, _ := two["text"].(string)
	result["text"] = textOne + textTwo

	return result
}

// mergeSlateElementNodes merges two Slate element nodes by concatenating their children
func mergeSlateElementNodes(one, two map[string]any) map[string]any {
	result := make(map[string]any)

	// Copy properties from first node
	for k, v := range one {
		if k != "children" {
			result[k] = v
		}
	}
	// Copy properties from second node (overwrites first)
	for k, v := range two {
		if k != "children" {
			result[k] = v
		}
	}

	// Concatenate children
	childrenOne, _ := one["children"].([]any)
	childrenTwo, _ := two["children"].([]any)
	mergedChildren := make([]any, 0, len(childrenOne)+len(childrenTwo))
	mergedChildren = append(mergedChildren, childrenOne...)
	mergedChildren = append(mergedChildren, childrenTwo...)
	result["children"] = mergedChildren

	return result
}
