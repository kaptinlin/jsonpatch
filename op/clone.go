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
		cloneOp, ok := predicateOp.(internal.CloneOp)
		if !ok {
			return nil, fmt.Errorf("%w: predicate %T cannot be cloned", errInvalid, operation)
		}
		cloned, err := cloneOp.Clone()
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

// Clone implements internal.CloneOp.
func (a *AddOperation) Clone() (internal.Op, error) {
	return &AddOperation{BaseOp: cloneBaseOp(a.BaseOp), Value: cloneValue(a.Value)}, nil
}

// Clone implements internal.CloneOp.
func (r *RemoveOperation) Clone() (internal.Op, error) {
	return &RemoveOperation{
		BaseOp:      cloneBaseOp(r.BaseOp),
		OldValue:    cloneValue(r.OldValue),
		HasOldValue: r.HasOldValue,
	}, nil
}

// Clone implements internal.CloneOp.
func (rp *ReplaceOperation) Clone() (internal.Op, error) {
	return &ReplaceOperation{
		BaseOp:   cloneBaseOp(rp.BaseOp),
		Value:    cloneValue(rp.Value),
		OldValue: cloneValue(rp.OldValue),
	}, nil
}

// Clone implements internal.CloneOp.
func (m *MoveOperation) Clone() (internal.Op, error) {
	return &MoveOperation{BaseOp: cloneBaseOp(m.BaseOp)}, nil
}

// Clone implements internal.CloneOp.
func (c *CopyOperation) Clone() (internal.Op, error) {
	return &CopyOperation{BaseOp: cloneBaseOp(c.BaseOp)}, nil
}

// Clone implements internal.CloneOp.
func (t *TestOperation) Clone() (internal.Op, error) {
	return &TestOperation{BaseOp: cloneBaseOp(t.BaseOp), Value: cloneValue(t.Value), NotFlag: t.NotFlag}, nil
}

// Clone implements internal.CloneOp.
func (co *ContainsOperation) Clone() (internal.Op, error) {
	return &ContainsOperation{BaseOp: cloneBaseOp(co.BaseOp), Value: co.Value, IgnoreCase: co.IgnoreCase}, nil
}

// Clone implements internal.CloneOp.
func (d *DefinedOperation) Clone() (internal.Op, error) {
	return &DefinedOperation{BaseOp: cloneBaseOp(d.BaseOp)}, nil
}

// Clone implements internal.CloneOp.
func (u *UndefinedOperation) Clone() (internal.Op, error) {
	return &UndefinedOperation{BaseOp: cloneBaseOp(u.BaseOp)}, nil
}

// Clone implements internal.CloneOp.
func (tp *TypeOperation) Clone() (internal.Op, error) {
	return &TypeOperation{BaseOp: cloneBaseOp(tp.BaseOp), TypeValue: tp.TypeValue}, nil
}

// Clone implements internal.CloneOp.
func (tt *TestTypeOperation) Clone() (internal.Op, error) {
	return &TestTypeOperation{BaseOp: cloneBaseOp(tt.BaseOp), Types: deepclone.Clone(tt.Types)}, nil
}

// Clone implements internal.CloneOp.
func (ts *TestStringOperation) Clone() (internal.Op, error) {
	return &TestStringOperation{
		BaseOp:     cloneBaseOp(ts.BaseOp),
		Str:        ts.Str,
		Pos:        ts.Pos,
		NotFlag:    ts.NotFlag,
		IgnoreCase: ts.IgnoreCase,
	}, nil
}

// Clone implements internal.CloneOp.
func (tl *TestStringLenOperation) Clone() (internal.Op, error) {
	return &TestStringLenOperation{BaseOp: cloneBaseOp(tl.BaseOp), Length: tl.Length, NotFlag: tl.NotFlag}, nil
}

// Clone implements internal.CloneOp.
func (e *EndsOperation) Clone() (internal.Op, error) {
	return &EndsOperation{BaseOp: cloneBaseOp(e.BaseOp), Value: e.Value, IgnoreCase: e.IgnoreCase}, nil
}

// Clone implements internal.CloneOp.
func (s *StartsOperation) Clone() (internal.Op, error) {
	return &StartsOperation{BaseOp: cloneBaseOp(s.BaseOp), Value: s.Value, IgnoreCase: s.IgnoreCase}, nil
}

