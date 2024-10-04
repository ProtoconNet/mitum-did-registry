package types

import (
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (d *Data) unmarshal(
	ht hint.Hint,
	pubKey, did string,
) error {
	d.BaseHinter = hint.NewBaseHinter(ht)
	d.pubKey = pubKey
	d.did = did

	return nil
}
