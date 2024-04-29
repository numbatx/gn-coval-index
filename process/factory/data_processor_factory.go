package factory

import (
	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process"
	"github.com/numbatx/gn-coval-index/process/accounts"
	blockCovalent "github.com/numbatx/gn-coval-index/process/block"
	"github.com/numbatx/gn-coval-index/process/block/miniblocks"
	"github.com/numbatx/gn-coval-index/process/logs"
	"github.com/numbatx/gn-coval-index/process/receipts"
	"github.com/numbatx/gn-coval-index/process/transactions"
	"github.com/numbatx/gn-core/core"
	"github.com/numbatx/gn-core/hashing"
	"github.com/numbatx/gn-core/marshal"
)

// ArgsDataProcessor holds all input dependencies required by data processor factory
// in order to create a new data handler instance of type data processor
type ArgsDataProcessor struct {
	PubKeyConvertor  core.PubkeyConverter
	Accounts         covalent.AccountsAdapter
	Hasher           hashing.Hasher
	Marshaller       marshal.Marshalizer
	ShardCoordinator process.ShardCoordinator
}

// CreateDataProcessor creates a new data handler instance of type data processor
func CreateDataProcessor(args *ArgsDataProcessor) (covalent.DataHandler, error) {
	miniBlocksHandler, err := miniblocks.NewMiniBlocksProcessor(args.Hasher, args.Marshaller)
	if err != nil {
		return nil, err
	}

	blockHandler, err := blockCovalent.NewBlockProcessor(args.Marshaller, miniBlocksHandler)
	if err != nil {
		return nil, err
	}

	transactionsHandler, err := transactions.NewTransactionProcessor(args.PubKeyConvertor, args.Hasher, args.Marshaller)
	if err != nil {
		return nil, err
	}

	receiptsHandler, err := receipts.NewReceiptsProcessor(args.PubKeyConvertor)
	if err != nil {
		return nil, err
	}

	scResultsHandler, err := transactions.NewSCResultsProcessor(args.PubKeyConvertor)
	if err != nil {
		return nil, err
	}

	logHandler, err := logs.NewLogsProcessor(args.PubKeyConvertor)
	if err != nil {
		return nil, err
	}

	accountsHandler, err := accounts.NewAccountsProcessor(args.ShardCoordinator, args.Accounts, args.PubKeyConvertor)
	if err != nil {
		return nil, err
	}

	return process.NewDataProcessor(
		blockHandler,
		transactionsHandler,
		scResultsHandler,
		receiptsHandler,
		logHandler,
		accountsHandler)
}
