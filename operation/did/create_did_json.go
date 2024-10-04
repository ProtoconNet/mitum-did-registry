package did

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type CreateDIDFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender   mitumbase.Address `json:"sender"`
	Contract mitumbase.Address `json:"contract"`
	PubKey   string            `json:"publicKey"`
	Currency types.CurrencyID  `json:"currency"`
}

func (fact CreateDIDFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateDIDFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		PubKey:                fact.pubKey,
		Currency:              fact.currency,
	})
}

type CreateDIDFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Contract string `json:"contract"`
	PubKey   string `json:"publicKey"`
	Currency string `json:"currency"`
}

func (fact *CreateDIDFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u CreateDIDFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, u.Sender, u.Contract, u.PubKey, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op CreateDID) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		op.BaseOperation.JSONMarshaler(),
	)
}

func (op *CreateDID) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
