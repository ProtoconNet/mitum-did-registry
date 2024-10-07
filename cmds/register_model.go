package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-did/operation/did"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type RegisterModelCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender         currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract       currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract account to register policy" required:"true"`
	DIDMethod      string                      `arg:"" name:"did-method" help:"did method" required:"true"`
	DocContext     string                      `arg:"" name:"doc-context" help:"doc context" required:"true"`
	DocAuthType    string                      `arg:"" name:"doc-authType" help:"doc authentication type" required:"true"`
	DocSvcType     string                      `arg:"" name:"doc-serviceType" help:"doc service type" required:"true"`
	DocSvcEndPoint string                      `arg:"" name:"doc-serviceEndpoint" help:"doc service endpoint" required:"true"`
	Currency       currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender         base.Address
	contract       base.Address
}

func (cmd *RegisterModelCommand) Run(pctx context.Context) error {
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

func (cmd *RegisterModelCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Contract.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid contract format; %q", cmd.Contract)
	} else {
		cmd.contract = a
	}

	if len(cmd.DIDMethod) < 1 {
		return errors.Errorf("invalid DID Method, %s", cmd.DIDMethod)
	}

	if len(cmd.DocContext) < 1 {
		return errors.Errorf("invalid Document context, %s", cmd.DocContext)
	}

	if len(cmd.DocAuthType) < 1 {
		return errors.Errorf("invalid Document authentication type, %s", cmd.DocAuthType)
	}

	if len(cmd.DocSvcType) < 1 {
		return errors.Errorf("invalid Document service type, %s", cmd.DocSvcType)
	}

	if len(cmd.DocSvcEndPoint) < 1 {
		return errors.Errorf("invalid Document service endpoint, %s", cmd.DocSvcType)
	}

	return nil
}

func (cmd *RegisterModelCommand) createOperation() (base.Operation, error) {
	e := util.StringError("failed to create register-model operation")

	fact := did.NewRegisterModelFact(
		[]byte(cmd.Token), cmd.sender, cmd.contract, cmd.DIDMethod,
		cmd.DocContext, cmd.DocAuthType, cmd.DocSvcType, cmd.DocSvcEndPoint, cmd.Currency.CID,
	)

	op, err := did.NewRegisterModel(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
