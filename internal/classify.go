package internal

// IsJSONPatchOperation reports whether op is a core JSON Patch
// (RFC 6902) operation.
func IsJSONPatchOperation(op string) bool {
	spec, ok := LookupOperation(OpType(op))
	return ok && spec.Families&FamilyJSONPatch != 0
}

// IsPredicateOperation reports whether op is any predicate operation.
func IsPredicateOperation(op string) bool {
	return IsFirstOrderPredicateOperation(op) || IsSecondOrderPredicateOperation(op)
}

// IsFirstOrderPredicateOperation reports whether op is a first-order
// predicate.
func IsFirstOrderPredicateOperation(op string) bool {
	spec, ok := LookupOperation(OpType(op))
	return ok && spec.Families&FamilyFirstOrderPredicate != 0
}

// IsSecondOrderPredicateOperation reports whether op is a
// second-order (composite) predicate.
func IsSecondOrderPredicateOperation(op string) bool {
	spec, ok := LookupOperation(OpType(op))
	return ok && spec.Families&FamilySecondOrderPredicate != 0
}

// IsJSONPatchExtendedOperation reports whether op is an extended
// operation.
func IsJSONPatchExtendedOperation(op string) bool {
	spec, ok := LookupOperation(OpType(op))
	return ok && spec.Families&FamilyExtended != 0
}
