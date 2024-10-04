package types

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DocumentJSONMarshaler struct {
	hint.BaseHinter
	DIDDoc DIDDocumentMarshaler `json:"did_doc"`
}

func (d Document) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DocumentJSONMarshaler{
		BaseHinter: d.BaseHinter,
		DIDDoc: DIDDocumentMarshaler{
			Context_: d.didDoc.context_,
			ID:       d.didDoc.id,
			Created:  d.didDoc.created,
			Status:   d.didDoc.status,
			Auth: AuthenticationMarshaler{
				ID:           d.didDoc.authentication.id,
				Type:         d.didDoc.authentication.authType,
				Controller:   d.didDoc.authentication.controller,
				PublicKeyHex: d.didDoc.authentication.publicKeyHex,
			},
			Service: ServiceMarshaler{
				ID:              d.didDoc.service.id,
				Type:            d.didDoc.service.serviceType,
				ServiceEndPoint: d.didDoc.service.serviceEndPoint,
			},
		},
	})
}

type DocumentJSONUnmarshaler struct {
	Hint   hint.Hint            `json:"_hint"`
	DIDDoc DIDDocumentMarshaler `json:"did_doc"`
}

func (d *Document) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Data")

	var u DocumentJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}
	var uf DocumentJSONUnmarshaler
	if err := json.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return d.unmarshal(
		u.Hint, uf.DIDDoc.Context_, uf.DIDDoc.ID, uf.DIDDoc.Created, uf.DIDDoc.Status,
		uf.DIDDoc.Auth.ID, uf.DIDDoc.Auth.Type, uf.DIDDoc.Auth.Controller, uf.DIDDoc.Auth.PublicKeyHex,
		uf.DIDDoc.Service.ID, uf.DIDDoc.Service.Type, uf.DIDDoc.Service.ServiceEndPoint,
	)
}

type DIDDocumentMarshaler struct {
	Context_ string                  `json:"@context"`
	ID       string                  `json:"id"`
	Created  string                  `json:"created"`
	Status   string                  `json:"status"`
	Auth     AuthenticationMarshaler `json:"authentication"`
	Service  ServiceMarshaler        `json:"service"`
}

type AuthenticationMarshaler struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Controller   string `json:"controller"`
	PublicKeyHex string `json:"publicKeyHex"`
}

type ServiceMarshaler struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndPoint string `json:"service_end_point"`
}
