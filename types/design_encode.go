package types

import (
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *Design) unmarshal(
	_ encoder.Encoder,
	ht hint.Hint,
	didMethod, docContext, docAuthType, docSvcType, docSvcEncPoint string,
) error {
	de.BaseHinter = hint.NewBaseHinter(ht)
	de.didMethod = didMethod
	de.docContext = docContext
	de.docAuthType = docAuthType
	de.docSvcType = docSvcType
	de.docSvcEncPoint = docSvcEncPoint

	return nil
}
