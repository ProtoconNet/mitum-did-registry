package types

import (
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (d *Document) unmarshal(
	ht hint.Hint,
	context_, docID, created, status,
	authid, authType, controller, publicKeyHex,
	serviceID, serviceType, serviceEndpoint string,
) error {
	d.BaseHinter = hint.NewBaseHinter(ht)

	d.didDoc = DIDDocument{
		context_: context_,
		id:       docID,
		created:  created,
		status:   status,
		authentication: Authentication{
			id:           authid,
			authType:     authType,
			controller:   controller,
			publicKeyHex: publicKeyHex,
		},
		service: Service{
			id:              serviceID,
			serviceType:     serviceType,
			serviceEndPoint: serviceEndpoint,
		},
	}

	return nil
}
