package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-did-registry/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (sv DesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  sv.Hint().String(),
			"design": sv.Design,
		},
	)
}

type DesignStateValueBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Design bson.Raw `bson:"design"`
}

func (sv *DesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of DesignStateValue")

	var u DesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	sv.BaseHinter = hint.NewBaseHinter(ht)

	var sd types.Design
	if err := sd.DecodeBSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}
	sv.Design = sd

	return nil
}

func (sv DataStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    sv.Hint().String(),
			"did_data": sv.Data,
		},
	)
}

type DataStateValueBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	DIDData bson.Raw `bson:"did_data"`
}

func (sv *DataStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of DataStateValue")

	var u DataStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	sv.BaseHinter = hint.NewBaseHinter(ht)

	var n types.Data
	if err := n.DecodeBSON(u.DIDData, enc); err != nil {
		return e.Wrap(err)
	}
	sv.Data = n

	return nil
}

func (sv DocumentStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":        sv.Hint().String(),
			"did_document": sv.Document,
		},
	)
}

type DocumentStateValueBSONUnmarshaler struct {
	Hint        string   `bson:"_hint"`
	DIDDocument bson.Raw `bson:"did_document"`
}

func (sv *DocumentStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of DocumentStateValue")

	var u DocumentStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	sv.BaseHinter = hint.NewBaseHinter(ht)

	var n types.Document
	if err := n.DecodeBSON(u.DIDDocument, enc); err != nil {
		return e.Wrap(err)
	}
	sv.Document = n

	return nil
}
