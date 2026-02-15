package internal

// IsJSONPatchOperation reports whether op is a core JSON Patch
// (RFC 6902) operation.
func IsJSONPatchOperation(op string) bool {
	switch OpType(op) { //nolint:exhaustive // intentionally matches only RFC 6902 operations
	case OpAddType, OpRemoveType, OpReplaceType,
		OpMoveType, OpCopyType, OpTestType:
		return true
	}
	return false
}

// IsPredicateOperation reports whether op is any predicate operation.
func IsPredicateOperation(op string) bool {
	return IsFirstOrderPredicateOperation(op) || IsSecondOrderPredicateOperation(op)
}

// IsFirstOrderPredicateOperation reports whether op is a first-order
// predicate.
func IsFirstOrderPredicateOperation(op string) bool {
	switch OpType(op) { //nolint:exhaustive // intentionally matches only first-order predicates
	case OpTestType, OpDefinedType, OpUndefinedType,
		OpTestTypeType, OpTestStringType, OpTestStringLenType,
		OpContainsType, OpEndsType, OpStartsType,
		OpInType, OpLessType, OpMoreType, OpMatchesType:
		return true
	}
	return false
}

// IsSecondOrderPredicateOperation reports whether op is a
// second-order (composite) predicate.
func IsSecondOrderPredicateOperation(op string) bool {
	switch OpType(op) { //nolint:exhaustive // intentionally matches only second-order predicates
	case OpAndType, OpOrType, OpNotType:
		return true
	}
	return false
}

// IsJSONPatchExtendedOperation reports whether op is an extended
// operation.
func IsJSONPatchExtendedOperation(op string) bool {
	switch OpType(op) { //nolint:exhaustive // intentionally matches only extended operations
	case OpStrInsType, OpStrDelType, OpFlipType,
		OpIncType, OpSplitType, OpMergeType, OpExtendType:
		return true
	}
	return false
}
