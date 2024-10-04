package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	DIDMethod      string `json:"didMethod"`
	DocContext     string `json:"docContext"`
	DocAuthType    string `json:"docAuthType"`
	DocSvcType     string `json:"docSvcType"`
	DocSvcEncPoint string `json:"docSvcEncPoint"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter:     de.BaseHinter,
		DIDMethod:      de.didMethod,
		DocContext:     de.docContext,
		DocAuthType:    de.docAuthType,
		DocSvcType:     de.docSvcType,
		DocSvcEncPoint: de.docSvcEncPoint,
	})
}

type DesignJSONUnmarshaler struct {
	Hint           hint.Hint `json:"_hint"`
	DIDMethod      string    `json:"didMethod"`
	DocContext     string    `json:"docContext"`
	DocAuthType    string    `json:"docAuthType"`
	DocSvcType     string    `json:"docSvcType"`
	DocSvcEncPoint string    `json:"docSvcEncPoint"`
}

func (de *Design) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Design")

	var u DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return de.unmarshal(enc, u.Hint, u.DIDMethod, u.DocContext, u.DocAuthType, u.DocSvcType, u.DocSvcEncPoint)
}
