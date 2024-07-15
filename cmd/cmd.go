package main

import "github.com/luongdev/fsgo/internal"

func main() {
	c, e := internal.NewConnection("103.141.141.53:65021", nil)
	if e != nil {
		panic(e)
	}

	internal.Loop(c)

	defer c.Close()
}
