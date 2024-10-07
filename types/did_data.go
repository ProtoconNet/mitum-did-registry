package types

import (
	"encoding/hex"
	"fmt"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"golang.org/x/crypto/sha3"
	"strings"
)

const MinKeyLen = 128
const Nid = "0000"

var DataHint = hint.MustNewHint("mitum-did-data-v0.0.1")
var DocumentHint = hint.MustNewHint("mitum-did-document-v0.0.1")

type Data struct {
	hint.BaseHinter
	pubKey string
	did    string
}

func NewData(
	pkey, method string,
) Data {
	data := Data{
		BaseHinter: hint.NewBaseHinter(DataHint),
	}
	// Detach 0x
	pubKey := strings.TrimPrefix(pkey, "0x")
	// reform pubkey
	pubKey = "04" + pubKey[len(pubKey)-128:]
	data.pubKey = pubKey

	digest1 := sha3.Sum256([]byte(pubKey))
	idString := hex.EncodeToString(digest1[:])[0:40]
	digest2 := sha3.Sum256([]byte(Nid + idString))
	checksum := hex.EncodeToString(digest2[:])[0:8]
	specificIdString := Nid + idString + checksum
	data.did = fmt.Sprintf("did:%s:%s", method, specificIdString)
	return data
}

func (d Data) IsValid([]byte) error {
	return nil
}

func (d Data) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(d.did),
	)
}

func (d Data) PubKey() string {
	return d.pubKey
}

func (d Data) DID() string {
	return d.did
}

func (d Data) Equal(ct Data) bool {
	if d.did != ct.did {
		return false
	}

	return true
}
