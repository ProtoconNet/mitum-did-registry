package digest

import (
	mongodb "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonutil "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-did-registry/state"
	"github.com/ProtoconNet/mitum-did-registry/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type DIDRegistryDesignDoc struct {
	mongodb.BaseDoc
	st     base.State
	design types.Design
}

// NewDIDRegistryDesignDoc get the State of DID Design
func NewDIDRegistryDesignDoc(st base.State, enc encoder.Encoder) (DIDRegistryDesignDoc, error) {
	design, err := state.GetDesignFromState(st)

	if err != nil {
		return DIDRegistryDesignDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DIDRegistryDesignDoc{}, err
	}

	return DIDRegistryDesignDoc{
		BaseDoc: b,
		st:      st,
		design:  design,
	}, nil
}

func (doc DIDRegistryDesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.DIDStateKeyPrefix, 3)

	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()

	return bsonutil.Marshal(m)
}

type DIDDataDoc struct {
	mongodb.BaseDoc
	st   base.State
	data types.Data
}

func NewDIDDataDoc(st base.State, enc encoder.Encoder) (DIDDataDoc, error) {
	data, err := state.GetDataFromState(st)
	if err != nil {
		return DIDDataDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DIDDataDoc{}, err
	}

	return DIDDataDoc{
		BaseDoc: b,
		st:      st,
		data:    data,
	}, nil
}

func (doc DIDDataDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.DIDStateKeyPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["publicKey"] = doc.data.PubKey()
	m["height"] = doc.st.Height()

	return bsonutil.Marshal(m)
}

type DIDDocumentDoc struct {
	mongodb.BaseDoc
	st       base.State
	document types.Document
}

func NewDIDDocumentDoc(st base.State, enc encoder.Encoder) (DIDDocumentDoc, error) {
	doc, err := state.GetDocumentFromState(st)
	if err != nil {
		return DIDDocumentDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DIDDocumentDoc{}, err
	}

	return DIDDocumentDoc{
		BaseDoc:  b,
		st:       st,
		document: doc,
	}, nil
}

func (doc DIDDocumentDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.DIDStateKeyPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["did"] = doc.document.DIDDoc().DID()
	m["height"] = doc.st.Height()

	return bsonutil.Marshal(m)
}
