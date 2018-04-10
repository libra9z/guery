package Plan

import (
	"github.com/xitongsys/guery/Context"
	"github.com/xitongsys/guery/DataSource"
	"github.com/xitongsys/guery/parser"
)

type BooleanValueNode struct {
	Bool bool
}

func NewBooleanValueNode(ctx *Context.Context, t parser.IBooleanValueContext) *BooleanValueNode {
	s := t.GetText()
	b := true
	if s != "TRUE" {
		b = false
	}
	return &BooleanValueNode{
		Tree: t,
		Bool: b,
	}
}

func (self *BooleanValueNode) Result(intput DataSource.DataSource) bool {
	return self.Bool
}
