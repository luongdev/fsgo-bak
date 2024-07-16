package internal

import (
	"bufio"
	"context"
	"errors"
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

	resChans   map[fsgo.ChannelType]chan fsgo.Response
	eventChans map[fsgo.ChannelType]chan fsgo.Event
	exitChan   chan bool
	errChan    chan error
}

func (c *connection) Send(cmd fsgo.Command) (fsgo.Response, error) {
	return nil, nil
}

func (c *connection) SendCtx(ctx context.Context, cmd fsgo.Command) (fsgo.Response, error) {
	if cmd == nil {
		return nil, fsgo.NewSyntaxError("command is nil")
	}
	var b []byte
	if raw, ok := cmd.(*rawCommand); ok {
		b = raw.Bytes()
	} else {
		b = []byte(fmt.Sprintf("%v%v%v", cmd.Raw(), fsgo.EOL, fsgo.EOL))
	}

	c.mu.Lock()
	if deadline, ok := ctx.Deadline(); ok {
		_ = c.conn.SetWriteDeadline(deadline)
	}

	_, err := c.conn.Write(b)
	c.mu.Unlock()

	if err != nil {
		return nil, err
	}

	c.resMu.RLock()
	defer c.resMu.RUnlock()
	select {
	case res := <-c.resChans[fsgo.TypCommandReply]:
		if res == nil {
			return nil, errors.New("connection closed")
		}
		return res, nil
	case res := <-c.resChans[fsgo.TypApiResponse]:
		if res == nil {
			return nil, errors.New("connection closed")
		}
		return res, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
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
	<-c.resChans[fsgo.TypAuthRequest]

	res, err := c.SendCtx(c.ctx, newAuthCommand(pass))
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("connection closed")
	}

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

	chann, ok := c.resChans[fsgo.ChannelType(t)]
	if !ok {
		//TODO: log error
	} else {
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

			chann <- msg
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
		resChans: map[fsgo.ChannelType]chan fsgo.Response{
			fsgo.TypCommandReply: make(chan fsgo.Response),
			fsgo.TypApiResponse:  make(chan fsgo.Response),
			fsgo.TypAuthRequest:  make(chan fsgo.Response, 1),
		},
		eventChans: map[fsgo.ChannelType]chan fsgo.Event{
			fsgo.TypEventPlain: make(chan fsgo.Event),
			fsgo.TypEventJson:  make(chan fsgo.Event),
		},
	}
	c.ctx, c.ctxCancel = context.WithCancel(opts.Context)
	go func() {
		_ = c.recvLoop()
		// TODO: log error when receive loop error
	}()
	return c, nil
}

var _ fsgo.Connection = (*connection)(nil)
