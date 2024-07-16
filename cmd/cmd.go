package main

import "github.com/luongdev/fsgo/internal"

type AuthCommand struct {
	*internal.RawCommand
}

func NewAuthCommand(password string) *AuthCommand {
	return &AuthCommand{
		RawCommand: internal.NewRawCommand("AUTH " + password),
	}
}

func main() {
	c, e := internal.NewConnection("103.141.141.53:65021", nil)
	if e != nil {
		panic(e)
	}

	internal.Loop(c)

	defer c.Close()
}
