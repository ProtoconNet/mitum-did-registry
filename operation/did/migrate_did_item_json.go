package did

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type MigrateDIDItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address             `json:"contract"`
	PubKey   string                   `json:"publicKey"`
	TxID     string                   `json:"txid"`
	Currency currencytypes.CurrencyID `json:"currency"`
}

func (it MigrateDIDItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MigrateDIDItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		PubKey:     it.pubKey,
		TxID:       it.txID,
		Currency:   it.currency,
	})
}

type MigrateDIDItemJSONUnMarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	PubKey   string    `json:"publicKey"`
	TxID     string    `json:"txid"`
	Currency string    `json:"currency"`
}

func (it *MigrateDIDItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uit MigrateDIDItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc,
		uit.Hint,
		uit.Contract,
		uit.PubKey,
		uit.TxID,
		uit.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
