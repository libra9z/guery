package plan

import (
	"github.com/xitongsys/guery/config"
	"github.com/xitongsys/guery/metadata"
)

type PlanUseNode struct {
	Catalog, Schema string
}

func NewPlanUseNode(runtime *config.ConfigRuntime, ct, sh string) *PlanUseNode {
	return &PlanUseNode{
		Catalog: ct,
		Schema:  sh,
	}
}

func (self *PlanUseNode) GetNodeType() PlanNodeType {
	return USENODE
}

func (self *PlanUseNode) SetMetadata() error {
	return nil
}

func (self *PlanUseNode) GetMetadata() *metadata.Metadata {
	return nil
}

func (self *PlanUseNode) GetOutput() PlanNode {
	return nil
}

func (self *PlanUseNode) SetOutput(output PlanNode) {
	return
}

func (self *PlanUseNode) GetInputs() []PlanNode {
	return nil
}

func (self *PlanUseNode) SetInputs(input []PlanNode) {
	return
}

func (self *PlanUseNode) String() string {
	res := "PlanUseNode  {\n"
	res += self.Catalog + "." + self.Schema
	res += "}\n"
	return res
}
