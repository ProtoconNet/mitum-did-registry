package did

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	crcytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-did-registry/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
	"strings"
)

var MigrateDIDItemHint = hint.MustNewHint("mitum-did-migrate-did-item-v0.0.1")

type MigrateDIDItem struct {
	hint.BaseHinter
	contract base.Address
	pubKey   string
	txID     string
	currency crcytypes.CurrencyID
}

func NewMigrateDIDItem(
	contract base.Address,
	pubKey string,
	txID string,
	currency crcytypes.CurrencyID,
) MigrateDIDItem {
	return MigrateDIDItem{
		BaseHinter: hint.NewBaseHinter(MigrateDIDItemHint),
		contract:   contract,
		pubKey:     pubKey,
		txID:       txID,
		currency:   currency,
	}
}

func (it MigrateDIDItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		[]byte(it.pubKey),
		[]byte(it.txID),
		it.currency.Bytes(),
	)
}

func (it MigrateDIDItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
	); err != nil {
		return common.ErrItemInvalid.Wrap(err)
	}

	pubKey := strings.TrimPrefix(it.pubKey, "0x")
	if len(pubKey) < types.MinKeyLen {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(
				errors.Errorf("invalid pub key length %v < %v", len(pubKey), types.MinKeyLen)))
	}

	if !crcytypes.ReValidSpcecialCh.Match([]byte(pubKey)) {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Wrap(
				errors.Errorf("pub key %s, must match regex `^[^\\s:/?#\\[\\]$@]*$`", pubKey)))
	}

	return nil
}

func (it MigrateDIDItem) Contract() base.Address {
	return it.contract
}

func (it MigrateDIDItem) PubKey() string {
	return it.pubKey
}

func (fact MigrateDIDItem) PubKeyDetatched() string {
	// Detach 0x
	pubKey := strings.TrimPrefix(fact.pubKey, "0x")

	return pubKey
}

func (fact MigrateDIDItem) PubKeyReformed() string {
	pubKey := fact.PubKeyDetatched()
	// reform pubkey
	pubKey = "04" + pubKey[len(pubKey)-128:]
	return pubKey
}

func (it MigrateDIDItem) TxID() string {
	return it.txID
}

func (it MigrateDIDItem) Currency() crcytypes.CurrencyID {
	return it.currency
}

func (it MigrateDIDItem) Addresses() []base.Address {
	ad := make([]base.Address, 1)

	ad[0] = it.contract

	return ad
}
