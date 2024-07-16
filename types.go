package fsgo

import (
	"context"
	"time"
)

type Command interface {
	Bytes() []byte
	Raw() string
}

type Message interface {
	Header(header string) (string, bool)
	Variable(variable string) (string, bool)
}

type Response interface {
	Message

	Error() error
}

type Event interface {
	Message

	UID() string
	CallID() string
}

type ConnectOptions struct {
	Context      context.Context
	Timeout      time.Duration
	OnDisconnect func(string)
}

type Connection interface {
	Send(cmd Command) (Response, error)
	SendCtx(ctx context.Context, cmd Command) (Response, error)
	Close() error
}

type Client interface {
}

type Server interface {
}
