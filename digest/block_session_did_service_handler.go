package digest

import (
	"github.com/ProtoconNet/mitum-did/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareDID() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var didModels []mongo.WriteModel
	var didDataModels []mongo.WriteModel
	var didDocumentModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case state.IsDesignStateKey(st.Key()):
			j, err := bs.handleDIDDesignState(st)
			if err != nil {
				return err
			}
			didModels = append(didModels, j...)
		case state.IsDataStateKey(st.Key()):
			j, err := bs.handleDIDDataState(st)
			if err != nil {
				return err
			}
			didDataModels = append(didDataModels, j...)
		case state.IsDocumentStateKey(st.Key()):
			j, err := bs.handleDIDDocumentState(st)
			if err != nil {
				return err
			}
			didDocumentModels = append(didDocumentModels, j...)
		default:
			continue
		}
	}

	bs.didModels = didModels
	bs.didDataModels = didDataModels
	bs.didDocumentModels = didDocumentModels

	return nil
}

func (bs *BlockSession) handleDIDDesignState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if DIDDesignDoc, err := NewDIDDesignDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(DIDDesignDoc),
		}, nil
	}
}

func (bs *BlockSession) handleDIDDataState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if DIDDataDoc, err := NewDIDDataDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(DIDDataDoc),
		}, nil
	}
}

func (bs *BlockSession) handleDIDDocumentState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if DIDDocumentDoc, err := NewDIDDocumentDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(DIDDocumentDoc),
		}, nil
	}
}
