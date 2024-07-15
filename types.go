package fsgo

import (
	"context"
	"time"
)

type Command interface {
}

type Message interface {
	Header(header string) (string, bool)
	Variable(variable string) (string, bool)
}

type Response interface {
	Message
}

type ConnectOptions struct {
	Context context.Context
	Timeout time.Duration
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
