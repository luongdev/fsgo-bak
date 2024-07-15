package internal

import "github.com/luongdev/fsgo"

func Main() {
	conn, err := newConnection(&fsgo.ConnectOptions{
		Host: "103.141.141.53",
		Port: 65021,
	})

	if err != nil {
		panic(err)
	}

	cmd := newCommand("show status")
	err = conn.Send(cmd)

	if err != nil {
		panic(err)
	}
}
