package slate

// SlateNode represents a Slate.js node that can be either a text node or an element node.
// It uses a union-type approach similar to TypeScript's SlateNode.
type SlateNode struct {
	// Text field is present only for text nodes
	Text *string `json:"text,omitempty"`
	// Children field is present only for element nodes
	Children []SlateNode `json:"children,omitempty"`
	// Properties contains additional properties (equivalent to [key: string]: unknown)
	Properties map[string]interface{} `json:",inline"`
}

// SlateTextNode is a type alias for better API clarity
type SlateTextNode = SlateNode

// SlateElementNode is a type alias for better API clarity
type SlateElementNode = SlateNode

// IsText returns true if this is a text node
func (n *SlateNode) IsText() bool {
	return n.Text != nil
}

// IsElement returns true if this is an element node
func (n *SlateNode) IsElement() bool {
	return n.Children != nil
}

// GetText returns the text content (empty string if not a text node)
func (n *SlateNode) GetText() string {
	if n.Text != nil {
		return *n.Text
	}
	return ""
}

// SetText sets the text content and converts to text node
func (n *SlateNode) SetText(text string) {
	n.Text = &text
	n.Children = nil // Clear children when setting text
}

// GetChildren returns the children (empty slice if not an element node)
func (n *SlateNode) GetChildren() []SlateNode {
	if n.Children != nil {
		return n.Children
	}
	return []SlateNode{}
}

// SetChildren sets the children and converts to element node
func (n *SlateNode) SetChildren(children []SlateNode) {
	n.Children = children
	n.Text = nil // Clear text when setting children
}

// GetProperty returns a property value
func (n *SlateNode) GetProperty(key string) interface{} {
	if n.Properties == nil {
		return nil
	}
	return n.Properties[key]
}

// SetProperty sets a property value
func (n *SlateNode) SetProperty(key string, value interface{}) {
	if n.Properties == nil {
		n.Properties = make(map[string]interface{})
	}
	n.Properties[key] = value
}

// ToMap converts the node to map[string]interface{} for backward compatibility
func (n *SlateNode) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	// Copy properties first
	for k, v := range n.Properties {
		result[k] = v
	}

	// Add type-specific fields
	if n.IsText() {
		result["text"] = n.GetText()
	} else if n.IsElement() {
		children := make([]interface{}, len(n.Children))
		for i, child := range n.Children {
			children[i] = child.ToMap()
		}
		result["children"] = children
	}

	return result
}

// FromMap creates a SlateNode from map[string]interface{} for backward compatibility
func FromMap(data map[string]interface{}) *SlateNode {
	node := &SlateNode{
		Properties: make(map[string]interface{}),
	}

	for k, v := range data {
		switch k {
		case "text":
			if textStr, ok := v.(string); ok {
				node.SetText(textStr)
			}
		case "children":
			if childrenData, ok := v.([]interface{}); ok {
				children := make([]SlateNode, len(childrenData))
				for i, childData := range childrenData {
					if childMap, ok := childData.(map[string]interface{}); ok {
						children[i] = *FromMap(childMap)
					}
				}
				node.SetChildren(children)
			}
		default:
			node.SetProperty(k, v)
		}
	}

	return node
}

// IsTextNode checks if a value is a Slate.js text node
func IsTextNode(value interface{}) bool {
	if node, ok := value.(*SlateNode); ok {
		return node.IsText()
	}
	if nodeMap, ok := value.(map[string]interface{}); ok {
		_, hasText := nodeMap["text"]
		return hasText
	}
	return false
}

// IsElementNode checks if a value is a Slate.js element node
func IsElementNode(value interface{}) bool {
	if node, ok := value.(*SlateNode); ok {
		return node.IsElement()
	}
	if nodeMap, ok := value.(map[string]interface{}); ok {
		if children, hasChildren := nodeMap["children"]; hasChildren {
			_, isArray := children.([]interface{})
			return isArray
		}
	}
	return false
}

// NewTextNode creates a new Slate text node
func NewTextNode(text string, properties map[string]interface{}) *SlateNode {
	node := &SlateNode{
		Properties: make(map[string]interface{}),
	}
	node.SetText(text)

	for k, v := range properties {
		if k != "text" { // Don't overwrite the text field
			node.SetProperty(k, v)
		}
	}

	return node
}

// NewElementNode creates a new Slate element node
func NewElementNode(children []SlateNode, properties map[string]interface{}) *SlateNode {
	node := &SlateNode{
		Properties: make(map[string]interface{}),
	}

	if children == nil {
		children = []SlateNode{}
	}
	node.SetChildren(children)

	for k, v := range properties {
		if k != "children" { // Don't overwrite the children field
			node.SetProperty(k, v)
		}
	}

	return node
}

