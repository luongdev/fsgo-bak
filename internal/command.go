package internal

type authCommand struct {
	*rawCommand
}

func (c *authCommand) Bytes() []byte {
	return []byte(c.Raw())
}

func newAuthCommand(password string) *authCommand {
	return &authCommand{rawCommand: newRawCommand("AUTH " + password)}
}
