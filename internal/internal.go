package internal

import "github.com/luongdev/fsgo"

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

type RawCommand = rawCommand

func NewRawCommand(raw string) *RawCommand {
	return newRawCommand(raw)
}
