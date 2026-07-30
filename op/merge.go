package op

import (
	"maps"
	"slices"

	"github.com/kaptinlin/jsonpatch/internal"
)

// MergeOperation represents an array merge operation.
// Path identifies the second element; Pos and Props carry reverse-split metadata.
type MergeOperation struct {
	BaseOp
	Pos   int            `json:"pos"`   // Merge position
	Props map[string]any `json:"props"` // Properties to apply after merge
}

// NewMerge creates a new merge operation.
func NewMerge(path []string, pos int, props map[string]any) *MergeOperation {
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

// Apply applies the merge operation to the document.
func (mg *MergeOperation) Apply(doc any) (internal.OpResult[any], error) {
	target, err := resolveTarget(doc, mg.Path(), targetPolicy{
		mustExist:     true,
		invalidParent: ErrInvalidTarget,
	})
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	if target.root || !target.array {
		return internal.OpResult[any]{}, ErrInvalidTarget
	}
	if target.index == 0 {
		return internal.OpResult[any]{}, ErrInvalidIndex
	}

	targetArray := target.container.([]any)
	one := targetArray[target.index-1]
	two := target.value
	merged := mg.mergeElements(one, two)

	newSlice := slices.Clone(targetArray)
	newSlice[target.index-1] = merged
	newSlice = slices.Delete(newSlice, target.index, target.index+1)
	newDoc, err := rebindContainer(doc, target.containerRef, newSlice)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	return internal.OpResult[any]{Doc: newDoc, Old: []any{one, two}}, nil
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
		return mg.mergeSlateNodes(one, two, mergeSlateTextNodes)
	}

	// Slate-like element node merging
	if isSlateElementNode(one) && isSlateElementNode(two) {
		return mg.mergeSlateNodes(one, two, mergeSlateElementNodes)
	}

	// Default: return array of both elements
	return []any{one, two}
}

func (mg *MergeOperation) mergeSlateNodes(one, two any, merge func(map[string]any, map[string]any) map[string]any) map[string]any {
	return merge(one.(map[string]any), two.(map[string]any))
}

// Validate validates the merge operation.
func (mg *MergeOperation) Validate() error {
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
	nodeMap, ok := value.(map[string]any)
	if !ok {
		return false
	}
	_, ok = nodeMap["children"].([]any)
	return ok
}

// copyPropsExcluding copies properties from two maps into a new map,
// skipping the specified key. Properties from the second map overwrite the first.
func copyPropsExcluding(one, two map[string]any, excludeKey string) map[string]any {
	result := maps.Clone(one)
	maps.Copy(result, two)
	delete(result, excludeKey)
	return result
}

// mergeSlateTextNodes merges two Slate text nodes by concatenating their text.
func mergeSlateTextNodes(one, two map[string]any) map[string]any {
	result := copyPropsExcluding(one, two, "text")
	textOne, _ := one["text"].(string)
	textTwo, _ := two["text"].(string)
	result["text"] = textOne + textTwo
	return result
}

// mergeSlateElementNodes merges two Slate element nodes by concatenating their children.
func mergeSlateElementNodes(one, two map[string]any) map[string]any {
	result := copyPropsExcluding(one, two, "children")
	childrenOne, _ := one["children"].([]any)
	childrenTwo, _ := two["children"].([]any)
	// Use slices.Concat (Go 1.22+) for efficient slice concatenation
	result["children"] = slices.Concat(childrenOne, childrenTwo)
	return result
}
