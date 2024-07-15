package internal

import (
	"github.com/luongdev/fsgo"
)

type command struct {
	s string
}

func (c *command) Raw() string {
	return c.s
}

func (c *command) Bytes() []byte {
	return []byte(c.Raw())
}

func newCommand(s string) fsgo.Command {
	return &command{}
}

var _ fsgo.Command = (*command)(nil)
