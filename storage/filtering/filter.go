package filtering

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Filter struct {
	Field string
	Op    string
	Value interface{}
}

func (f *Filter) GetExp() exp.Expression {
	switch f.Op {
	case ">":
		return goqu.C(f.Field).Gt(f.Value)
	case "=":
		return goqu.C(f.Field).Eq(f.Value)
	}
	return exp.Default()
}
