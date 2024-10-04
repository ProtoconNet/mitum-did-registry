package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (d Data) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":     d.Hint().String(),
		"publicKey": d.pubKey,
		"did":       d.did,
	})
}

type DataBSONUnmarshaler struct {
	Hint      string `bson:"_hint"`
	PublicKey string `bson:"publicKey"`
	DID       string `bson:"did"`
}

func (d *Data) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of Item")

	var u DataBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return d.unmarshal(ht, u.PublicKey, u.DID)
}