// MergeTextNodes merges two Slate text nodes by concatenating their text and merging properties
func MergeTextNodes(one, two *SlateNode) *SlateNode {
	if !one.IsText() || !two.IsText() {
		return nil // Can't merge non-text nodes
	}

	result := NewTextNode(one.GetText()+two.GetText(), nil)

	// Copy properties from both nodes (two overwrites one)
	for k, v := range one.Properties {
		result.SetProperty(k, v)
	}
	for k, v := range two.Properties {
		result.SetProperty(k, v)
	}

	return result
}

// MergeElementNodes merges two Slate element nodes by concatenating their children and merging properties
func MergeElementNodes(one, two *SlateNode) *SlateNode {
	if !one.IsElement() || !two.IsElement() {
		return nil // Can't merge non-element nodes
	}

	mergedChildren := make([]SlateNode, 0, len(one.Children)+len(two.Children))
	mergedChildren = append(mergedChildren, one.Children...)
	mergedChildren = append(mergedChildren, two.Children...)

	result := NewElementNode(mergedChildren, nil)

	// Copy properties from both nodes (two overwrites one)
	for k, v := range one.Properties {
		result.SetProperty(k, v)
	}
	for k, v := range two.Properties {
		result.SetProperty(k, v)
	}

	return result
}

// SplitTextNode splits a Slate text node at the specified position
func SplitTextNode(node *SlateNode, pos int, props map[string]interface{}) []*SlateNode {
	if !node.IsText() {
		return nil // Can't split non-text node
	}

	text := node.GetText()
	runes := []rune(text)

	// Handle edge cases
	if pos > len(runes) {
		pos = len(runes)
	}
	if pos < 0 {
		pos = 0
	}

	before := string(runes[:pos])
	after := string(runes[pos:])

	// Create two new nodes
	beforeNode := NewTextNode(before, node.Properties)
	afterNode := NewTextNode(after, node.Properties)

	// Apply extra properties if specified
	for k, v := range props {
		beforeNode.SetProperty(k, v)
		afterNode.SetProperty(k, v)
	}

	return []*SlateNode{beforeNode, afterNode}
}

// SplitElementNode splits a Slate element node at the specified position in its children
func SplitElementNode(node *SlateNode, pos int, props map[string]interface{}) []*SlateNode {
	if !node.IsElement() {
		return nil // Can't split non-element node
	}

	children := node.GetChildren()

	// Handle edge cases
	if pos > len(children) {
		pos = len(children)
	}
	if pos < 0 {
		pos = 0
	}

	before := children[:pos]
	after := children[pos:]

	// Create two new nodes
	beforeNode := NewElementNode(before, node.Properties)
	afterNode := NewElementNode(after, node.Properties)

	// Apply extra properties if specified
	for k, v := range props {
		beforeNode.SetProperty(k, v)
		afterNode.SetProperty(k, v)
	}

	return []*SlateNode{beforeNode, afterNode}
}

// Legacy functions for backward compatibility with map[string]interface{}

// MergeTextNodesFromMaps merges two Slate text nodes from maps (backward compatibility)
func MergeTextNodesFromMaps(one, two map[string]interface{}) map[string]interface{} {
	nodeOne := FromMap(one)
	nodeTwo := FromMap(two)
	result := MergeTextNodes(nodeOne, nodeTwo)
	if result == nil {
		return nil
	}
	return result.ToMap()
}

// MergeElementNodesFromMaps merges two Slate element nodes from maps (backward compatibility)
func MergeElementNodesFromMaps(one, two map[string]interface{}) map[string]interface{} {
	nodeOne := FromMap(one)
	nodeTwo := FromMap(two)
	result := MergeElementNodes(nodeOne, nodeTwo)
	if result == nil {
		return nil
	}
	return result.ToMap()
}

// SplitTextNodeFromMap splits a Slate text node from map (backward compatibility)
func SplitTextNodeFromMap(nodeMap map[string]interface{}, pos int, props map[string]interface{}) []map[string]interface{} {
	node := FromMap(nodeMap)
	results := SplitTextNode(node, pos, props)
	if results == nil {
		return nil
	}

	maps := make([]map[string]interface{}, len(results))
	for i, result := range results {
		maps[i] = result.ToMap()
	}
	return maps
}

// SplitElementNodeFromMap splits a Slate element node from map (backward compatibility)
func SplitElementNodeFromMap(nodeMap map[string]interface{}, pos int, props map[string]interface{}) []map[string]interface{} {
	node := FromMap(nodeMap)
	results := SplitElementNode(node, pos, props)
	if results == nil {
		return nil
	}

	maps := make([]map[string]interface{}, len(results))
	for i, result := range results {
		maps[i] = result.ToMap()
	}
	return maps
}
