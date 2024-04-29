package accounts_test

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process"
	"github.com/numbatx/gn-coval-index/process/accounts"
	"github.com/numbatx/gn-coval-index/process/utility"
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-coval-index/testscommon"
	"github.com/numbatx/gn-coval-index/testscommon/mock"
	"github.com/numbatx/gn-core/core"
	vmcommon "github.com/numbatx/gn-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewAccountsProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter)
		expectedErr error
	}{
		{
			args: func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter) {
				return nil, &mock.AccountsAdapterStub{}, &mock.PubKeyConverterStub{}
			},
			expectedErr: covalent.ErrNilShardCoordinator,
		},
		{
			args: func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter) {
				return &mock.ShardCoordinatorMock{}, nil, &mock.PubKeyConverterStub{}
			},
			expectedErr: covalent.ErrNilAccountsAdapter,
		},
		{
			args: func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter) {
				return &mock.ShardCoordinatorMock{}, &mock.AccountsAdapterStub{}, nil
			},
			expectedErr: covalent.ErrNilPubKeyConverter,
		},
		{
			args: func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter) {
				return &mock.ShardCoordinatorMock{}, &mock.AccountsAdapterStub{}, &mock.PubKeyConverterStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := accounts.NewAccountsProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestAccountsProcessor_ProcessAccounts_InvalidUserAccountHandler_ExpectZeroAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{
			LoadAccountCalled: func(address []byte) (vmcommon.AccountHandler, error) {
				return nil, nil
			}},
		&mock.PubKeyConverterStub{})

	tx := &schema.Transaction{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes()}
	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 0)
}

func TestAccountsProcessor_ProcessAccounts_InvalidLoadAccount_ExpectZeroAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{
			LoadAccountCalled: func(address []byte) (vmcommon.AccountHandler, error) {
				return nil, errors.New("load account error")
			}},
		&mock.PubKeyConverterStub{})

	tx := &schema.Transaction{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes()}
	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 0)
}

func TestAccountsProcessor_ProcessAccounts_NotInSameShard_ExpectZeroAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{SelfID: 4},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountMock{}},
		&mock.PubKeyConverterStub{})

	tx := &schema.Transaction{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes()}
	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 0)
}

func TestAccountsProcessor_ProcessAccounts_OneSender_NilReceiver_ExpectOneAccount(t *testing.T) {
	addresses := generateAddresses(1)
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountMock{}},
		&mock.PubKeyConverterStub{
			DecodeCalled: func(humanReadable string) ([]byte, error) {
				if len(humanReadable) == 0 {
					return nil, errors.New("nil address")
				}
				return make([]byte, 0), nil
			},
		})

	tx := &schema.Transaction{
		Sender:   addresses[0],
		Receiver: nil,
	}

	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 1)
	checkProcessedAccounts(t, addresses, ret)
}

func TestAccountsProcessor_ProcessAccounts_NilSender_OneReceiver_ExpectOneAccount(t *testing.T) {
	addresses := generateAddresses(1)
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountMock{}},
		&mock.PubKeyConverterStub{
			DecodeCalled: func(humanReadable string) ([]byte, error) {
				if len(humanReadable) == 0 {
					return nil, errors.New("nil address")
				}
				return make([]byte, 0), nil
			},
		})

	tx := &schema.Transaction{
		Sender:   nil,
		Receiver: addresses[0],
	}

	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 1)
	checkProcessedAccounts(t, addresses, ret)
}

func TestAccountsProcessor_ProcessAccounts_FourAddresses_TwoIdentical_ExpectTwoAccounts(t *testing.T) {
	addresses := generateAddresses(2)
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountMock{}},
		&mock.PubKeyConverterStub{})

	tx1 := &schema.Transaction{
		Sender:   addresses[0],
		Receiver: addresses[1],
	}
	tx2 := &schema.Transaction{
		Sender:   addresses[1],
		Receiver: addresses[0],
	}

	ret := ap.ProcessAccounts([]*schema.Transaction{tx1, tx2}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 2)
	checkProcessedAccounts(t, addresses, ret)
}

func TestAccountsProcessor_ProcessAccounts_OneAddress_OneMetaChainShardAddress_ExpectOneAccount(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountMock{}},
		&mock.PubKeyConverterStub{})

	tx := &schema.Transaction{
		Receiver: []byte("adr1"),
		Sender:   utility.MetaChainShardAddress()}

	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 1)
	require.Equal(t, []byte("adr1"), ret[0].Address)
	require.Equal(t, big.NewInt(int64(1)).Bytes(), ret[0].Balance)
	require.Equal(t, int64(1), ret[0].Nonce)
}

func TestAccountsProcessor_ProcessAccounts_TwoAddresses_OneInvalidBech32_ExpectOneAccount(t *testing.T) {
	invalidAddress := "adr_invalid"
	errPubKeyDecode := errors.New("error invalid address")

	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountMock{}},
		&mock.PubKeyConverterStub{
			DecodeCalled: func(humanReadable string) ([]byte, error) {
				if humanReadable == invalidAddress {
					return nil, errPubKeyDecode
				}
				return make([]byte, 0), nil
			},
		})

	tx := &schema.Transaction{
		Receiver: []byte("adr1"),
		Sender:   []byte(invalidAddress)}

	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 1)
	require.Equal(t, []byte("adr1"), ret[0].Address)
	require.Equal(t, big.NewInt(int64(1)).Bytes(), ret[0].Balance)
	require.Equal(t, int64(1), ret[0].Nonce)
}

func TestAccountsProcessor_ProcessAccounts_SevenAddresses_ExpectSevenAccounts(t *testing.T) {
	addresses := generateAddresses(7)

	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountMock{}},
		&mock.PubKeyConverterStub{})

	tx1 := &schema.Transaction{
		Receiver: addresses[0],
		Sender:   addresses[1],
	}
	tx2 := &schema.Transaction{
		Sender:   addresses[2],
		Receiver: addresses[3],
	}
	txs := []*schema.Transaction{tx1, tx2}

	scr := &schema.SCResult{
		Sender:   addresses[4],
		Receiver: addresses[5],
	}
	scrs := []*schema.SCResult{scr}

	receipt := &schema.Receipt{
		Sender: addresses[6],
	}
	receipts := []*schema.Receipt{receipt}

	ret := ap.ProcessAccounts(txs, scrs, receipts)

	require.Len(t, ret, 7)
	checkProcessedAccounts(t, addresses, ret)
}

func generateAddresses(n int) [][]byte {
	addresses := make([][]byte, n)

	for i := 0; i < n; i++ {
		addresses[i] = []byte("adr" + strconv.Itoa(i))
	}

	return addresses
}

// This function only works if accounts.NewAccountsProcessor is called with
// &mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountMock{}}
func checkProcessedAccounts(t *testing.T, addresses [][]byte, processedAcc []*schema.AccountBalanceUpdate) {
	require.Equal(t, len(addresses), len(processedAcc), "should have the same number of processed accounts as initial addresses")

	allProcessedAddr := make(map[string]struct{})

	for _, currAccount := range processedAcc {
		allProcessedAddr[string(currAccount.Address)] = struct{}{}
	}

	for _, addr := range addresses {
		_, exists := allProcessedAddr[string(addr)]
		require.True(t, exists, fmt.Sprintf("%s not processed successfully", addr))
	}

	for idx, account := range processedAcc {
		require.Equal(t, big.NewInt(int64(idx+1)).Bytes(), account.Balance)
		require.Equal(t, int64(idx+1), account.Nonce)
	}
}
