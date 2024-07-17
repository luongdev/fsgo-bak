package internal

import "github.com/luongdev/fsgo"

type sendMessage struct {
}

func (s *sendMessage) Bytes() []byte {
	//TODO implement me
	panic("implement me")
}

func (s *sendMessage) Raw() string {
	//TODO implement me
	panic("implement me")
}

var _ fsgo.Command = (*sendMessage)(nil)
