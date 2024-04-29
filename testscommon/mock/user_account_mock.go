package mock

import (
	"math/big"
	"strconv"
)

// UserAccountMock -
type UserAccountMock struct {
	CurrentBalance int64
	CurrentNonce   uint64
}

// IncreaseNonce -
func (uas *UserAccountMock) IncreaseNonce(_ uint64) {
}

// GetBalance increments CurrentBalance and returns it as a big int
func (uas *UserAccountMock) GetBalance() *big.Int {
	uas.CurrentBalance++
	return big.NewInt(uas.CurrentBalance)
}

// AddressBytes returns a byte slice of ("addr" + CurrentBalance)
func (uas *UserAccountMock) AddressBytes() []byte {
	return []byte("addr" + strconv.Itoa(int(uas.CurrentBalance)))
}

// GetNonce increments CurrentNonce and returns it
func (uas *UserAccountMock) GetNonce() uint64 {
	uas.CurrentNonce++
	return uas.CurrentNonce
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (uas *UserAccountMock) IsInterfaceNil() bool {
	return uas == nil
}

// RetrieveValueFromDataTrieTracker -
func (uas *UserAccountMock) RetrieveValueFromDataTrieTracker([]byte) ([]byte, error) {
	return nil, nil
}
