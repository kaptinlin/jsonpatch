package op

import "github.com/kaptinlin/jsonpatch/internal"

func codeFor(opType internal.OpType) int {
	spec, ok := internal.LookupOperation(opType)
	if !ok {
		return -1
	}
	return spec.Code
}

// Code returns the operation code.
func (a *AddOperation) Code() int {
	return codeFor(a.Op())
}

// ToJSON serializes the operation to JSON format.
func (a *AddOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpAddType),
		Path:  formatPath(a.path),
		Value: a.Value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (a *AddOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpAddType), a.path, a.Value}, nil
}

// Code returns the operation code.
func (ao *AndOperation) Code() int {
	return codeFor(ao.Op())
}

// ToJSON serializes the operation to JSON format.
func (ao *AndOperation) ToJSON() (internal.Operation, error) {
	operations, err := predicateOpsToJSON(ao.Operations, ErrInvalidPredicateInAnd)
	if err != nil {
		return internal.Operation{}, err
	}
	return internal.Operation{
		Op:    string(internal.OpAndType),
		Path:  formatPath(ao.path),
		Apply: operations,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (ao *AndOperation) ToCompact() (internal.CompactOperation, error) {
	opsCompact, err := predicateOpsToCompact(ao.Operations, ao.Path(), ErrInvalidPredicateInAnd)
	if err != nil {
		return nil, err
	}
	return internal.CompactOperation{codeFor(internal.OpAndType), ao.path, opsCompact}, nil
}

// Code returns the operation code.
func (co *ContainsOperation) Code() int {
	return codeFor(co.Op())
}

// ToJSON serializes the operation to JSON format.
func (co *ContainsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpContainsType),
		Path:       formatPath(co.Path()),
		Value:      co.Value,
		IgnoreCase: co.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (co *ContainsOperation) ToCompact() (internal.CompactOperation, error) {
	compact := internal.CompactOperation{codeFor(internal.OpContainsType), co.Path(), co.Value}
	if co.IgnoreCase {
		compact = append(compact, true)
	}
	return compact, nil
}

// Code returns the operation code.
func (c *CopyOperation) Code() int {
	return codeFor(c.Op())
}

// ToJSON serializes the operation to JSON format.
func (c *CopyOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpCopyType),
		Path: formatPath(c.path),
		From: formatPath(c.from),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (c *CopyOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpCopyType), c.path, c.from}, nil
}

// Code returns the operation code.
func (d *DefinedOperation) Code() int {
	return codeFor(d.Op())
}

// ToJSON serializes the operation to JSON format.
func (d *DefinedOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpDefinedType),
		Path: formatPath(d.Path()),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (d *DefinedOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpDefinedType), d.Path()}, nil
}

// Code returns the operation code.
func (e *EndsOperation) Code() int {
	return codeFor(e.Op())
}

// ToJSON serializes the operation to JSON format.
func (e *EndsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpEndsType),
		Path:       formatPath(e.Path()),
		Value:      e.Value,
		IgnoreCase: e.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (e *EndsOperation) ToCompact() (internal.CompactOperation, error) {
	compact := internal.CompactOperation{codeFor(internal.OpEndsType), e.Path(), e.Value}
	if e.IgnoreCase {
		compact = append(compact, true)
	}
	return compact, nil
}

// Code returns the operation code.
func (ex *ExtendOperation) Code() int {
	return codeFor(ex.Op())
}

// ToJSON serializes the operation to JSON format.
func (ex *ExtendOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpExtendType),
		Path:       formatPath(ex.Path()),
		Props:      ex.Properties,
		DeleteNull: ex.DeleteNull,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (ex *ExtendOperation) ToCompact() (internal.CompactOperation, error) {
	compact := internal.CompactOperation{codeFor(internal.OpExtendType), ex.Path(), ex.Properties}
	if ex.DeleteNull {
		compact = append(compact, true)
	}
	return compact, nil
}

