package internal

// OperationFamily describes how an operation participates in the patch
// vocabulary. Families are flags because "test" is both an RFC 6902 operation
// and a first-order predicate.
type OperationFamily uint8

const (
	// FamilyJSONPatch marks RFC 6902 operations.
	FamilyJSONPatch OperationFamily = 1 << iota
	// FamilyFirstOrderPredicate marks predicates that test one condition.
	FamilyFirstOrderPredicate
	// FamilySecondOrderPredicate marks predicates that combine predicates.
	FamilySecondOrderPredicate
	// FamilyRegexPredicate marks predicates that require regex support.
	FamilyRegexPredicate
	// FamilyExtended marks JSON Patch+ extended operations.
	FamilyExtended
)

// OperationCapability is the compile-time capability required by an operation.
type OperationCapability uint8

const (
	// CapabilityJSONPatch identifies the RFC 6902 compile capability.
	CapabilityJSONPatch OperationCapability = iota
	// CapabilityPredicate identifies the non-regex predicate compile capability.
	CapabilityPredicate
	// CapabilityRegexPredicate identifies the regex predicate compile capability.
	CapabilityRegexPredicate
	// CapabilityExtended identifies the extended operation compile capability.
	CapabilityExtended
)

// OperationSpec is the small executable spine shared by compile policy and
// compact/binary opcode resolution.
type OperationSpec struct {
	Type       OpType
	Families   OperationFamily
	Capability OperationCapability
	Code       int
}

var operationSpecs = []OperationSpec{
	{Type: OpAddType, Families: FamilyJSONPatch, Capability: CapabilityJSONPatch, Code: OpAddCode},
	{Type: OpRemoveType, Families: FamilyJSONPatch, Capability: CapabilityJSONPatch, Code: OpRemoveCode},
	{Type: OpReplaceType, Families: FamilyJSONPatch, Capability: CapabilityJSONPatch, Code: OpReplaceCode},
	{Type: OpMoveType, Families: FamilyJSONPatch, Capability: CapabilityJSONPatch, Code: OpMoveCode},
	{Type: OpCopyType, Families: FamilyJSONPatch, Capability: CapabilityJSONPatch, Code: OpCopyCode},
	{Type: OpTestType, Families: FamilyJSONPatch | FamilyFirstOrderPredicate, Capability: CapabilityJSONPatch, Code: OpTestCode},

	{Type: OpContainsType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpContainsCode},
	{Type: OpDefinedType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpDefinedCode},
	{Type: OpUndefinedType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpUndefinedCode},
	{Type: OpTypeType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpTypeCode},
	{Type: OpTestTypeType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpTestTypeCode},
	{Type: OpTestStringType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpTestStringCode},
	{Type: OpTestStringLenType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpTestStringLenCode},
	{Type: OpEndsType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpEndsCode},
	{Type: OpStartsType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpStartsCode},
	{Type: OpInType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpInCode},
	{Type: OpLessType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpLessCode},
	{Type: OpMoreType, Families: FamilyFirstOrderPredicate, Capability: CapabilityPredicate, Code: OpMoreCode},
	{Type: OpMatchesType, Families: FamilyFirstOrderPredicate | FamilyRegexPredicate, Capability: CapabilityRegexPredicate, Code: OpMatchesCode},

	{Type: OpAndType, Families: FamilySecondOrderPredicate, Capability: CapabilityPredicate, Code: OpAndCode},
	{Type: OpOrType, Families: FamilySecondOrderPredicate, Capability: CapabilityPredicate, Code: OpOrCode},
	{Type: OpNotType, Families: FamilySecondOrderPredicate, Capability: CapabilityPredicate, Code: OpNotCode},

	{Type: OpFlipType, Families: FamilyExtended, Capability: CapabilityExtended, Code: OpFlipCode},
	{Type: OpIncType, Families: FamilyExtended, Capability: CapabilityExtended, Code: OpIncCode},
	{Type: OpStrInsType, Families: FamilyExtended, Capability: CapabilityExtended, Code: OpStrInsCode},
	{Type: OpStrDelType, Families: FamilyExtended, Capability: CapabilityExtended, Code: OpStrDelCode},
	{Type: OpSplitType, Families: FamilyExtended, Capability: CapabilityExtended, Code: OpSplitCode},
	{Type: OpMergeType, Families: FamilyExtended, Capability: CapabilityExtended, Code: OpMergeCode},
	{Type: OpExtendType, Families: FamilyExtended, Capability: CapabilityExtended, Code: OpExtendCode},
}

var (
	operationByType = make(map[OpType]OperationSpec, len(operationSpecs))
	operationByCode = make(map[int]OperationSpec, len(operationSpecs))
)

func init() {
	for _, spec := range operationSpecs {
		operationByType[spec.Type] = spec
		operationByCode[spec.Code] = spec
	}
}

// OperationSpecs returns the operation vocabulary spine.
func OperationSpecs() []OperationSpec {
	specs := make([]OperationSpec, len(operationSpecs))
	copy(specs, operationSpecs)
	return specs
}

// LookupOperation returns the vocabulary entry for opType.
func LookupOperation(opType OpType) (OperationSpec, bool) {
	spec, ok := operationByType[opType]
	return spec, ok
}

// LookupOperationCode returns the vocabulary entry for code.
func LookupOperationCode(code int) (OperationSpec, bool) {
	spec, ok := operationByCode[code]
	return spec, ok
}
