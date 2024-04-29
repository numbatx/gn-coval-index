package block

import (
	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process"
	"github.com/numbatx/gn-coval-index/process/utility"
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-core/core"
	"github.com/numbatx/gn-core/core/check"
	"github.com/numbatx/gn-core/data"
	moaBlock "github.com/numbatx/gn-core/data/block"
	"github.com/numbatx/gn-core/data/indexer"
	"github.com/numbatx/gn-core/marshal"
)

const ProposerIndex = int64(0)

type blockProcessor struct {
	marshaller        marshal.Marshalizer
	miniBlocksHandler process.MiniBlockHandler
}

// NewBlockProcessor creates a new instance of block processor
func NewBlockProcessor(marshaller marshal.Marshalizer, mbHandler process.MiniBlockHandler) (*blockProcessor, error) {
	if check.IfNil(marshaller) {
		return nil, covalent.ErrNilMarshaller
	}
	if mbHandler == nil {
		return nil, covalent.ErrNilMiniBlockHandler
	}

	return &blockProcessor{
		marshaller:        marshaller,
		miniBlocksHandler: mbHandler,
	}, nil
}

// ProcessBlock converts block data to a specific structure defined by avro schema
func (bp *blockProcessor) ProcessBlock(args *indexer.ArgsSaveBlockData) (*schema.Block, error) {
	blockSizeInBytes, err := bp.computeBlockSize(args.Header, args.Body)
	if err != nil {
		return nil, err
	}

	miniBlocks, err := bp.miniBlocksHandler.ProcessMiniBlocks(args.Header, args.Body)
	if err != nil {
		return nil, err
	}

	notarizedBlockHashes, err := utility.HexSliceToByteSlice(args.NotarizedHeadersHashes)
	if err != nil {
		return nil, err
	}

	header := args.Header
	return &schema.Block{
		Nonce:                 int64(header.GetNonce()),
		Round:                 int64(header.GetRound()),
		Epoch:                 int32(header.GetEpoch()),
		Hash:                  args.HeaderHash,
		MiniBlocks:            miniBlocks,
		NotarizedBlocksHashes: notarizedBlockHashes,
		Proposer:              getProposerIndex(args.SignersIndexes),
		Validators:            utility.UIntSliceToIntSlice(args.SignersIndexes),
		PubKeysBitmap:         header.GetPubKeysBitmap(),
		Size:                  blockSizeInBytes,
		Timestamp:             int64(header.GetTimeStamp()),
		StateRootHash:         header.GetRootHash(),
		PrevHash:              header.GetPrevHash(),
		ShardID:               int32(header.GetShardID()),
		TxCount:               int32(header.GetTxCount()),
		AccumulatedFees:       utility.GetBytes(header.GetAccumulatedFees()),
		DeveloperFees:         utility.GetBytes(header.GetDeveloperFees()),
		EpochStartBlock:       header.IsStartOfEpochBlock(),
		EpochStartInfo:        getEpochStartInfo(header),
	}, nil
}

func (bp *blockProcessor) computeBlockSize(header data.HeaderHandler, body data.BodyHandler) (int64, error) {
	headerBytes, err := bp.marshaller.Marshal(header)
	if err != nil {
		return 0, err
	}
	bodyBytes, err := bp.marshaller.Marshal(body)
	if err != nil {
		return 0, err
	}

	blockSize := len(headerBytes) + len(bodyBytes)

	return int64(blockSize), nil
}

func getProposerIndex(signersIndexes []uint64) int64 {
	if len(signersIndexes) > 0 {
		return int64(signersIndexes[ProposerIndex])
	}

	return ProposerIndex
}

func getEpochStartInfo(header data.HeaderHandler) *schema.EpochStartInfo {
	if header.GetShardID() != core.MetachainShardId {
		return nil
	}

	metaHeader, ok := header.(*moaBlock.MetaBlock)
	if !ok {
		return nil
	}

	if !metaHeader.IsStartOfEpochBlock() {
		return nil
	}

	economics := metaHeader.EpochStart.Economics

	return &schema.EpochStartInfo{
		TotalSupply:                      utility.GetBytes(economics.TotalSupply),
		TotalToDistribute:                utility.GetBytes(economics.TotalToDistribute),
		TotalNewlyMinted:                 utility.GetBytes(economics.TotalNewlyMinted),
		RewardsPerBlock:                  utility.GetBytes(economics.RewardsPerBlock),
		RewardsForProtocolSustainability: utility.GetBytes(economics.RewardsForProtocolSustainability),
		NodePrice:                        utility.GetBytes(economics.NodePrice),
		PrevEpochStartRound:              int32(economics.PrevEpochStartRound),
		PrevEpochStartHash:               economics.PrevEpochStartHash,
	}
}
