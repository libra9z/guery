package executor

import (
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"time"

	"github.com/vmihailenco/msgpack"
	"github.com/xitongsys/guery/eplan"
	"github.com/xitongsys/guery/logger"
	"github.com/xitongsys/guery/metadata"
	"github.com/xitongsys/guery/pb"
	"github.com/xitongsys/guery/row"
	"github.com/xitongsys/guery/util"
)

func (self *Executor) SetInstructionAggregateFuncGlobal(instruction *pb.Instruction) (err error) {
	var enode eplan.EPlanAggregateFuncGlobalNode
	if err = msgpack.Unmarshal(instruction.EncodedEPlanNodeBytes, &enode); err != nil {
		return err
	}
	self.Instruction = instruction
	self.EPlanNode = &enode
	self.InputLocations = []*pb.Location{}
	for i := 0; i < len(enode.Inputs); i++ {
		self.InputLocations = append(self.InputLocations, &enode.Inputs[i])
	}
	self.OutputLocations = []*pb.Location{&enode.Output}
	return nil
}

func (self *Executor) RunAggregateFuncGlobal() (err error) {
	fname := fmt.Sprintf("executor_%v_aggregatefunclocal_%v_cpu.pprof", self.Name, time.Now().Format("20060102150405"))
	f, _ := os.Create(fname)
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	defer func() {
		if err != nil {
			self.AddLogInfo(err, pb.LogLevel_ERR)
		}
		self.Clear()
	}()

	writer := self.Writers[0]
	enode := self.EPlanNode.(*eplan.EPlanAggregateFuncGlobalNode)
	md := &metadata.Metadata{}

	//read md
	for _, reader := range self.Readers {
		if err = util.ReadObject(reader, md); err != nil {
			return err
		}
	}

	//write md
	if err = util.WriteObject(writer, enode.Metadata); err != nil {
		return err
	}

	rbWriter := row.NewRowsBuffer(enode.Metadata, nil, writer)

	defer func() {
		rbWriter.Flush()
	}()

	//init
	if err := enode.Init(enode.Metadata); err != nil {
		return err
	}

	//write rows
	var rg *row.RowsGroup
	res := make([]map[string]interface{}, len(enode.FuncNodes))
	for i := 0; i < len(res); i++ {
		res[i] = make(map[string]interface{})
	}

	keys := map[string]*row.Row{}

	for _, reader := range self.Readers { //TODO: concurrent?
		rbReader := row.NewRowsBuffer(md, reader, nil)
		for {
			rg, err = rbReader.Read()

			if err == io.EOF {
				err = nil
				break
			}
			if err != nil {
				break
			}

			for i := 0; i < rg.GetRowsNumber(); i++ {
				key := rg.GetKeyString(i)
				if _, ok := keys[key]; !ok {
					keys[key] = rg.GetRow(i)
				}
			}

			if err = self.CalAggregateFuncGlobal(enode, rg, &res); err != nil {
				break
			}
		}
	}

	for key, row := range keys {
		for i := 0; i < len(res); i++ {
			row.Vals[len(row.Vals)-i-1] = res[i][key]
		}
		rbWriter.WriteRow(row)
	}

	logger.Infof("RunAggregateFuncGlobal finished")
	return err
}

func (self *Executor) CalAggregateFuncGlobal(enode *eplan.EPlanAggregateFuncGlobalNode, rg *row.RowsGroup, res *[]map[string]interface{}) error {
	var err error
	var resc map[string]interface{}
	var resci interface{}
	for i, item := range enode.FuncNodes {
		if resci, err = item.Result(rg); err != nil {
			break
		}
		resc = resci.(map[string]interface{})
		for k, v := range resc {
			(*res)[i][k] = v
		}
	}
	return err
}
