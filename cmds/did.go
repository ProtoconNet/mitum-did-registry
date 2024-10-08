package cmds

type DIDCommand struct {
	CreateDID     CreateDIDCommand     `cmd:"" name:"create-did" help:"create new did"`
	MigrateDID    MigrateDIDCommand    `cmd:"" name:"migrate-did" help:"migrate did"`
	ReactivateDID ReactivateDIDCommand `cmd:"" name:"reactivate-did" help:"reactivate did"`
	DeactivateDID DeactivateDIDCommand `cmd:"" name:"deactivate-did" help:"deactivate did"`
	RegisterModel RegisterModelCommand `cmd:"" name:"register-model" help:"register did model"`
}
