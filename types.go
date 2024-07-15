package fsgo

import (
	"context"
	"time"
)

type Command interface {
	Raw() string
	Bytes() []byte
}

type ConnectOptions struct {
	Host    string
	Port    uint16
	Timeout time.Duration
	Context context.Context
}

type InboundOptions struct {
	*ConnectOptions

	Password string
}

type OutboundOptions struct {
	*ConnectOptions
}

type Connection interface {
	Send(commands ...Command) error
	SendCtx(ctx context.Context, commands ...Command) error

	Close() error
}

type Message interface {
	Header(header string) (string, bool)
	Variable(key string) (string, bool)
}
