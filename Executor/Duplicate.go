package Executor

import (
	"io"

	"github.com/vmihailenco/msgpack"
	"github.com/xitongsys/guery/EPlan"
	"github.com/xitongsys/guery/Util"
	"github.com/xitongsys/guery/pb"
)

func (self *Executor) SetInstructionDuplicate(instruction *pb.Instruction) (err error) {
	var enode EPlan.EPlanDuplicateNode
	if err = msgpack.Unmarshal(instruction.EncodedEPlanNodeBytes, &enode); err != nil {
		return err
	}
	self.Instruction = instruction
	self.EPlanNode = &enode
	self.InputLocations = []*pb.Location{}
	for i := 0; i < len(enode.Inputs); i++ {
		loc := enode.Inputs[i]
		self.InputLocations = append(self.InputLocations, &loc)
	}
	self.OutputLocations = []*pb.Location{}
	for i := 0; i < len(enode.Outputs); i++ {
		loc := enode.Outputs[i]
		self.OutputLocations = append(self.OutputLocations, &loc)
	}
	return nil
}

func (self *Executor) RunDuplicate() (err error) {
	defer self.Clear()
	enode := self.EPlanNode.(*EPlan.EPlanDuplicateNode)
	//read md
	md := &Util.Metadata{}
	for _, reader := range self.Readers {
		if err = Util.ReadObject(reader, md); err != nil {
			return err
		}
	}

	//write md
	if enode.Keys != nil && len(enode.Keys) > 0 {
		md.ClearKeys()
		md.AppendKeyByType(Util.STRING)
	}
	for _, writer := range self.Writers {
		if err = Util.WriteObject(writer, md); err != nil {
			return err
		}
	}

	rbWriters := make([]*Util.RowsBuffer, len(self.Writers))
	for i, writer := range self.Writers {
		rbWriters[i] = Util.NewRowsBuffer(md, nil, writer)
	}

	//write rows
	var row *Util.Row
	for _, reader := range self.Readers {
		rbReader := Util.NewRowsBuffer(md, reader, nil)
		for {
			row, err = rbReader.ReadRow()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			if enode.Keys != nil && len(enode.Keys) > 0 {
				rg := Util.NewRowsGroup(md)
				rg.Write(row)
				key, err := CalHashKey(enode.Keys, rg)
				if err != nil {
					return err
				}
				row.AppendKeys(key)
			}

			for _, rbWriter := range rbWriters {
				if err = rbWriter.WriteRow(row); err != nil {
					return err
				}
			}
		}
	}

	for _, rbWriter := range rbWriters {
		if err = rbWriter.Flush(); err != nil {
			return err
		}
	}
	return nil
}
