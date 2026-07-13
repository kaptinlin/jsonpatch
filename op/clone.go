package op

import (
	"fmt"

	"github.com/kaptinlin/deepclone"

	"github.com/kaptinlin/jsonpatch/internal"
)

func cloneBaseOp(base BaseOp) BaseOp {
	return NewBaseOpWithFrom(base.path, base.from)
}

func cloneValue(value any) any {
	return deepclone.Clone(value)
}

func clonePredicateOps(operations []any, errInvalid error) ([]any, error) {
	result := make([]any, len(operations))
	for i, operation := range operations {
		predicateOp, ok := operation.(internal.PredicateOp)
		if !ok {
			return nil, errInvalid
		}
		cloned, err := predicateOp.Clone()
		if err != nil {
			return nil, err
		}
		clonedPredicate, ok := cloned.(internal.PredicateOp)
		if !ok {
			return nil, fmt.Errorf("%w: cloned predicate %T is not a predicate", errInvalid, cloned)
		}
		result[i] = clonedPredicate
	}
	return result, nil
}

// Clone returns an independent copy of the operation.
func (a *AddOperation) Clone() (internal.Op, error) {
	return &AddOperation{BaseOp: cloneBaseOp(a.BaseOp), Value: cloneValue(a.Value)}, nil
}

// Clone returns an independent copy of the operation.
func (r *RemoveOperation) Clone() (internal.Op, error) {
	return &RemoveOperation{
		BaseOp:      cloneBaseOp(r.BaseOp),
		OldValue:    cloneValue(r.OldValue),
		HasOldValue: r.HasOldValue,
	}, nil
}

// Clone returns an independent copy of the operation.
func (rp *ReplaceOperation) Clone() (internal.Op, error) {
	return &ReplaceOperation{
		BaseOp:      cloneBaseOp(rp.BaseOp),
		Value:       cloneValue(rp.Value),
		OldValue:    cloneValue(rp.OldValue),
		HasOldValue: rp.HasOldValue,
	}, nil
}

// Clone returns an independent copy of the operation.
func (m *MoveOperation) Clone() (internal.Op, error) {
	return &MoveOperation{BaseOp: cloneBaseOp(m.BaseOp)}, nil
}

// Clone returns an independent copy of the operation.
func (c *CopyOperation) Clone() (internal.Op, error) {
	return &CopyOperation{BaseOp: cloneBaseOp(c.BaseOp)}, nil
}

// Clone returns an independent copy of the operation.
func (t *TestOperation) Clone() (internal.Op, error) {
	return &TestOperation{BaseOp: cloneBaseOp(t.BaseOp), Value: cloneValue(t.Value), NotFlag: t.NotFlag}, nil
}

// Clone returns an independent copy of the operation.
func (co *ContainsOperation) Clone() (internal.Op, error) {
	return &ContainsOperation{BaseOp: cloneBaseOp(co.BaseOp), Value: co.Value, IgnoreCase: co.IgnoreCase}, nil
}

// Clone returns an independent copy of the operation.
func (d *DefinedOperation) Clone() (internal.Op, error) {
	return &DefinedOperation{BaseOp: cloneBaseOp(d.BaseOp)}, nil
}

// Clone returns an independent copy of the operation.
func (u *UndefinedOperation) Clone() (internal.Op, error) {
	return &UndefinedOperation{BaseOp: cloneBaseOp(u.BaseOp)}, nil
}

// Clone returns an independent copy of the operation.
func (tp *TypeOperation) Clone() (internal.Op, error) {
	return &TypeOperation{BaseOp: cloneBaseOp(tp.BaseOp), TypeValue: tp.TypeValue}, nil
}

// Clone returns an independent copy of the operation.
func (tt *TestTypeOperation) Clone() (internal.Op, error) {
	return &TestTypeOperation{BaseOp: cloneBaseOp(tt.BaseOp), Types: deepclone.Clone(tt.Types)}, nil
}

// Clone returns an independent copy of the operation.
func (ts *TestStringOperation) Clone() (internal.Op, error) {
	return &TestStringOperation{
		BaseOp:     cloneBaseOp(ts.BaseOp),
		Str:        ts.Str,
		Pos:        ts.Pos,
		NotFlag:    ts.NotFlag,
		IgnoreCase: ts.IgnoreCase,
	}, nil
}

// Clone returns an independent copy of the operation.
func (tl *TestStringLenOperation) Clone() (internal.Op, error) {
	return &TestStringLenOperation{BaseOp: cloneBaseOp(tl.BaseOp), Length: tl.Length, NotFlag: tl.NotFlag}, nil
}

