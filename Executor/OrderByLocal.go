package Executor

import (
	"io"

	"github.com/vmihailenco/msgpack"
	"github.com/xitongsys/guery/EPlan"
	"github.com/xitongsys/guery/Logger"
	"github.com/xitongsys/guery/Plan"
	"github.com/xitongsys/guery/Util"
	"github.com/xitongsys/guery/pb"
)

func (self *Executor) SetInstructionOrderByLocal(instruction *pb.Instruction) (err error) {
	var enode EPlan.EPlanOrderByLocalNode
	if err = msgpack.Unmarshal(instruction.EncodedEPlanNodeBytes, &enode); err != nil {
		return err
	}
	self.Instruction = instruction
	self.EPlanNode = &enode
	self.InputLocations = []*pb.Location{&enode.Input}
	self.OutputLocations = []*pb.Location{&enode.Output}
	return nil
}

func (self *Executor) RunOrderByLocal() (err error) {
	defer self.Clear()

	reader, writer := self.Readers[0], self.Writers[0]
	md := &Util.Metadata{}
	//read md
	if err = Util.ReadObject(reader, md); err != nil {
		return err
	}

	//write md
	if err = Util.WriteObject(writer, md); err != nil {
		return err
	}

	//write rows
	var row *Util.Row
	rows := Util.NewRows()
	enode := self.EPlanNode.(*EPlan.EPlanOrderByLocalNode)
	for {
		row, err = Util.ReadRow(reader)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return err
		}
		rb := Util.NewRowsBuffer(md)
		rb.Write(row)
		key, err := self.CalSortKey(rb)
		if err != nil {
			return err
		}
		row.Key = key
		rows.Append(row)
	}

	switch enode.OrderType {
	case Plan.ASC:
		rows.SortASC()
	case Plan.DESC:
		rows.SortDesc()
	}

	Util.WriteEOFMessage(writer)
	Logger.Infof("RunOrderByLocal finished")
	return nil
}

func (self *Executor) CalSortKey(enode *EPlan.EPlanOrderByLocalNode, rowsBuf *Util.RowsBuffer) (string, error) {
	var err error
	res := ""
	for _, item := range enode.SortItems {
		key, err := item.Result(rowsBuf)
		if err == io.EOF {
			return res, nil
		}
		if err != nil {
			return "", err
		}
		res = res + fmt.Sprintf("%v", key)
	}

	return res, err

}