// Code returns the operation code.
func (f *FlipOperation) Code() int {
	return codeFor(f.Op())
}

// ToJSON serializes the operation to JSON format.
func (f *FlipOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpFlipType),
		Path: formatPath(f.Path()),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (f *FlipOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpFlipType), f.Path()}, nil
}

// Code returns the operation code.
func (in *InOperation) Code() int {
	return codeFor(in.Op())
}

// ToJSON serializes the operation to JSON format.
func (in *InOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpInType),
		Path:  formatPath(in.Path()),
		Value: in.Value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (in *InOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpInType), in.Path(), in.Value}, nil
}

// Code returns the operation code.
func (ic *IncOperation) Code() int {
	return codeFor(ic.Op())
}

// ToJSON serializes the operation to JSON format.
func (ic *IncOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpIncType),
		Path: formatPath(ic.path),
		Inc:  ic.Inc,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (ic *IncOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpIncType), ic.path, ic.Inc}, nil
}

// Code returns the operation code.
func (l *LessOperation) Code() int {
	return codeFor(l.Op())
}

// ToJSON serializes the operation to JSON format.
func (l *LessOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpLessType),
		Path:  formatPath(l.Path()),
		Value: floatToJSONValue(l.Value),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (l *LessOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpLessType), l.Path(), l.Value}, nil
}

// Code returns the operation code.
func (ma *MatchesOperation) Code() int {
	return codeFor(ma.Op())
}

// ToJSON converts the operation to JSON representation.
func (ma *MatchesOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpMatchesType),
		Path:       formatPath(ma.Path()),
		Value:      ma.Pattern,
		IgnoreCase: ma.IgnoreCase,
	}, nil
}

// ToCompact converts the operation to compact array representation.
func (ma *MatchesOperation) ToCompact() (internal.CompactOperation, error) {
	compact := internal.CompactOperation{codeFor(internal.OpMatchesType), ma.Path(), ma.Pattern}
	if ma.IgnoreCase {
		compact = append(compact, true)
	}
	return compact, nil
}

// Code returns the operation code.
func (mg *MergeOperation) Code() int {
	return codeFor(mg.Op())
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
	compact := internal.CompactOperation{codeFor(internal.OpMergeType), mg.Path(), mg.Pos}
	if mg.Props != nil {
		compact = append(compact, mg.Props)
	}
	return compact, nil
}

// Code returns the operation code.
func (mo *MoreOperation) Code() int {
	return codeFor(mo.Op())
}

// ToJSON converts the operation to JSON representation.
func (mo *MoreOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpMoreType),
		Path:  formatPath(mo.Path()),
		Value: floatToJSONValue(mo.Value),
	}, nil
}

// ToCompact converts the operation to compact array representation.
func (mo *MoreOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpMoreType), mo.Path(), mo.Value}, nil
}

// Code returns the operation code.
func (m *MoveOperation) Code() int {
	return codeFor(m.Op())
}

// ToJSON serializes the operation to JSON format.
func (m *MoveOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpMoveType),
		Path: formatPath(m.path),
		From: formatPath(m.from),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (m *MoveOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpMoveType), m.path, m.from}, nil
}

// Code returns the operation code.
func (n *NotOperation) Code() int {
	return codeFor(n.Op())
}

