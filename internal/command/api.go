package command

import (
	"github.com/luongdev/fsgo"
)

type api struct {
	rawCommand
	name       string
	args       string
	background bool
}

func (a *api) Bytes() []byte {

}

func (a *api) Raw() string {

}

func newApi(name, args string) fsgo.Command {
	return &api{name: name, args: args, rawCommand: newRawCommand("API " + name + " " + args)}
}

var _ fsgo.Command = (*api)(nil)