// Clone returns an independent copy of the operation.
func (e *EndsOperation) Clone() (internal.Op, error) {
	return &EndsOperation{BaseOp: cloneBaseOp(e.BaseOp), Value: e.Value, IgnoreCase: e.IgnoreCase}, nil
}

// Clone returns an independent copy of the operation.
func (s *StartsOperation) Clone() (internal.Op, error) {
	return &StartsOperation{BaseOp: cloneBaseOp(s.BaseOp), Value: s.Value, IgnoreCase: s.IgnoreCase}, nil
}

// Clone returns an independent copy of the operation.
func (in *InOperation) Clone() (internal.Op, error) {
	return &InOperation{BaseOp: cloneBaseOp(in.BaseOp), Value: deepclone.Clone(in.Value)}, nil
}

// Clone returns an independent copy of the operation.
func (l *LessOperation) Clone() (internal.Op, error) {
	return &LessOperation{BaseOp: cloneBaseOp(l.BaseOp), Value: l.Value}, nil
}

// Clone returns an independent copy of the operation.
func (mo *MoreOperation) Clone() (internal.Op, error) {
	return &MoreOperation{BaseOp: cloneBaseOp(mo.BaseOp), Value: mo.Value}, nil
}

// Clone returns an independent copy of the operation.
func (ma *MatchesOperation) Clone() (internal.Op, error) {
	return &MatchesOperation{
		BaseOp:     cloneBaseOp(ma.BaseOp),
		Pattern:    ma.Pattern,
		IgnoreCase: ma.IgnoreCase,
		matcher:    ma.matcher,
		matcherErr: ma.matcherErr,
	}, nil
}

// Clone returns an independent copy of the operation.
func (ao *AndOperation) Clone() (internal.Op, error) {
	operations, err := clonePredicateOps(ao.Operations, ErrInvalidPredicateInAnd)
	if err != nil {
		return nil, err
	}
	return &AndOperation{BaseOp: cloneBaseOp(ao.BaseOp), Operations: operations}, nil
}

// Clone returns an independent copy of the operation.
func (oo *OrOperation) Clone() (internal.Op, error) {
	operations, err := clonePredicateOps(oo.Operations, ErrInvalidPredicateInOr)
	if err != nil {
		return nil, err
	}
	return &OrOperation{BaseOp: cloneBaseOp(oo.BaseOp), Operations: operations}, nil
}

// Clone returns an independent copy of the operation.
func (n *NotOperation) Clone() (internal.Op, error) {
	operations, err := clonePredicateOps(n.Operations, ErrInvalidPredicateInNot)
	if err != nil {
		return nil, err
	}
	return &NotOperation{BaseOp: cloneBaseOp(n.BaseOp), Operations: operations}, nil
}

// Clone returns an independent copy of the operation.
func (f *FlipOperation) Clone() (internal.Op, error) {
	return &FlipOperation{BaseOp: cloneBaseOp(f.BaseOp)}, nil
}

// Clone returns an independent copy of the operation.
func (ic *IncOperation) Clone() (internal.Op, error) {
	return &IncOperation{BaseOp: cloneBaseOp(ic.BaseOp), Inc: ic.Inc}, nil
}

// Clone returns an independent copy of the operation.
func (si *StrInsOperation) Clone() (internal.Op, error) {
	return &StrInsOperation{BaseOp: cloneBaseOp(si.BaseOp), Pos: si.Pos, Str: si.Str}, nil
}

// Clone returns an independent copy of the operation.
func (sd *StrDelOperation) Clone() (internal.Op, error) {
	return &StrDelOperation{
		BaseOp: cloneBaseOp(sd.BaseOp),
		Pos:    sd.Pos,
		Len:    sd.Len,
		Str:    sd.Str,
		HasStr: sd.HasStr,
	}, nil
}

// Clone returns an independent copy of the operation.
func (sp *SplitOperation) Clone() (internal.Op, error) {
	return &SplitOperation{BaseOp: cloneBaseOp(sp.BaseOp), Pos: sp.Pos, Props: cloneValue(sp.Props)}, nil
}

// Clone returns an independent copy of the operation.
func (mg *MergeOperation) Clone() (internal.Op, error) {
	return &MergeOperation{BaseOp: cloneBaseOp(mg.BaseOp), Pos: mg.Pos, Props: deepclone.Clone(mg.Props)}, nil
}

// Clone returns an independent copy of the operation.
func (ex *ExtendOperation) Clone() (internal.Op, error) {
	return &ExtendOperation{
		BaseOp:     cloneBaseOp(ex.BaseOp),
		Properties: deepclone.Clone(ex.Properties),
		DeleteNull: ex.DeleteNull,
	}, nil
}
