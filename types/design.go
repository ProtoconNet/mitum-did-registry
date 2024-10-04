package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var DesignHint = hint.MustNewHint("mitum-did-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	didMethod      string
	docContext     string
	docAuthType    string
	docSvcType     string
	docSvcEncPoint string
}

func NewDesign(didMethod, docContext, docAuthType, docSvcType, docSvcEndpoint string) Design {
	return Design{
		BaseHinter:     hint.NewBaseHinter(DesignHint),
		didMethod:      didMethod,
		docContext:     docContext,
		docAuthType:    docAuthType,
		docSvcType:     docSvcType,
		docSvcEncPoint: docSvcEndpoint,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
	); err != nil {
		return err
	}

	return nil
}

func (de Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(de.didMethod),
		[]byte(de.docContext),
		[]byte(de.docAuthType),
		[]byte(de.docSvcType),
		[]byte(de.docSvcEncPoint),
	)
}

func (de Design) Hash() util.Hash {
	return de.GenerateHash()
}

func (de Design) GenerateHash() util.Hash {
	return valuehash.NewSHA256(de.Bytes())
}

func (de Design) DIDMethod() string {
	return de.didMethod
}

func (de Design) DocContext() string {
	return de.docContext
}

func (de Design) DocAuthType() string {
	return de.docAuthType
}

func (de Design) DocSvcType() string {
	return de.docSvcType
}

func (de Design) DocSvcEncPoint() string {
	return de.docSvcEncPoint
}

func (de Design) Equal(cd Design) bool {
	if de.didMethod != cd.didMethod {
		return false
	}

	if de.docContext != cd.docContext {
		return false
	}

	if de.docAuthType != cd.docAuthType {
		return false
	}

	if de.docSvcType != cd.docSvcType {
		return false
	}

	if de.docSvcEncPoint != cd.docSvcEncPoint {
		return false
	}

	return true
}
