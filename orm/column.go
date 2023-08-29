package orm

type Column struct {
	name string
}

func C(name string) *Column {
	return &Column{
		name: name,
	}
}

func (c *Column) EQ(arg any) *Predicate {
	return &Predicate{
		left:  c,
		op:    OpEQ,
		right: valueOf(arg),
	}
}

// GT 大于
func (c *Column) GT(arg any) *Predicate {
	return &Predicate{
		left:  c,
		op:    OpGT,
		right: valueOf(arg),
	}
}

// LT 小于
func (c *Column) LT(arg any) *Predicate {
	return &Predicate{
		left:  c,
		op:    OpLT,
		right: valueOf(arg),
	}
}

func (c *Column) expr() {}

type value struct {
	val any
}

func valueOf(val any) *value {
	return &value{val: val}
}
func (v *value) expr() {}
