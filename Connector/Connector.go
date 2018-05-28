package Connector

import (
	"fmt"
	"strings"

	"github.com/xitongsys/guery/Connector/FileConnector"
	"github.com/xitongsys/guery/Connector/HiveConnector"
	"github.com/xitongsys/guery/Connector/TestConnector"
	"github.com/xitongsys/guery/Util"
)

type Connector interface {
	GetMetadata() *Util.Metadata
	GetPartitionInfo() *Util.PartitionInfo
	Read() (*Util.Row, error)
	ReadByColumns(colIndexes []int) (*Util.Row, error)
	ReadPartitionByColumns(parIndex int, colIndexes []int) (*Util.Row, error)
}

func NewConnector(catalog string, schema string, table string) (Connector, error) {
	catalog, schema, table = strings.ToUpper(catalog), strings.ToUpper(schema), strings.ToUpper(table)
	switch catalog {
	case "TEST":
		return TestConnector.NewTestConnector(schema, table)
	case "FILE":
		return FileConnector.NewFileConnector(schema, table)
	case "HIVE":
		return HiveConnector.NewHiveConnector(schema, table)

	}
	return nil, fmt.Errorf("NewConnector failed")
}