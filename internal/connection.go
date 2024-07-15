package internal

import (
	"bufio"
	"context"
	"fmt"
	"github.com/luongdev/fsgo"
	"github.com/luongdev/fsgo/errors"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

var defaultOpts = fsgo.ConnectOptions{
	Context: context.Background(),
	Timeout: 10 * time.Second,
}

var endOfLine = fmt.Sprintf("%v%v", fsgo.EOF, fsgo.EOF)

type connection struct {
	conn net.Conn
	ctx  context.Context

	m sync.Mutex
}

func (c *connection) Handle() error {
	done := make(chan bool)

	rbuf := bufio.NewReaderSize(c.conn, fsgo.ReadBufferSize)

	go func() {
		for {
			//msg, err := newMessage(rbuf, true)
			//
			//if err != nil {
			//	c.err <- err
			//	done <- true
			//	break
			//}
			//
			//c.m <- msg
		}
	}()

	<-done

	return c.Close()
}

func (c *connection) Send(commands ...fsgo.Command) error {
	if len(commands) == 0 {
		return errors.ErrInvalidCommand
	}

	for _, cmd := range commands {
		if err := c.send(cmd); err != nil {
			return err
		}
	}

	return nil
}

func (c *connection) SendCtx(ctx context.Context, commands ...fsgo.Command) error {
	if len(commands) == 0 {
		return errors.ErrInvalidCommand
	}

	for _, cmd := range commands {
		if err := c.send(cmd); err != nil {
			return err
		}
	}

	return nil
}

func (c *connection) send(cmd fsgo.Command) error {
	if cmd == nil || strings.HasSuffix(cmd.Raw(), fsgo.EOF) {
		return errors.ErrInvalidCommand
	}

	c.m.Lock()
	defer c.m.Unlock()

	if _, err := c.conn.Write(cmd.Bytes()); err != nil {
		return err
	}

	if _, err := io.WriteString(c.conn, endOfLine); err != nil {
		return err
	}

	return nil
}

func (c *connection) Close() error {
	return c.conn.Close()
}

func newConnection(opts *fsgo.ConnectOptions) (fsgo.Connection, error) {
	if opts == nil {
		opts = &defaultOpts
	}

	addr := fmt.Sprintf("%v:%v", opts.Host, opts.Port)
	conn, err := net.DialTimeout("tcp", addr, opts.Timeout)
	if err != nil {
		return nil, err
	}

	c := &connection{ctx: opts.Context, conn: conn}

	return c, nil
}

var _ fsgo.Connection = (*connection)(nil)
