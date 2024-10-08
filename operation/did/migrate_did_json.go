package did

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type MigrateDIDFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address     `json:"sender"`
	Items  []MigrateDIDItem `json:"items"`
}

func (fact MigrateDIDFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MigrateDIDFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type MIgrateDIDFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *MigrateDIDFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf MIgrateDIDFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Sender, uf.Items); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type MigrateDIDMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op MigrateDID) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MigrateDIDMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *MigrateDID) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
