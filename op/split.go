package op

import (
	"maps"
	"slices"

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
func (sp *SplitOperation) Op() internal.OpType {
	return internal.OpSplitType
}

// Code returns the operation code.
func (sp *SplitOperation) Code() int {
	return internal.OpSplitCode
}

// Apply applies the split operation following TypeScript reference.
func (sp *SplitOperation) Apply(doc any) (internal.OpResult[any], error) {
	var target any
	var err error

	if len(sp.Path()) == 0 {
		target = doc
	} else {
		target, err = value(doc, sp.Path())
		if err != nil {
			return internal.OpResult[any]{}, err
		}
	}

	parts := sp.splitValue(target)

	if len(sp.Path()) == 0 {
		return internal.OpResult[any]{Doc: parts, Old: target}, nil
	}

	parent, key, err := navigateToParent(doc, sp.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if slice, ok := parent.([]any); ok {
		if index, ok := key.(int); ok {
			splitResult := parts.([]any)
			newSlice := slices.Clone(slice)
			newSlice[index] = splitResult[0]
			newSlice = slices.Insert(newSlice, index+1, splitResult[1])

			parentPath := sp.Path()[:len(sp.Path())-1]
			if len(parentPath) == 0 {
				return internal.OpResult[any]{Doc: newSlice, Old: target}, nil
			}
			err = setValueAtPath(doc, parentPath, newSlice)
			if err != nil {
				return internal.OpResult[any]{}, err
			}
		}
	} else {
		err = setValueAtPath(doc, sp.Path(), parts)
		if err != nil {
			return internal.OpResult[any]{}, err
		}
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

func (sp *SplitOperation) splitValue(value any) any {
	switch v := value.(type) {
	case string:
		return sp.splitString(v)
	case float64:
		return []any{sp.Pos, v - sp.Pos}
	case int:
		return []any{sp.Pos, float64(v) - sp.Pos}
	case bool:
		return []any{v, v}
	case map[string]any:
		if isSlateTextNode(v) {
			propsMap, _ := sp.Props.(map[string]any)
			results := splitSlateTextNode(v, int(sp.Pos), propsMap)
			if results != nil {
				return []any{results[0], results[1]}
			}
		}
		if isSlateElementNode(v) {
			propsMap, _ := sp.Props.(map[string]any)
			results := splitSlateElementNode(v, int(sp.Pos), propsMap)
			if results != nil {
				return []any{results[0], results[1]}
			}
		}
		return []any{v, v}
	default:
		return []any{value, value}
	}
}

func (sp *SplitOperation) splitString(s string) []any {
	runes := []rune(s)
	pos := int(sp.Pos)

	if pos < 0 {
		pos = max(len(runes)+pos, 0)
	}

	if pos > len(runes) {
		pos = len(runes)
	}

	before := string(runes[:pos])
	after := string(runes[pos:])

	if propsMap, ok := sp.Props.(map[string]any); ok {
		beforeNode := map[string]any{"text": before}
		afterNode := map[string]any{"text": after}

		for k, v := range propsMap {
			beforeNode[k] = v
			afterNode[k] = v
		}

		return []any{beforeNode, afterNode}
	}

	return []any{before, after}
}

// ToJSON serializes the operation to JSON format.
func (sp *SplitOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpSplitType),
		Path: formatPath(sp.Path()),
		Pos:  int(sp.Pos),
	}
	if sp.Props != nil {
		if props, ok := sp.Props.(map[string]any); ok {
			result.Props = props
		}
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (sp *SplitOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpSplitCode, sp.Path(), sp.Pos, sp.Props}, nil
}

// Validate validates the split operation.
func (sp *SplitOperation) Validate() error {
	return nil
}

// splitNodePair creates two new nodes from an original, copying all properties
// except excludeKey, and applying optional extra properties to both.
func splitNodePair(nodeMap map[string]any, excludeKey string, props map[string]any) (map[string]any, map[string]any) {
	beforeNode := maps.Clone(nodeMap)
	delete(beforeNode, excludeKey)
	maps.Copy(beforeNode, props)

	afterNode := maps.Clone(beforeNode)
	return beforeNode, afterNode
}

// splitSlateTextNode splits a Slate text node at the specified position.
func splitSlateTextNode(nodeMap map[string]any, pos int, props map[string]any) []map[string]any {
	text, ok := nodeMap["text"].(string)
	if !ok {
		return nil
	}

	runes := []rune(text)
	pos = max(0, min(pos, len(runes)))

	beforeNode, afterNode := splitNodePair(nodeMap, "text", props)
	beforeNode["text"] = string(runes[:pos])
	afterNode["text"] = string(runes[pos:])

	return []map[string]any{beforeNode, afterNode}
}

// splitSlateElementNode splits a Slate element node at the specified position in its children.
func splitSlateElementNode(nodeMap map[string]any, pos int, props map[string]any) []map[string]any {
	children, ok := nodeMap["children"].([]any)
	if !ok {
		return nil
	}

	pos = max(0, min(pos, len(children)))

	beforeNode, afterNode := splitNodePair(nodeMap, "children", props)

	// Clone children slices to avoid mutation (Go 1.21+)
	beforeNode["children"] = slices.Clone(children[:pos])
	afterNode["children"] = slices.Clone(children[pos:])

	return []map[string]any{beforeNode, afterNode}
}
