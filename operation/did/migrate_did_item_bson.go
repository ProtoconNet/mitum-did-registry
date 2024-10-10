package did

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it MigrateDIDItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     it.Hint().String(),
			"contract":  it.contract,
			"publicKey": it.pubKey,
			"txid":      it.txID,
			"currency":  it.currency,
		},
	)
}

type MigrateDIDItemBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Contract string `bson:"contract"`
	PubKey   string `bson:"publicKey"`
	TxID     string `bson:"txid"`
	Currency string `bson:"currency"`
}

func (it *MigrateDIDItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit MigrateDIDItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unpack(enc, ht,
		uit.Contract,
		uit.PubKey,
		uit.TxID,
		uit.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	return nil
}