// ToJSON serializes the operation to JSON format.
func (n *NotOperation) ToJSON() (internal.Operation, error) {
	if _, err := n.operand(); err != nil {
		return internal.Operation{}, err
	}
	opsJSON, err := predicateOpsToJSON(n.Operations, ErrInvalidPredicateInNot)
	if err != nil {
		return internal.Operation{}, err
	}

	return internal.Operation{
		Op:    string(internal.OpNotType),
		Path:  formatPath(n.Path()),
		Apply: opsJSON,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (n *NotOperation) ToCompact() (internal.CompactOperation, error) {
	if _, err := n.operand(); err != nil {
		return nil, err
	}
	opsCompact, err := predicateOpsToCompact(n.Operations, n.Path(), ErrInvalidPredicateInNot)
	if err != nil {
		return nil, err
	}

	return internal.CompactOperation{codeFor(internal.OpNotType), n.Path(), opsCompact}, nil
}

// Code returns the operation code.
func (oo *OrOperation) Code() int {
	return codeFor(oo.Op())
}

// ToJSON serializes the operation to JSON format.
func (oo *OrOperation) ToJSON() (internal.Operation, error) {
	opsJSON, err := predicateOpsToJSON(oo.Operations, ErrInvalidPredicateInOr)
	if err != nil {
		return internal.Operation{}, err
	}
	return internal.Operation{
		Op:    string(internal.OpOrType),
		Path:  formatPath(oo.Path()),
		Apply: opsJSON,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (oo *OrOperation) ToCompact() (internal.CompactOperation, error) {
	opsCompact, err := predicateOpsToCompact(oo.Operations, oo.Path(), ErrInvalidPredicateInOr)
	if err != nil {
		return nil, err
	}
	return internal.CompactOperation{codeFor(internal.OpOrType), oo.Path(), opsCompact}, nil
}

// Code returns the operation code.
func (r *RemoveOperation) Code() int {
	return codeFor(r.Op())
}

// ToJSON serializes the operation to JSON format.
func (r *RemoveOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpRemoveType),
		Path: formatPath(r.path),
	}

	if r.HasOldValue {
		result.OldValue = r.OldValue
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (r *RemoveOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpRemoveType), r.path}, nil
}

// Code returns the operation code.
func (rp *ReplaceOperation) Code() int {
	return codeFor(rp.Op())
}

// ToJSON serializes the operation to JSON format.
func (rp *ReplaceOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:    string(internal.OpReplaceType),
		Path:  formatPath(rp.path),
		Value: rp.Value,
	}

	if rp.OldValue != nil {
		result.OldValue = rp.OldValue
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (rp *ReplaceOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpReplaceType), rp.path, rp.Value}, nil
}

// Code returns the operation code.
func (sp *SplitOperation) Code() int {
	return codeFor(sp.Op())
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
	compact := internal.CompactOperation{codeFor(internal.OpSplitType), sp.Path(), sp.Pos}
	if sp.Props != nil {
		compact = append(compact, sp.Props)
	}
	return compact, nil
}

// Code returns the operation code.
func (s *StartsOperation) Code() int {
	return codeFor(s.Op())
}

// ToJSON serializes the operation to JSON format.
func (s *StartsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpStartsType),
		Path:       formatPath(s.Path()),
		Value:      s.Value,
		IgnoreCase: s.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (s *StartsOperation) ToCompact() (internal.CompactOperation, error) {
	compact := internal.CompactOperation{codeFor(internal.OpStartsType), s.Path(), s.Value}
	if s.IgnoreCase {
		compact = append(compact, true)
	}
	return compact, nil
}

// Code returns the operation code.
func (sd *StrDelOperation) Code() int {
	return codeFor(sd.Op())
}

// ToJSON serializes the operation to JSON format.
func (sd *StrDelOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpStrDelType),
		Path: formatPath(sd.Path()),
		Pos:  sd.Pos,
		Len:  sd.Len,
	}

	if sd.HasStr {
		result.Str = sd.Str
		result.Len = 0
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (sd *StrDelOperation) ToCompact() (internal.CompactOperation, error) {
	if sd.HasStr {
		return internal.CompactOperation{codeFor(internal.OpStrDelType), sd.Path(), sd.Pos, sd.Str}, nil
	}
	return internal.CompactOperation{codeFor(internal.OpStrDelType), sd.Path(), sd.Pos, 0, sd.Len}, nil
}

// Code returns the operation code.
func (si *StrInsOperation) Code() int {
	return codeFor(si.Op())
}

// ToJSON serializes the operation to JSON format.
func (si *StrInsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpStrInsType),
		Path: formatPath(si.Path()),
		Pos:  si.Pos,
		Str:  si.Str,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (si *StrInsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpStrInsType), si.Path(), si.Pos, si.Str}, nil
}

