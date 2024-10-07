package did

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type RegisterModelFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender         mitumbase.Address `json:"sender"`
	Contract       mitumbase.Address `json:"contract"`
	DIDMethod      string            `json:"didMethod"`
	DocContext     string            `json:"docContext"`
	DocAuthType    string            `json:"docAuthType"`
	DocSvcType     string            `json:"docSvcType"`
	DocSvcEndPoint string            `json:"docSvcEndPoint"`
	Currency       types.CurrencyID  `json:"currency"`
}

func (fact RegisterModelFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterModelFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		DIDMethod:             fact.didMethod,
		DocContext:            fact.docContext,
		DocAuthType:           fact.docAuthType,
		DocSvcType:            fact.docSvcType,
		DocSvcEndPoint:        fact.docSvcEndPoint,
		Currency:              fact.currency,
	})
}

type RegisterModelFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender         string `json:"sender"`
	Contract       string `json:"contract"`
	DIDMethod      string `json:"didMethod"`
	DocContext     string `json:"docContext"`
	DocAuthType    string `json:"docAuthType"`
	DocSvcType     string `json:"docSvcType"`
	DocSvcEndPoint string `json:"docSvcEndPoint"`
	Currency       string `json:"currency"`
}

func (fact *RegisterModelFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u RegisterModelFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unpack(
		enc, u.Sender, u.Contract, u.DIDMethod, u.DocContext, u.DocAuthType, u.DocSvcType, u.DocSvcEndPoint, u.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op RegisterModel) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(op.BaseOperation.JSONMarshaler())
}

func (op *RegisterModel) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
