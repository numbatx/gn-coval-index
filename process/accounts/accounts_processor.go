package accounts

import (
	"bytes"

	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process"
	"github.com/numbatx/gn-coval-index/process/utility"
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-core/core"
	"github.com/numbatx/gn-core/core/check"
	"github.com/numbatx/gn-core/data"
	logger "github.com/numbatx/gn-logger"
)

var log = logger.GetOrCreate("covalent/process/accounts")

type accountsProcessor struct {
	shardCoordinator process.ShardCoordinator
	pubKeyConverter  core.PubkeyConverter
	accounts         covalent.AccountsAdapter
}

// NewAccountsProcessor creates a new instance of accounts processor
func NewAccountsProcessor(
	shardCoordinator process.ShardCoordinator,
	accounts covalent.AccountsAdapter,
	pubKeyConverter core.PubkeyConverter,
) (*accountsProcessor, error) {

	if check.IfNil(shardCoordinator) {
		return nil, covalent.ErrNilShardCoordinator
	}
	if check.IfNil(accounts) {
		return nil, covalent.ErrNilAccountsAdapter
	}
	if check.IfNil(pubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}

	return &accountsProcessor{
		accounts:         accounts,
		pubKeyConverter:  pubKeyConverter,
		shardCoordinator: shardCoordinator,
	}, nil
}

// ProcessAccounts converts accounts data to a specific structure defined by avro schema
func (ap *accountsProcessor) ProcessAccounts(
	processedTxs []*schema.Transaction,
	processedSCRs []*schema.SCResult,
	processedReceipts []*schema.Receipt,
) []*schema.AccountBalanceUpdate {
	addresses := ap.getAllAddresses(processedTxs, processedSCRs, processedReceipts)
	accounts := make([]*schema.AccountBalanceUpdate, 0, len(addresses))

	for address := range addresses {
		account, err := ap.processAccount(address)
		if err != nil || account == nil {
			log.Warn("cannot get account address", "address", address, "error", err)
			continue
		}

		accounts = append(accounts, account)
	}

	return accounts
}

func (ap *accountsProcessor) getAllAddresses(
	processedTxs []*schema.Transaction,
	processedSCRs []*schema.SCResult,
	processedReceipts []*schema.Receipt,
) map[string]struct{} {
	addresses := make(map[string]struct{})

	for _, tx := range processedTxs {
		ap.addAddressIfInSelfShard(addresses, tx.Sender)
		ap.addAddressIfInSelfShard(addresses, tx.Receiver)
	}

	for _, scr := range processedSCRs {
		ap.addAddressIfInSelfShard(addresses, scr.Sender)
		ap.addAddressIfInSelfShard(addresses, scr.Receiver)
	}

	for _, receipt := range processedReceipts {
		ap.addAddressIfInSelfShard(addresses, receipt.Sender)
	}

	return addresses
}

func (ap *accountsProcessor) addAddressIfInSelfShard(addresses map[string]struct{}, address []byte) {
	if bytes.Equal(address, utility.MetaChainShardAddress()) {
		return
	}
	if ap.shardCoordinator.SelfId() == ap.shardCoordinator.ComputeId(address) {
		addresses[string(address)] = struct{}{}
	}
}

func (ap *accountsProcessor) processAccount(address string) (*schema.AccountBalanceUpdate, error) {
	//TODO: This only works as long as covalent indexer is part of numbat node binary.
	// This needs to be changed, so that account content is given as an input parameter, not loaded.
	pubKey, err := ap.pubKeyConverter.Decode(address)
	if err != nil {
		return nil, err
	}

	acc, err := ap.accounts.LoadAccount(pubKey)
	if err != nil {
		return nil, err
	}

	account, castOk := acc.(data.UserAccountHandler)
	if !castOk {
		return nil, covalent.ErrCannotCastAccountHandlerToUserAccount
	}
	return &schema.AccountBalanceUpdate{
		Address: []byte(address),
		Balance: utility.GetBytes(account.GetBalance()),
		Nonce:   int64(account.GetNonce()),
	}, nil
}
