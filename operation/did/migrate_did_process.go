package did

import (
	"context"
	"strings"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-did/state"
	"github.com/ProtoconNet/mitum-did/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var migrateDIDItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MigrateDIDItemProcessor)
	},
}

var migrateDIDProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MigrateDIDProcessor)
	},
}

func (MigrateDID) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type MigrateDIDItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   MigrateDIDItem
}

func (ipp *MigrateDIDItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess MigrateDIDItemProcessor")
	it := ipp.item

	if err := it.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.DesignStateKey(it.Currency()), getStateFunc); err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id %v", it.Currency())))
	}

	_, cSt, aErr, cErr := currencystate.ExistsCAccount(it.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(cErr)
	}

	_, err := extensioncurrency.CheckCAAuthFromState(cSt, ipp.sender)
	if err != nil {
		return e.Wrap(err)
	}

	if st, err := currencystate.ExistsState(state.DesignStateKey(it.Contract()), "design", getStateFunc); err != nil {
		return e.Wrap(
			common.ErrServiceNF.Errorf("service design state for contract account %v", it.Contract()))
	} else if _, err := state.GetDesignFromState(st); err != nil {
		return e.Wrap(
			common.ErrServiceNF.Errorf("service design value for contract account %v", it.Contract()))
	}

	if found, _ := currencystate.CheckNotExistsState(state.DataStateKey(it.Contract(), it.PubKeyReformed()), getStateFunc); found {
		return e.Wrap(
			common.ErrStateE.Errorf(
				"pubKey %v in contract account %v",
				it.PubKeyReformed(), it.Contract()),
		)
	}

	return nil
}

func (ipp *MigrateDIDItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	var sts []base.StateMergeValue

	st, _ := currencystate.ExistsState(state.DesignStateKey(it.Contract()), "did design", getStateFunc)
	design, _ := state.GetDesignFromState(st)

	didData := types.NewData(
		it.PubKey(), design.DIDMethod(),
	)
	if err := didData.IsValid(nil); err != nil {
		return nil, err
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		state.DataStateKey(it.Contract(), it.PubKeyReformed()),
		state.NewDataStateValue(didData),
	))

	txHash := strings.TrimPrefix(it.TxID(), "0x")
	didDocument := types.NewDIDDocument(design.DocContext(), didData.DID(), txHash, "1",
		design.DocAuthType(), it.PubKeyReformed(), design.DocSvcType(), design.DocSvcEndPoint())
	document := types.NewDocument(didDocument)
	if err := document.IsValid(nil); err != nil {
		return nil, err
	}
	sts = append(sts, currencystate.NewStateMergeValue(
		state.DocumentStateKey(it.Contract(), didData.DID()),
		state.NewDocumentStateValue(document),
	))

	return sts, nil
}

func (ipp *MigrateDIDItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = MigrateDIDItem{}

	migrateDIDItemProcessorPool.Put(ipp)
}

type MigrateDIDProcessor struct {
	*base.BaseOperationProcessor
}

func NewMigrateDIDProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new MigrateDIDProcessor")

		nopp := migrateDIDProcessorPool.Get()
		opp, ok := nopp.(*MigrateDIDProcessor)
		if !ok {
			return nil, e.Errorf("expected %T, not %T", MigrateDIDProcessor{}, nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *MigrateDIDProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(MigrateDIDFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", MigrateDIDFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	_, _, aErr, cErr := currencystate.ExistsCAccount(
		fact.Sender(), "sender", true, false, getStateFunc)
	if aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: sender %v is contract account", cErr, fact.Sender())), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	for _, it := range fact.Items() {
		ip := migrateDIDItemProcessorPool.Get()
		ipc, ok := ip.(*MigrateDIDItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMTypeMismatch.Errorf("expected %T, not %T", MigrateDIDItemProcessor{}, ip)), nil
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Errorf("%v", err),
			), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *MigrateDIDProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process MigrateDID")

	fact, _ := op.Fact().(MigrateDIDFact)

	var sts []base.StateMergeValue // nolint:prealloc
	for _, it := range fact.Items() {
		ip := migrateDIDItemProcessorPool.Get()
		ipc, _ := ip.(*MigrateDIDItemProcessor)

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it

		st, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process MigrateDIDItem; %w", err), nil
		}

		sts = append(sts, st...)

		ipc.Close()
	}

	items := make([]DIDItem, len(fact.Items()))
	for i := range fact.Items() {
		items[i] = fact.Items()[i]
	}

	feeReceiverBalSts, required, err := calculateDIDItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee; %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
	}

	for cid := range sb {
		v, ok := sb[cid].Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[cid].Value())
		}

		_, feeReceiverFound := feeReceiverBalSts[cid]

		if feeReceiverFound && (sb[cid].Key() != feeReceiverBalSts[cid].Key()) {
			stmv := common.NewBaseStateMergeValue(
				sb[cid].Key(),
				statecurrency.NewDeductBalanceStateValue(v.Amount.WithBig(required[cid][1])),
				func(height base.Height, st base.State) base.StateValueMerger {
					return statecurrency.NewBalanceStateValueMerger(height, sb[cid].Key(), cid, st)
				},
			)

			r, ok := feeReceiverBalSts[cid].Value().(statecurrency.BalanceStateValue)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", statecurrency.BalanceStateValue{}, feeReceiverBalSts[cid].Value()), nil
			}
			sts = append(
				sts,
				common.NewBaseStateMergeValue(
					feeReceiverBalSts[cid].Key(),
					statecurrency.NewAddBalanceStateValue(r.Amount.WithBig(required[cid][1])),
					func(height base.Height, st base.State) base.StateValueMerger {
						return statecurrency.NewBalanceStateValueMerger(height, feeReceiverBalSts[cid].Key(), cid, st)
					},
				),
			)

			sts = append(sts, stmv)
		}
	}

	return sts, nil, nil
}

func (opp *MigrateDIDProcessor) Close() error {
	migrateDIDProcessorPool.Put(opp)

	return nil
}

func calculateDIDItemsFee(getStateFunc base.GetStateFunc, items []DIDItem) (
	map[currencytypes.CurrencyID]base.State, map[currencytypes.CurrencyID][2]common.Big, error) {
	feeReceiveSts := map[currencytypes.CurrencyID]base.State{}
	required := map[currencytypes.CurrencyID][2]common.Big{}

	for _, item := range items {
		rq := [2]common.Big{common.ZeroBig, common.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := currencystate.ExistsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, nil, err
		}

		switch k, err := policy.Feeer().Fee(common.ZeroBig); {
		case err != nil:
			return nil, nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]common.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]common.Big{rq[0].Add(k), rq[1].Add(k)}
		}

		if policy.Feeer().Receiver() == nil {
			continue
		}

		if err := currencystate.CheckExistsState(statecurrency.AccountStateKey(policy.Feeer().Receiver()), getStateFunc); err != nil {
			return nil, nil, err
		} else if st, found, err := getStateFunc(statecurrency.BalanceStateKey(policy.Feeer().Receiver(), item.Currency())); err != nil {
			return nil, nil, err
		} else if !found {
			return nil, nil, errors.Errorf("feeer receiver account not found, %s", policy.Feeer().Receiver())
		} else {
			feeReceiveSts[item.Currency()] = st
		}

	}

	return feeReceiveSts, required, nil

}