// Code returns the operation code.
func (t *TestOperation) Code() int {
	return codeFor(t.Op())
}

// ToJSON serializes the operation to JSON format.
func (t *TestOperation) ToJSON() (internal.Operation, error) {
	op := internal.Operation{
		Op:    string(internal.OpTestType),
		Path:  formatPath(t.path),
		Value: t.Value,
	}
	if t.NotFlag {
		op.Not = true
	}
	return op, nil
}

// ToCompact serializes the operation to compact format.
func (t *TestOperation) ToCompact() (internal.CompactOperation, error) {
	if t.NotFlag {
		return internal.CompactOperation{codeFor(internal.OpTestType), t.path, t.Value, 1}, nil
	}
	return internal.CompactOperation{codeFor(internal.OpTestType), t.path, t.Value}, nil
}

// Code returns the operation code.
func (ts *TestStringOperation) Code() int {
	return codeFor(ts.Op())
}

// ToJSON serializes the operation to JSON format.
func (ts *TestStringOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpTestStringType),
		Path:       formatPath(ts.Path()),
		Str:        ts.Str,
		Pos:        ts.Pos,
		Not:        ts.NotFlag,
		IgnoreCase: ts.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (ts *TestStringOperation) ToCompact() (internal.CompactOperation, error) {
	compact := internal.CompactOperation{codeFor(internal.OpTestStringType), ts.Path(), ts.Pos, ts.Str}
	if ts.NotFlag {
		compact = append(compact, true)
	}
	return compact, nil
}

// Code returns the operation code.
func (tl *TestStringLenOperation) Code() int {
	return codeFor(tl.Op())
}

// ToJSON serializes the operation to JSON format.
func (tl *TestStringLenOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpTestStringLenType),
		Path: formatPath(tl.Path()),
		Len:  int(tl.Length),
		Not:  tl.NotFlag,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (tl *TestStringLenOperation) ToCompact() (internal.CompactOperation, error) {
	compact := internal.CompactOperation{codeFor(internal.OpTestStringLenType), tl.Path(), tl.Length}
	if tl.NotFlag {
		compact = append(compact, true)
	}
	return compact, nil
}

// Code returns the operation code.
func (tt *TestTypeOperation) Code() int {
	return codeFor(tt.Op())
}

// ToJSON serializes the operation to JSON format.
func (tt *TestTypeOperation) ToJSON() (internal.Operation, error) {
	if len(tt.Types) == 1 {
		return internal.Operation{
			Op:   string(internal.OpTestTypeType),
			Path: formatPath(tt.Path()),
			Type: tt.Types[0],
		}, nil
	}
	return internal.Operation{
		Op:   string(internal.OpTestTypeType),
		Path: formatPath(tt.Path()),
		Type: tt.Types,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (tt *TestTypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpTestTypeType), tt.Path(), tt.Types}, nil
}

// Code returns the operation code.
func (tp *TypeOperation) Code() int {
	return codeFor(tp.Op())
}

// ToJSON serializes the operation to JSON format.
func (tp *TypeOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpTypeType),
		Path:  formatPath(tp.Path()),
		Value: tp.TypeValue,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (tp *TypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpTypeType), tp.Path(), tp.TypeValue}, nil
}

// Code returns the operation code.
func (u *UndefinedOperation) Code() int {
	return codeFor(u.Op())
}

// ToJSON serializes the operation to JSON format.
func (u *UndefinedOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpUndefinedType),
		Path: formatPath(u.path),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (u *UndefinedOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{codeFor(internal.OpUndefinedType), u.path}, nil
}
