package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table                   string
	where                   []*Predicate
	args                    []any
	sb                      *strings.Builder
	lastSpaceIndexTag       int
	lastLeftBracketIndexTag int
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Build() (*Query, error) {
	sb := &strings.Builder{}
	s.sb = sb
	sb.WriteString("SELECT * FROM ")
	if s.table == "" {
		var t T
		sb.WriteString("`" + strings.ToLower(reflect.TypeOf(t).Name()) + "`")
	} else {
		sb.WriteString(s.table)
	}

	if len(s.where) > 0 {
		sb.WriteString(" WHERE")
		//合并多个条件
		p := s.where[0]
		for _, v := range s.where[1:] {
			p = p.AND(v)
		}
		s.BuildExpression(p)
	}

	sb.WriteByte(';')
	return &Query{
			SQL:  sb.String(),
			Args: s.args,
		},
		nil
}
func (s *Selector[T]) BuildExpression(expr Expression) (err error) {
	sb := s.sb
	switch exp := expr.(type) {
	case *Predicate:
		s.AddSpace()
		s.AddLeftBracket()
		if err = s.BuildExpression(exp.left); err != nil {
			return err
		}
		s.AddSpace()
		sb.WriteString(exp.op.String())
		s.AddSpace()
		if err = s.BuildExpression(exp.right); err != nil {
			return err
		}
		s.AddRightBracket()
	case *Column:
		sb.WriteByte('`')
		sb.WriteString(exp.name)
		sb.WriteByte('`')
		s.AddSpace()
	case *value:
		s.AddSpace()
		sb.WriteByte('?')
		//s.AddSpace()
		s.AddArg(exp.val)
	}
	return err
}
func (s *Selector[T]) AddArg(arg any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, arg)
}
func (s *Selector[T]) AddSpace() {
	if s.sb == nil {
		s.sb = &strings.Builder{}
		s.lastSpaceIndexTag = -1
	}
	if s.lastSpaceIndexTag != s.sb.Len()-1 && s.lastLeftBracketIndexTag != s.sb.Len()-1 {
		s.sb.WriteByte(' ')
		s.lastSpaceIndexTag = s.sb.Len() - 1
	}
}
func (s *Selector[T]) AddLeftBracket() {
	if s.sb == nil {
		s.sb = &strings.Builder{}
		s.lastLeftBracketIndexTag = -1
	}
	s.sb.WriteByte('(')
	s.lastLeftBracketIndexTag = s.sb.Len() - 1
}
func (s *Selector[T]) AddRightBracket() {
	if s.sb == nil {
		s.sb = &strings.Builder{}
		s.lastLeftBracketIndexTag = -1
	}
	s.sb.WriteByte(')')
}
func (s *Selector[T]) From(table string) *Selector[T] {
	strs := strings.Split(table, ".")
	if len(strs) == 2 {
		s.table = fmt.Sprintf("`%s`.`%s`", strs[0], strs[1])
	} else {
		s.table = fmt.Sprintf("`%s`", table)
	}
	return s
}
func (s *Selector[T]) Where(ps ...*Predicate) *Selector[T] {
	s.where = ps
	return s
}
