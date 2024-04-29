package transactions_test

import (
	"math/rand"
	"testing"

	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process/transactions"
	"github.com/numbatx/gn-coval-index/process/utility"
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-coval-index/testscommon"
	"github.com/numbatx/gn-coval-index/testscommon/mock"
	"github.com/numbatx/gn-core/core"
	"github.com/numbatx/gn-core/data"
	"github.com/numbatx/gn-core/data/smartContractResult"
	"github.com/numbatx/gn-core/data/vm"
	"github.com/stretchr/testify/require"
)

func TestNewSCProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() core.PubkeyConverter
		expectedErr error
	}{
		{
			args: func() core.PubkeyConverter {
				return nil
			},
			expectedErr: covalent.ErrNilPubKeyConverter,
		},
		{
			args: func() core.PubkeyConverter {
				return &mock.PubKeyConverterStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := transactions.NewSCResultsProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func generateRandomSCR() *smartContractResult.SmartContractResult {
	return &smartContractResult.SmartContractResult{
		Nonce:          rand.Uint64(),
		Value:          testscommon.GenerateRandomBigInt(),
		RcvAddr:        testscommon.GenerateRandomBytes(),
		SndAddr:        testscommon.GenerateRandomBytes(),
		RelayerAddr:    testscommon.GenerateRandomBytes(),
		RelayedValue:   testscommon.GenerateRandomBigInt(),
		Code:           testscommon.GenerateRandomBytes(),
		Data:           testscommon.GenerateRandomBytes(),
		PrevTxHash:     testscommon.GenerateRandomBytes(),
		OriginalTxHash: testscommon.GenerateRandomBytes(),
		GasLimit:       rand.Uint64(),
		GasPrice:       rand.Uint64(),
		CallType:       vm.CallType(rand.Int()),
		CodeMetadata:   testscommon.GenerateRandomBytes(),
		ReturnMessage:  testscommon.GenerateRandomBytes(),
		OriginalSender: testscommon.GenerateRandomBytes(),
	}
}

func TestScProcessor_ProcessSCs_TwoSCRs_OneNormalTx_ExpectTwoProcessedSCRs(t *testing.T) {
	// TODO refactor this test: the processing is done by iterating on a map and the result is a slice that might,
	// sometimes, have the transactions in another order

	scp, _ := transactions.NewSCResultsProcessor(&mock.PubKeyConverterStub{})

	tx1 := generateRandomSCR()
	tx2 := generateRandomSCR()
	tx3 := generateRandomTx()

	txPool := map[string]data.TransactionHandler{
		"hash1": tx1,
		"hash2": tx2,
		"hash3": tx3,
	}

	ret := scp.ProcessSCRs(txPool, 123)

	require.Len(t, ret, 2)
	requireProcessedSCREqual(t, ret[0], tx1, "hash1", 123, &mock.PubKeyConverterStub{})
	requireProcessedSCREqual(t, ret[1], tx2, "hash2", 123, &mock.PubKeyConverterStub{})
}

func requireProcessedSCREqual(
	t *testing.T,
	processedSCR *schema.SCResult,
	scr *smartContractResult.SmartContractResult,
	hash string,
	timeStamp uint64,
	pubKeyConverter core.PubkeyConverter) {

	require.Equal(t, []byte(hash), processedSCR.Hash)
	require.Equal(t, int64(scr.GetNonce()), processedSCR.Nonce)
	require.Equal(t, int64(scr.GetGasLimit()), processedSCR.GasLimit)
	require.Equal(t, int64(scr.GetGasPrice()), processedSCR.GasPrice)
	require.Equal(t, scr.GetValue().Bytes(), processedSCR.Value)
	require.Equal(t, utility.EncodePubKey(pubKeyConverter, scr.GetSndAddr()), processedSCR.Sender)
	require.Equal(t, utility.EncodePubKey(pubKeyConverter, scr.GetRcvAddr()), processedSCR.Receiver)
	require.Equal(t, utility.EncodePubKey(pubKeyConverter, scr.GetRelayerAddr()), processedSCR.RelayerAddr)
	require.Equal(t, scr.GetRelayedValue().Bytes(), processedSCR.RelayedValue)
	require.Equal(t, scr.GetCode(), processedSCR.Code)
	require.Equal(t, scr.GetData(), processedSCR.Data)
	require.Equal(t, scr.GetPrevTxHash(), processedSCR.PrevTxHash)
	require.Equal(t, scr.GetOriginalTxHash(), processedSCR.OriginalTxHash)
	require.Equal(t, int32(scr.GetCallType()), processedSCR.CallType)
	require.Equal(t, scr.GetCodeMetadata(), processedSCR.CodeMetadata)
	require.Equal(t, scr.GetReturnMessage(), processedSCR.ReturnMessage)
	require.Equal(t, int64(timeStamp), processedSCR.Timestamp)
}