// Clone implements internal.CloneOp.
func (in *InOperation) Clone() (internal.Op, error) {
	return &InOperation{BaseOp: cloneBaseOp(in.BaseOp), Value: deepclone.Clone(in.Value)}, nil
}

// Clone implements internal.CloneOp.
func (l *LessOperation) Clone() (internal.Op, error) {
	return &LessOperation{BaseOp: cloneBaseOp(l.BaseOp), Value: l.Value}, nil
}

// Clone implements internal.CloneOp.
func (mo *MoreOperation) Clone() (internal.Op, error) {
	return &MoreOperation{BaseOp: cloneBaseOp(mo.BaseOp), Value: mo.Value}, nil
}

// Clone implements internal.CloneOp.
func (ma *MatchesOperation) Clone() (internal.Op, error) {
	return &MatchesOperation{
		BaseOp:     cloneBaseOp(ma.BaseOp),
		Pattern:    ma.Pattern,
		IgnoreCase: ma.IgnoreCase,
		matcher:    ma.matcher,
	}, nil
}

// Clone implements internal.CloneOp.
func (ao *AndOperation) Clone() (internal.Op, error) {
	operations, err := clonePredicateOps(ao.Operations, ErrInvalidPredicateInAnd)
	if err != nil {
		return nil, err
	}
	return &AndOperation{BaseOp: cloneBaseOp(ao.BaseOp), Operations: operations}, nil
}

// Clone implements internal.CloneOp.
func (oo *OrOperation) Clone() (internal.Op, error) {
	operations, err := clonePredicateOps(oo.Operations, ErrInvalidPredicateInOr)
	if err != nil {
		return nil, err
	}
	return &OrOperation{BaseOp: cloneBaseOp(oo.BaseOp), Operations: operations}, nil
}

// Clone implements internal.CloneOp.
func (n *NotOperation) Clone() (internal.Op, error) {
	operations, err := clonePredicateOps(n.Operations, ErrInvalidPredicateInNot)
	if err != nil {
		return nil, err
	}
	return &NotOperation{BaseOp: cloneBaseOp(n.BaseOp), Operations: operations}, nil
}

// Clone implements internal.CloneOp.
func (f *FlipOperation) Clone() (internal.Op, error) {
	return &FlipOperation{BaseOp: cloneBaseOp(f.BaseOp)}, nil
}

// Clone implements internal.CloneOp.
func (ic *IncOperation) Clone() (internal.Op, error) {
	return &IncOperation{BaseOp: cloneBaseOp(ic.BaseOp), Inc: ic.Inc}, nil
}

// Clone implements internal.CloneOp.
func (si *StrInsOperation) Clone() (internal.Op, error) {
	return &StrInsOperation{BaseOp: cloneBaseOp(si.BaseOp), Pos: si.Pos, Str: si.Str}, nil
}

// Clone implements internal.CloneOp.
func (sd *StrDelOperation) Clone() (internal.Op, error) {
	return &StrDelOperation{
		BaseOp: cloneBaseOp(sd.BaseOp),
		Pos:    sd.Pos,
		Len:    sd.Len,
		Str:    sd.Str,
		HasStr: sd.HasStr,
	}, nil
}

// Clone implements internal.CloneOp.
func (sp *SplitOperation) Clone() (internal.Op, error) {
	return &SplitOperation{BaseOp: cloneBaseOp(sp.BaseOp), Pos: sp.Pos, Props: cloneValue(sp.Props)}, nil
}

// Clone implements internal.CloneOp.
func (mg *MergeOperation) Clone() (internal.Op, error) {
	return &MergeOperation{BaseOp: cloneBaseOp(mg.BaseOp), Pos: mg.Pos, Props: deepclone.Clone(mg.Props)}, nil
}

// Clone implements internal.CloneOp.
func (ex *ExtendOperation) Clone() (internal.Op, error) {
	return &ExtendOperation{
		BaseOp:     cloneBaseOp(ex.BaseOp),
		Properties: deepclone.Clone(ex.Properties),
		DeleteNull: ex.DeleteNull,
	}, nil
}
