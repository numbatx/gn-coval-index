package receipts

import (
	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process/utility"
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-core/core"
	"github.com/numbatx/gn-core/core/check"
	"github.com/numbatx/gn-core/data"
	"github.com/numbatx/gn-core/data/receipt"
)

type receiptsProcessor struct {
	pubKeyConverter core.PubkeyConverter
}

// NewReceiptsProcessor creates a new instance of receipts processor
func NewReceiptsProcessor(pubKeyConverter core.PubkeyConverter) (*receiptsProcessor, error) {
	if check.IfNil(pubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}

	return &receiptsProcessor{
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// ProcessReceipts converts receipts data to a specific structure defined by avro schema
func (rp *receiptsProcessor) ProcessReceipts(receipts map[string]data.TransactionHandler, timeStamp uint64) []*schema.Receipt {
	allReceipts := make([]*schema.Receipt, 0, len(receipts))

	for currHash, currReceipt := range receipts {
		rec := rp.processReceipt(currReceipt, currHash, timeStamp)
		if rec != nil {
			allReceipts = append(allReceipts, rec)
		}
	}

	return allReceipts
}

func (rp *receiptsProcessor) processReceipt(
	tx data.TransactionHandler,
	receiptHash string,
	timeStamp uint64,
) *schema.Receipt {

	rec, castOk := tx.(*receipt.Receipt)
	if !castOk {
		return nil
	}

	return &schema.Receipt{
		Hash:      []byte(receiptHash),
		Value:     utility.GetBytes(rec.GetValue()),
		Sender:    utility.EncodePubKey(rp.pubKeyConverter, rec.GetSndAddr()),
		Data:      rec.GetData(),
		TxHash:    rec.GetTxHash(),
		Timestamp: int64(timeStamp),
	}
}
