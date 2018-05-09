package Plan

import (
	"context"
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/xitongsys/guery/Util"
	"github.com/xitongsys/guery/parser"
)

type PlanFiliterNode struct {
	Input             PlanNode
	Metadata          *Util.Metadata
	BooleanExpression *BooleanExpressionNode
}

func NewPlanFiliterNode(input PlanNode, t parser.IBooleanExpressionContext) *PlanFiliterNode {
	res := &PlanFiliterNode{
		Input:             input,
		Metadata:          Util.NewDefaultMetadata,
		BooleanExpression: NewBooleanExpressionNode(t),
	}
	return res
}

func (self *PlanFiliterNode) GetNodeType() PlanNodeType {
	return FILTERNODE
}

func (self *PlanFiliterNode) String() string {
	res := "PlanFiliterNode {\n"
	res += "Input: " + self.Input.String() + "\n"
	res += "BooleanExpression: " + fmt.Sprint(self.BooleanExpression) + "\n"
	res += "}\n"
	return res
}