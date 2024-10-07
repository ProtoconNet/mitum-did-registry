package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":          de.Hint().String(),
			"didMethod":      de.didMethod,
			"docContext":     de.docContext,
			"docAuthType":    de.docAuthType,
			"docSvcType":     de.docSvcType,
			"docSvcEndPoint": de.docSvcEndPoint,
		})
}

type DesignBSONUnmarshaler struct {
	Hint           string `bson:"_hint"`
	DIDMethod      string `bson:"didMethod"`
	DocContext     string `bson:"docContext"`
	DocAuthType    string `bson:"docAuthType"`
	DocSvcType     string `bson:"docSvcType"`
	DocSvcEndPoint string `bson:"docSvcEndPoint"`
}

func (de *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of Design")

	var u DesignBSONUnmarshaler
	if err := bson.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return de.unmarshal(enc, ht, u.DIDMethod, u.DocContext, u.DocAuthType, u.DocSvcType, u.DocSvcEndPoint)

}
