package internal

import (
	"bufio"
	"context"
	"fmt"
	"github.com/luongdev/fsgo"
	"net"
	"sync"
	"time"
)

var defaultOpts = &fsgo.ConnectOptions{Timeout: 10 * time.Second, Context: context.Background()}

type connection struct {
	conn net.Conn

	ctx       context.Context
	ctxCancel context.CancelFunc
	mu        sync.Mutex
	resMu     sync.RWMutex

	errChan  chan error
	exitChan chan bool
	authChan chan fsgo.Response
}

func (c *connection) Send(cmd fsgo.Command) (fsgo.Response, error) {
	return nil, nil
}

func (c *connection) SendCtx(ctx context.Context, cmd fsgo.Command) (fsgo.Response, error) {
	return nil, nil
}

func (c *connection) recvLoop() error {
	buff := bufio.NewReaderSize(c.conn, fsgo.ReadBuffSize)
	for {
		err := c.read(buff)
		if err != nil {
			c.errChan <- err
		}
	}
}

func (c *connection) Close() error {
	c.ctxCancel()
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *connection) Auth(pass string) error {
	<-c.authChan

	return fmt.Errorf("not implemented")
}

func (c *connection) read(r *bufio.Reader) error {
	msg, err := newMessage(r)
	if err != nil {
		return err
	}

	t, ok := msg.Header("Content-Type")
	if !ok {
		return fsgo.NewSyntaxError("content-type missing")
	}

	ctx, cancel := context.WithTimeout(c.ctx, time.Second*5)
	defer cancel()

	select {
	case <-ctx.Done():
		//TODO: handle when message handle timeout
	case <-c.ctx.Done():
		if c.ctx.Err() != nil {
			return c.ctx.Err()
		}
		return nil
	default:
		c.resMu.RLock()
		defer c.resMu.RUnlock()
		switch fsgo.MessageType(t) {
		case fsgo.TypeAuthRequest:
			c.authChan <- msg
		}
	}

	return nil
}

func newConnection(addr string, opts *fsgo.ConnectOptions) (*connection, error) {
	if opts == nil {
		opts = defaultOpts
	} else {
		if opts.Context == nil {
			opts.Context = defaultOpts.Context
		}

		if opts.Timeout <= 0 {
			opts.Timeout = defaultOpts.Timeout
		}
	}

	netConn, err := net.DialTimeout("tcp", addr, opts.Timeout)
	if err != nil {
		return nil, err
	}

	c := &connection{
		conn:     netConn,
		exitChan: make(chan bool),
		errChan:  make(chan error),
		authChan: make(chan fsgo.Response),
	}
	c.ctx, c.ctxCancel = context.WithCancel(opts.Context)
	go func() {
		_ = c.recvLoop()
		// TODO: log error when receive loop error
	}()
	return c, nil
}

var _ fsgo.Connection = (*connection)(nil)
