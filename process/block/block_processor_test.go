package block_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process"
	"github.com/numbatx/gn-coval-index/process/block"
	"github.com/numbatx/gn-coval-index/process/block/miniblocks"
	"github.com/numbatx/gn-coval-index/process/utility"
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-coval-index/testscommon/mock"
	"github.com/numbatx/gn-core/data"
	moaBlock "github.com/numbatx/gn-core/data/block"
	"github.com/numbatx/gn-core/data/indexer"
	"github.com/numbatx/gn-core/marshal"
	"github.com/stretchr/testify/require"
)

func TestBlockProcessor_NewBlockProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (marshal.Marshalizer, process.MiniBlockHandler)
		expectedErr error
	}{
		{
			args: func() (marshal.Marshalizer, process.MiniBlockHandler) {
				return nil, &mock.MiniBlockHandlerStub{}
			},
			expectedErr: covalent.ErrNilMarshaller,
		},
		{
			args: func() (marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.MarshallerStub{}, nil
			},
			expectedErr: covalent.ErrNilMiniBlockHandler,
		},
		{
			args: func() (marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := block.NewBlockProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestBlockProcessor_ProcessBlock_InvalidBodyAndHeaderMarshaller_ExpectProcessError(t *testing.T) {
	errMarshallHeader := errors.New("err header marshall")
	errMarshallBody := errors.New("err body marshall")

	tests := []struct {
		Marshaller  func(obj interface{}) ([]byte, error)
		expectedErr error
	}{
		{
			Marshaller: func(obj interface{}) ([]byte, error) {
				_, ok := obj.(*moaBlock.Header)
				if ok {
					return nil, errMarshallHeader
				}
				return json.Marshal(obj)
			},
			expectedErr: errMarshallHeader,
		},
		{
			Marshaller: func(obj interface{}) ([]byte, error) {
				_, ok := obj.(*moaBlock.Body)
				if ok {
					return nil, errMarshallBody
				}
				return json.Marshal(obj)
			},
			expectedErr: errMarshallBody,
		},
	}

	for _, currTest := range tests {
		bp, _ := block.NewBlockProcessor(
			&mock.MarshallerStub{MarshalCalled: currTest.Marshaller},
			&mock.MiniBlockHandlerStub{})

		args := getInitializedArgs(false)
		_, err := bp.ProcessBlock(args)

		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestBlockProcessor_ProcessBlock_InvalidBody_ExpectErrBlockBodyAssertion(t *testing.T) {
	mbp, _ := miniblocks.NewMiniBlocksProcessor(&mock.HasherMock{}, &mock.MarshallerStub{})
	bp, _ := block.NewBlockProcessor(&mock.MarshallerStub{}, mbp)

	args := getInitializedArgs(false)
	args.Body = nil
	_, err := bp.ProcessBlock(args)

	require.Equal(t, covalent.ErrBlockBodyAssertion, err)
}

func TestNewBlockProcessor_ProcessBlock_InvalidMBHandler_ExpectErr(t *testing.T) {
	errMBHandler := errors.New("error mb handler")

	bp, _ := block.NewBlockProcessor(
		&mock.MarshallerStub{},
		&mock.MiniBlockHandlerStub{
			ProcessMiniBlockCalled: func(header data.HeaderHandler, body data.BodyHandler) ([]*schema.MiniBlock, error) {
				return nil, errMBHandler
			}})

	args := getInitializedArgs(false)
	_, err := bp.ProcessBlock(args)

	require.Equal(t, errMBHandler, err)
}

func TestNewBlockProcessor_ProcessBlock_NoSigners_ExpectDefaultProposerIndex(t *testing.T) {
	bp, _ := block.NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})

	args := getInitializedArgs(false)
	args.SignersIndexes = nil
	ret, _ := bp.ProcessBlock(args)

	require.Equal(t, block.ProposerIndex, ret.Proposer)
}

func TestBlockProcessor_ProcessBlock(t *testing.T) {
	t.Parallel()

	bp, _ := block.NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})
	args := getInitializedArgs(false)
	ret, _ := bp.ProcessBlock(args)
	expectedNotarizedHeaderHashes, _ := utility.HexSliceToByteSlice(args.NotarizedHeadersHashes)

	require.Equal(t, int64(args.Header.GetNonce()), ret.Nonce)
	require.Equal(t, int64(args.Header.GetRound()), ret.Round)
	require.Equal(t, int32(args.Header.GetEpoch()), ret.Epoch)
	require.Equal(t, args.HeaderHash, ret.Hash)
	require.Equal(t, expectedNotarizedHeaderHashes, ret.NotarizedBlocksHashes)
	require.Equal(t, int64(args.SignersIndexes[0]), ret.Proposer)
	require.Equal(t, utility.UIntSliceToIntSlice(args.SignersIndexes), ret.Validators)
	require.Equal(t, args.Header.GetPubKeysBitmap(), ret.PubKeysBitmap)
	require.Equal(t, int64(485), ret.Size)
	require.Equal(t, int64(args.Header.GetTimeStamp()), ret.Timestamp)
	require.Equal(t, args.Header.GetRootHash(), ret.StateRootHash)
	require.Equal(t, args.Header.GetPrevHash(), ret.PrevHash)
	require.Equal(t, int32(args.Header.GetShardID()), ret.ShardID)
	require.Equal(t, int32(args.Header.GetTxCount()), ret.TxCount)
	require.Equal(t, args.Header.GetAccumulatedFees().Bytes(), ret.AccumulatedFees)
	require.Equal(t, args.Header.GetDeveloperFees().Bytes(), ret.DeveloperFees)

	require.Equal(t, ret.EpochStartInfo, (*schema.EpochStartInfo)(nil))
}

func TestBlockProcessor_ProcessMetaBlock(t *testing.T) {
	t.Parallel()

	bp, _ := block.NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})
	args := getInitializedArgs(true)
	ret, _ := bp.ProcessBlock(args)
	expectedNotarizedHeaderHashes, _ := utility.HexSliceToByteSlice(args.NotarizedHeadersHashes)

	require.Equal(t, int64(args.Header.GetNonce()), ret.Nonce)
	require.Equal(t, int64(args.Header.GetRound()), ret.Round)
	require.Equal(t, int32(args.Header.GetEpoch()), ret.Epoch)
	require.Equal(t, args.HeaderHash, ret.Hash)
	require.Equal(t, expectedNotarizedHeaderHashes, ret.NotarizedBlocksHashes)
	require.Equal(t, int64(args.SignersIndexes[0]), ret.Proposer)
	require.Equal(t, utility.UIntSliceToIntSlice(args.SignersIndexes), ret.Validators)
	require.Equal(t, args.Header.GetPubKeysBitmap(), ret.PubKeysBitmap)
	require.Equal(t, int64(args.Header.GetTimeStamp()), ret.Timestamp)
	require.Equal(t, args.Header.GetRootHash(), ret.StateRootHash)
	require.Equal(t, args.Header.GetPrevHash(), ret.PrevHash)
	require.Equal(t, int32(args.Header.GetShardID()), ret.ShardID)
	require.Equal(t, int32(args.Header.GetTxCount()), ret.TxCount)
	require.Equal(t, args.Header.GetAccumulatedFees().Bytes(), ret.AccumulatedFees)
	require.Equal(t, args.Header.GetDeveloperFees().Bytes(), ret.DeveloperFees)

	metaBlockEconomics := args.Header.(*moaBlock.MetaBlock).GetEpochStart().Economics

	require.Equal(t, metaBlockEconomics.TotalSupply.Bytes(), ret.EpochStartInfo.TotalSupply)
	require.Equal(t, metaBlockEconomics.TotalToDistribute.Bytes(), ret.EpochStartInfo.TotalToDistribute)
	require.Equal(t, metaBlockEconomics.TotalNewlyMinted.Bytes(), ret.EpochStartInfo.TotalNewlyMinted)
	require.Equal(t, metaBlockEconomics.RewardsPerBlock.Bytes(), ret.EpochStartInfo.RewardsPerBlock)
	require.Equal(t, metaBlockEconomics.RewardsForProtocolSustainability.Bytes(), ret.EpochStartInfo.RewardsForProtocolSustainability)
	require.Equal(t, metaBlockEconomics.NodePrice.Bytes(), ret.EpochStartInfo.NodePrice)
	require.Equal(t, int32(metaBlockEconomics.PrevEpochStartRound), ret.EpochStartInfo.PrevEpochStartRound)
	require.Equal(t, metaBlockEconomics.PrevEpochStartHash, ret.EpochStartInfo.PrevEpochStartHash)
}

