package internal

import (
	"github.com/luongdev/fsgo"
	"github.com/luongdev/fsgo/internal/command"
)

func NewConnection(addr string, opts *fsgo.ConnectOptions) (fsgo.Connection, error) {
	return newConnection(addr, opts)
}

func Loop(con fsgo.Connection) {
	conn, ok := con.(*connection)
	if !ok {
		return
	}

	err := conn.Auth("Simplefs!!")
	if err != nil {
		panic(err)
	}
}

type RawCommand = command.RawCommand

func NewRawCommand(raw string) *RawCommand {
	return command.NewRawCommand(raw)
}
