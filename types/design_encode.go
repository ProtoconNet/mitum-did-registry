package types

import (
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *Design) unmarshal(
	_ encoder.Encoder,
	ht hint.Hint,
	didMethod, docContext, docAuthType, docSvcType, docSvcEndPoint string,
) error {
	de.BaseHinter = hint.NewBaseHinter(ht)
	de.didMethod = didMethod
	de.docContext = docContext
	de.docAuthType = docAuthType
	de.docSvcType = docSvcType
	de.docSvcEndPoint = docSvcEndPoint

	return nil
}
