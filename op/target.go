package op

import (
	"errors"
	"slices"
	"strconv"

	"github.com/kaptinlin/jsonpointer"
)

type targetPolicy struct {
	mustExist     bool
	allowArrayEnd bool
	invalidParent error
}

type resolvedTarget struct {
	root         bool
	container    any
	containerRef *jsonpointer.Reference
	key          string
	index        int
	array        bool
	exists       bool
	value        any
}

func resolveTarget(doc any, path []string, policy targetPolicy) (*resolvedTarget, error) {
	if len(path) == 0 {
		return &resolvedTarget{root: true, exists: true, value: doc}, nil
	}

	containerRef, err := jsonpointer.Find(doc, path[:len(path)-1]...)
	if err != nil {
		if errors.Is(err, jsonpointer.ErrIndexOutOfBounds) {
			return nil, ErrIndexOutOfRange
		}
		return nil, ErrPathNotFound
	}

	key := path[len(path)-1]
	switch container := containerRef.Val.(type) {
	case map[string]any:
		value, exists := container[key]
		if policy.mustExist && !exists {
			return nil, ErrPathNotFound
		}
		return &resolvedTarget{
			container:    container,
			containerRef: containerRef,
			key:          key,
			exists:       exists,
			value:        value,
		}, nil
	case []any:
		index, err := resolveArrayIndex(key, len(container), policy.allowArrayEnd)
		if err != nil {
			return nil, err
		}
		exists := index < len(container)
		if policy.mustExist && !exists {
			return nil, ErrIndexOutOfRange
		}
		var value any
		if exists {
			value = container[index]
		}
		return &resolvedTarget{
			container:    container,
			containerRef: containerRef,
			key:          key,
			index:        index,
			array:        true,
			exists:       exists,
			value:        value,
		}, nil
	default:
		if policy.invalidParent != nil {
			return nil, policy.invalidParent
		}
		return nil, ErrUnsupportedParentType
	}
}

func resolveArrayIndex(token string, length int, allowEnd bool) (int, error) {
	if token == "-" {
		if allowEnd {
			return length, nil
		}
		return 0, ErrPathNotFound
	}
	if !jsonpointer.IsValidIndex(token) {
		return 0, ErrPathNotFound
	}
	parsed, err := strconv.ParseInt(token, 10, strconv.IntSize)
	if err != nil {
		return 0, ErrPathNotFound
	}
	index := int(parsed)
	if index > length {
		return 0, ErrIndexOutOfRange
	}
	return index, nil
}

func addAtPath(doc any, path []string, value any) (any, any, error) {
	target, err := resolveTarget(doc, path, targetPolicy{
		allowArrayEnd: true,
		invalidParent: ErrCannotAddToValue,
	})
	if err != nil {
		return nil, nil, err
	}
	if target.root {
		return value, target.value, nil
	}
	if !target.array {
		container := target.container.(map[string]any)
		container[target.key] = value
		return doc, target.value, nil
	}

	container := target.container.([]any)
	updated := slices.Insert(slices.Clone(container), target.index, value)
	newDoc, err := rebindContainer(doc, target.containerRef, updated)
	if err != nil {
		return nil, nil, err
	}
	return newDoc, target.value, nil
}

func removeAtPath(doc any, path []string) (any, any, error) {
	target, err := resolveTarget(doc, path, targetPolicy{
		mustExist:     true,
		invalidParent: ErrCannotRemoveFromValue,
	})
	if err != nil {
		return nil, nil, err
	}
	if target.root {
		return nil, target.value, nil
	}
	if !target.array {
		delete(target.container.(map[string]any), target.key)
		return doc, target.value, nil
	}

	container := target.container.([]any)
	updated := slices.Delete(slices.Clone(container), target.index, target.index+1)
	newDoc, err := rebindContainer(doc, target.containerRef, updated)
	if err != nil {
		return nil, nil, err
	}
	return newDoc, target.value, nil
}

func replaceAtPath(doc any, path []string, value any) (any, any, error) {
	target, err := resolveTarget(doc, path, targetPolicy{
		mustExist:     true,
		invalidParent: ErrUnsupportedParentType,
	})
	if err != nil {
		return nil, nil, err
	}
	if target.root {
		return value, target.value, nil
	}
	if !target.array {
		target.container.(map[string]any)[target.key] = value
		return doc, target.value, nil
	}

	target.container.([]any)[target.index] = value
	return doc, target.value, nil
}

func rebindContainer(doc any, ref *jsonpointer.Reference, value any) (any, error) {
	if ref.Obj == nil {
		return value, nil
	}
	switch parent := ref.Obj.(type) {
	case map[string]any:
		parent[ref.Key] = value
		return doc, nil
	case []any:
		index, err := resolveArrayIndex(ref.Key, len(parent), false)
		if err != nil || index == len(parent) {
			return nil, ErrPathNotFound
		}
		parent[index] = value
		return doc, nil
	default:
		return nil, ErrCannotUpdateGrandparent
	}
}
