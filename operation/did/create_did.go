package did

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-did-registry/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
	"strings"
)

var (
	CreateDIDFactHint = hint.MustNewHint("mitum-did-create-did-operation-fact-v0.0.1")
	CreateDIDHint     = hint.MustNewHint("mitum-did-create-did-operation-v0.0.1")
)

type CreateDIDFact struct {
	mitumbase.BaseFact
	sender   mitumbase.Address
	contract mitumbase.Address
	pubKey   string
	currency currencytypes.CurrencyID
}

func NewCreateDIDFact(
	token []byte, sender, contract mitumbase.Address,
	pubKey string, currency currencytypes.CurrencyID) CreateDIDFact {
	bf := mitumbase.NewBaseFact(CreateDIDFactHint, token)
	fact := CreateDIDFact{
		BaseFact: bf,
		sender:   sender,
		contract: contract,
		pubKey:   pubKey,
		currency: currency,
	}

	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact CreateDIDFact) IsValid(b []byte) error {
	pubKey := strings.TrimPrefix(fact.pubKey, "0x")
	if len(pubKey) < types.MinKeyLen {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(
				errors.Errorf("invalid pub key length %v < %v", len(fact.pubKey), types.MinKeyLen)))
	}

	if !currencytypes.ReValidSpcecialCh.Match([]byte(pubKey)) {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Wrap(
				errors.Errorf("pub key %s, must match regex `^[^\\s:/?#\\[\\]$@]*$`", pubKey)))
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact CreateDIDFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateDIDFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateDIDFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		[]byte(fact.pubKey),
		fact.currency.Bytes(),
	)
}

func (fact CreateDIDFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact CreateDIDFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact CreateDIDFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact CreateDIDFact) PubKey() string {
	return fact.pubKey
}

func (fact CreateDIDFact) PubKeyDetatched() string {
	// Detach 0x
	pubKey := strings.TrimPrefix(fact.pubKey, "0x")

	return pubKey
}

func (fact CreateDIDFact) PubKeyReformed() string {
	pubKey := fact.PubKeyDetatched()
	// reform pubkey
	pubKey = "04" + pubKey[len(pubKey)-128:]
	return pubKey
}

func (fact CreateDIDFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact CreateDIDFact) Addresses() ([]mitumbase.Address, error) {
	as := []mitumbase.Address{fact.sender}

	return as, nil
}

type CreateDID struct {
	common.BaseOperation
}

func NewCreateDID(fact CreateDIDFact) (CreateDID, error) {
	return CreateDID{BaseOperation: common.NewBaseOperation(CreateDIDHint, fact)}, nil
}
