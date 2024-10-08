package cmds

import (
	"context"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-did/operation/did"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type MigrateDIDCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	PubKey   string                      `arg:"" name:"pubKey" help:"public key" required:"true"`
	TxID     string                      `arg:"" name:"txID" help:"txID" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
}

func (cmd *MigrateDIDCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *MigrateDIDCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	a, err = cmd.Contract.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid contract format, %q", cmd.Contract)
	} else {
		cmd.contract = a
	}

	if len(cmd.PubKey) < 1 {
		return errors.Errorf("invalid Public Key, %s", cmd.PubKey)
	}

	if len(cmd.TxID) < 1 {
		return errors.Errorf("invalid TxID, %s", cmd.TxID)
	}

	return nil
}

func (cmd *MigrateDIDCommand) createOperation() (base.Operation, error) { // nolint:dupl
	e := util.StringError("failed to create MigrateDID operation")

	item := did.NewMigrateDIDItem(cmd.contract, cmd.PubKey, cmd.TxID, cmd.Currency.CID)

	fact := did.NewMigrateDIDFact([]byte(cmd.Token), cmd.sender, []did.MigrateDIDItem{item})

	op, err := did.NewMigrateDID(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
