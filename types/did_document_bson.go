package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (d Document) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint": d.Hint().String(),
		"did_doc": bson.M{
			"@context": d.didDoc.context_,
			"id":       d.didDoc.id,
			"created":  d.didDoc.created,
			"status":   d.didDoc.status,
			"authentication": bson.M{
				"id":           d.didDoc.authentication.id,
				"type":         d.didDoc.authentication.authType,
				"controller":   d.didDoc.authentication.controller,
				"publicKeyHex": d.didDoc.authentication.publicKeyHex,
			},
			"service": bson.M{
				"id":                d.didDoc.service.id,
				"type":              d.didDoc.service.serviceType,
				"service_end_point": d.didDoc.service.serviceEndPoint,
			},
		},
	})
}

//func (d DIDDocument) MarshalBSON() ([]byte, error) {
//	authBytes, err := d.authentication.MarshalBSON()
//	if err != nil {
//		return nil, err
//	}
//	serviceBytes, err := d.service.MarshalBSON()
//	if err != nil {
//		return nil, err
//	}
//	return bsonenc.Marshal(bson.M{
//		"@context":       d.context_,
//		"id":             d.id,
//		"created":        d.created,
//		"status":         d.status,
//		"authentication": authBytes,
//		"service":        serviceBytes,
//	})
//}
//
//func (d Authentication) MarshalBSON() ([]byte, error) {
//	return bsonenc.Marshal(bson.M{
//		"id":           d.id,
//		"type":         d.authType,
//		"controller":   d.controller,
//		"publicKeyHex": d.publicKeyHex,
//	})
//}
//
//func (d Service) MarshalBSON() ([]byte, error) {
//	return bsonenc.Marshal(bson.M{
//		"id":                d.id,
//		"type":              d.serviceType,
//		"service_end_point": d.serviceEndPoint,
//	})
//}

type DocumentBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	DIDDoc bson.Raw `bson:"did_doc"`
}

func (d *Document) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of Document")

	var u DocumentBSONUnmarshaler
	if err := bson.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}
	var ud DIDDocumentBSONUnmarshaler
	if err := bson.Unmarshal(u.DIDDoc, &ud); err != nil {
		return e.Wrap(err)
	}

	var ua AuthenticationBSONUnmarshaler
	if err := bson.Unmarshal(ud.Auth, &ua); err != nil {
		return e.Wrap(err)
	}

	var us ServiceBSONUnmarshaler
	if err := bson.Unmarshal(ud.Service, &us); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return d.unmarshal(ht, ud.Context_, ud.ID, ud.Created, ud.Status,
		ua.ID, ua.Type, ua.Controller, ua.PublicKeyHex,
		us.ID, us.Type, us.ServiceEndPoint,
	)
}

type DIDDocumentBSONUnmarshaler struct {
	Context_ string   `bson:"@context"`
	ID       string   `bson:"id"`
	Created  string   `bson:"created"`
	Status   string   `bson:"status"`
	Auth     bson.Raw `bson:"authentication"`
	Service  bson.Raw `bson:"service"`
}

type AuthenticationBSONUnmarshaler struct {
	ID           string `bson:"id"`
	Type         string `bson:"type"`
	Controller   string `bson:"controller"`
	PublicKeyHex string `bson:"publicKeyHex"`
}

type ServiceBSONUnmarshaler struct {
	ID              string `bson:"id"`
	Type            string `bson:"type"`
	ServiceEndPoint string `bson:"service_end_point"`
}
