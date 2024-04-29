package mock

import (
	vmcommon "github.com/numbatx/gn-vm-common"
)

// AccountsAdapterStub -
type AccountsAdapterStub struct {
	UserAccountHandler vmcommon.AccountHandler
	LoadAccountCalled  func(address []byte) (vmcommon.AccountHandler, error)
}

// LoadAccount calls a custom load account function if defined, otherwise returns UserAccountMock, nil
func (aas *AccountsAdapterStub) LoadAccount(address []byte) (vmcommon.AccountHandler, error) {
	if aas.LoadAccountCalled != nil {
		return aas.LoadAccountCalled(address)
	}
	return aas.UserAccountHandler, nil
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (aas *AccountsAdapterStub) IsInterfaceNil() bool {
	return aas == nil
}
