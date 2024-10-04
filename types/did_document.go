package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type Document struct {
	hint.BaseHinter
	didDoc DIDDocument
}

func NewDocument(
	didDoc DIDDocument,
) Document {
	return Document{
		BaseHinter: hint.NewBaseHinter(DataHint),
		didDoc:     didDoc,
	}
}

func (d Document) IsValid([]byte) error {
	return nil
}

func (d Document) Bytes() []byte {
	return util.ConcatBytesSlice(
		d.didDoc.Bytes(),
	)
}

func (d Document) DIDDoc() DIDDocument {
	return d.didDoc
}

type DIDDocument struct {
	context_       string
	id             string
	created        string
	status         string
	authentication Authentication
	service        Service
}

func NewDIDDocument(
	context_, did, created, status, authType, publicKeyHex, serviceType, serviceEndPoint string,
) DIDDocument {
	return DIDDocument{
		context_: context_,
		id:       did,
		created:  created,
		status:   status,
		authentication: Authentication{
			id:           did,
			authType:     authType,
			controller:   did,
			publicKeyHex: publicKeyHex,
		},
		service: Service{
			id:              did,
			serviceType:     serviceType,
			serviceEndPoint: serviceEndPoint,
		},
	}
}

func (d DIDDocument) IsValid([]byte) error {
	return nil
}

func (d DIDDocument) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(d.context_),
		[]byte(d.id),
		[]byte(d.created),
		[]byte(d.status),
		d.authentication.Bytes(),
		d.service.Bytes(),
	)
}

func (d DIDDocument) DID() string {
	return d.id
}

func (d DIDDocument) Status() string {
	return d.status
}

func (d *DIDDocument) SetStatus(status string) {
	d.status = status
}

type Authentication struct {
	id           string
	authType     string
	controller   string
	publicKeyHex string
}

func NewAuthentication(
	id, authType, controller, publicKeyHex string,
) Authentication {
	return Authentication{
		id:           id,
		authType:     authType,
		controller:   controller,
		publicKeyHex: publicKeyHex,
	}
}

func (d Authentication) IsValid([]byte) error {
	return nil
}

func (d Authentication) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(d.id),
		[]byte(d.authType),
		[]byte(d.controller),
		[]byte(d.publicKeyHex),
	)
}

type Service struct {
	id              string
	serviceType     string
	serviceEndPoint string
}

func NewService(
	id, serviceType, serviceEndPoint string,
) Service {
	return Service{
		id:              id,
		serviceType:     serviceType,
		serviceEndPoint: serviceEndPoint,
	}
}

func (d Service) IsValid([]byte) error {
	return nil
}

func (d Service) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(d.id),
		[]byte(d.serviceType),
		[]byte(d.serviceEndPoint),
	)
}
