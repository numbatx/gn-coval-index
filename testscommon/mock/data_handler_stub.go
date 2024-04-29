package mock

import (
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-core/data/indexer"
)

type DataHandlerStub struct {
	ProcessDataCalled func(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error)
}

func (dhs *DataHandlerStub) ProcessData(args *indexer.ArgsSaveBlockData) (*schema.BlockResult, error) {
	if dhs.ProcessDataCalled != nil {
		return dhs.ProcessDataCalled(args)
	}
	return nil, nil
}
