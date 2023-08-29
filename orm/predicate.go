package orm

type operator string

const (
	OpADD operator = "AND"
	OpOR  operator = "OR"
	OpNOT operator = "NOT"
	OpEQ  operator = "="
	OpLT  operator = "<"
	OpGT  operator = ">"
)

func (o operator) String() string {
	return string(o)
}

type Predicate struct {
	left  Expression
	op    operator
	right Expression
}

func (p *Predicate) expr() {}

func (p *Predicate) AND(right *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    OpADD,
		right: right,
	}
}
func (p *Predicate) OR(right *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    OpOR,
		right: right,
	}
}
func NOT(right *Predicate) *Predicate {
	return &Predicate{
		op:    OpNOT,
		right: right,
	}
}

type Expression interface {
	expr()
}