func TestBlockProcessor_ProcessMetaBlock_NotStartOfEpochBlock_ExpectNilEpochStartInfo(t *testing.T) {
	bp, _ := block.NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})

	metaBlockHeader := getInitializedMetaBlockHeader()
	metaBlockHeader.EpochStart.LastFinalizedHeaders = nil

	ret, _ := bp.ProcessBlock(&indexer.ArgsSaveBlockData{
		Header: metaBlockHeader,
		Body:   &moaBlock.Body{}})

	require.Equal(t, (*schema.EpochStartInfo)(nil), ret.EpochStartInfo)
}

func getInitializedArgs(metaBlock bool) *indexer.ArgsSaveBlockData {
	var header data.HeaderHandler

	if metaBlock {
		header = getInitializedMetaBlockHeader()
	} else {
		header = getInitialisedBlockHeader()
	}

	return &indexer.ArgsSaveBlockData{
		HeaderHash:             []byte("header hash"),
		Body:                   &moaBlock.Body{},
		Header:                 header,
		SignersIndexes:         []uint64{1, 2, 3},
		NotarizedHeadersHashes: []string{"0a", "1f"},
		TransactionsPool:       nil,
	}
}

func getInitializedMetaBlockHeader() *moaBlock.MetaBlock {
	return &moaBlock.MetaBlock{
		Nonce:           1,
		Epoch:           2,
		Round:           3,
		TimeStamp:       4,
		LeaderSignature: []byte("meta leader signature"),
		PubKeysBitmap:   []byte("meta pub keys bitmap"),
		PrevHash:        []byte("meta prev hash"),
		PrevRandSeed:    []byte("meta prev rand seed"),
		RandSeed:        []byte("meta rand seed"),
		RootHash:        []byte("meta root hash"),
		EpochStart: moaBlock.EpochStart{
			LastFinalizedHeaders: []moaBlock.EpochStartShardData{{}},
			Economics: moaBlock.Economics{
				TotalSupply:                      big.NewInt(5),
				TotalToDistribute:                big.NewInt(6),
				TotalNewlyMinted:                 big.NewInt(7),
				RewardsPerBlock:                  big.NewInt(8),
				RewardsForProtocolSustainability: big.NewInt(9),
				NodePrice:                        big.NewInt(10),
				PrevEpochStartRound:              11,
				PrevEpochStartHash:               []byte("meta prev epoch hash"),
			},
		},
		ChainID:         []byte("meta chain id"),
		AccumulatedFees: big.NewInt(11),
		DeveloperFees:   big.NewInt(12),
		TxCount:         13,
	}
}

func getInitialisedBlockHeader() *moaBlock.Header {
	return &moaBlock.Header{
		Nonce:              1,
		PrevHash:           []byte("prev hash"),
		PrevRandSeed:       []byte("prev rand seed"),
		RandSeed:           []byte("rand seed"),
		PubKeysBitmap:      []byte("pub keys bitmap"),
		ShardID:            2,
		TimeStamp:          3,
		Round:              4,
		Epoch:              5,
		BlockBodyType:      6,
		LeaderSignature:    []byte("leader signature"),
		RootHash:           []byte("root hash"),
		TxCount:            7,
		EpochStartMetaHash: []byte("epoch start meta hash"),
		ReceiptsHash:       []byte("receipts hash"),
		ChainID:            []byte("chain id"),
		AccumulatedFees:    big.NewInt(8),
		DeveloperFees:      big.NewInt(9),
	}
}
