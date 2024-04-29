package miniblocks_test

import (
	"errors"
	"testing"

	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process/block/miniblocks"
	"github.com/numbatx/gn-coval-index/testscommon/mock"
	"github.com/numbatx/gn-core/data/block"
	"github.com/numbatx/gn-core/hashing"
	"github.com/numbatx/gn-core/marshal"
	"github.com/stretchr/testify/require"
)

func TestNewMiniBlocksProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (hashing.Hasher, marshal.Marshalizer)
		expectedErr error
	}{
		{
			args: func() (hashing.Hasher, marshal.Marshalizer) {
				return nil, &mock.MarshallerStub{}
			},
			expectedErr: covalent.ErrNilHasher,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer) {
				return &mock.HasherMock{}, nil
			},
			expectedErr: covalent.ErrNilMarshaller,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer) {
				return &mock.HasherMock{}, &mock.MarshallerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := miniblocks.NewMiniBlocksProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestMiniBlocksProcessor_ProcessMiniBlocks(t *testing.T) {
	mbp, _ := miniblocks.NewMiniBlocksProcessor(&mock.HasherMock{}, &mock.MarshallerStub{})

	header := &block.Header{TimeStamp: 123}
	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{[]byte("x"), []byte("y")},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            3},
		{
			TxHashes:        [][]byte{[]byte("y"), []byte("z")},
			ReceiverShardID: 4,
			SenderShardID:   5,
			Type:            6},
	}}

	ret, _ := mbp.ProcessMiniBlocks(header, body)

	require.Len(t, ret, 2)

	require.Equal(t, []byte("ok"), ret[0].Hash)
	require.Equal(t, [][]byte{[]byte("x"), []byte("y")}, ret[0].TxHashes)
	require.Equal(t, int64(123), ret[0].Timestamp)
	require.Equal(t, int32(1), ret[0].ReceiverShardID)
	require.Equal(t, int32(2), ret[0].SenderShardID)
	require.Equal(t, int32(3), ret[0].Type)

	require.Equal(t, []byte("ok"), ret[1].Hash)
	require.Equal(t, [][]byte{[]byte("y"), []byte("z")}, ret[1].TxHashes)
	require.Equal(t, int64(123), ret[1].Timestamp)
	require.Equal(t, int32(4), ret[1].ReceiverShardID)
	require.Equal(t, int32(5), ret[1].SenderShardID)
	require.Equal(t, int32(6), ret[1].Type)
}

func TestMiniBlocksProcessor_ProcessMiniBlocks_InvalidMarshaller_ExpectZeroMBProcessed(t *testing.T) {
	mbp, _ := miniblocks.NewMiniBlocksProcessor(
		&mock.HasherMock{},
		&mock.MarshallerStub{
			MarshalCalled: func(obj interface{}) ([]byte, error) {
				return nil, errors.New("error marshaller stub")
			},
		})

	header := &block.Header{TimeStamp: 123}
	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            3},
		{
			ReceiverShardID: 4,
			SenderShardID:   5,
			Type:            6},
	}}

	ret, err := mbp.ProcessMiniBlocks(header, body)

	require.Nil(t, err)
	require.Len(t, ret, 0)
}
