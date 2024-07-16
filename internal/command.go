package internal

import (
	"fmt"
	"github.com/luongdev/fsgo"
)

type rawCommand struct {
	raw string
}

func (c *rawCommand) Bytes() []byte {
	return []byte(c.Raw())
}

func (c *rawCommand) Raw() string {
	return fmt.Sprintf("%v%v%v", c.raw, fsgo.EOL, fsgo.EOL)
}

func newRawCommand(raw string) *rawCommand {
	return &rawCommand{raw: raw}
}

var _ fsgo.Command = (*rawCommand)(nil)

type authCommand struct {
	*rawCommand
}

func (c *authCommand) Bytes() []byte {
	return []byte(c.Raw())
}

func newAuthCommand(password string) *authCommand {
	return &authCommand{rawCommand: newRawCommand("AUTH " + password)}
}
