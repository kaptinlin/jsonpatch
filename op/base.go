package op

// BaseOp provides common functionality for all operations.
type BaseOp struct {
	path []string
	from []string
}

// NewBaseOp creates a new BaseOp with the given path.
func NewBaseOp(path []string) BaseOp {
	return BaseOp{path: path}
}

// NewBaseOpWithFrom creates a new BaseOp with path and from path.
func NewBaseOpWithFrom(path, from []string) BaseOp {
	return BaseOp{path: path, from: from}
}

// Path returns the operation path.
func (b *BaseOp) Path() []string {
	return b.path
}

// From returns the from path for move/copy operations.
func (b *BaseOp) From() []string {
	return b.from
}

// HasFrom returns true if the operation has a from path.
func (b *BaseOp) HasFrom() bool {
	return len(b.from) > 0
}

// PredicateOpBase provides common functionality for predicate operations.
type PredicateOpBase struct {
	BaseOp
	not bool
}

// NewPredicateOp creates a new PredicateOpBase.
func NewPredicateOp(path []string, not bool) PredicateOpBase {
	return PredicateOpBase{
		BaseOp: NewBaseOp(path),
		not:    not,
	}
}

// Not returns the negation flag.
func (p *PredicateOpBase) Not() bool {
	return p.not
}
