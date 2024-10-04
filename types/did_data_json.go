package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DataJSONMarshaler struct {
	hint.BaseHinter
	PublicKey string `json:"publicKey"`
	DID       string `json:"did"`
}

func (d Data) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DataJSONMarshaler{
		BaseHinter: d.BaseHinter,
		PublicKey:  d.pubKey,
		DID:        d.did,
	})
}

type DataJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	PublicKey string    `json:"publicKey"`
	DID       string    `json:"did"`
}

func (d *Data) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Data")

	var u DataJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return d.unmarshal(u.Hint, u.PublicKey, u.DID)
}
