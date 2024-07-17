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

	ctx          context.Context
	ctxCancel    context.CancelFunc
	onDisconnect func(string)

	mu    sync.Mutex
	resMu sync.RWMutex

	authChan   chan fsgo.Message
	resChans   map[fsgo.ChannelType]chan fsgo.Response
	eventChans map[fsgo.ChannelType]chan fsgo.Event
	exitChan   chan bool
	errChan    chan error
	closed     bool
}

func (c *connection) Send(cmd fsgo.Command) (fsgo.Response, error) {
	return c.SendCtx(c.ctx, cmd)
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
			select {
			case <-c.exitChan:
				return nil
			case c.errChan <- err:
			}
		}
	}
}

func (c *connection) Close() error {
	if c.closed {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.closed = true
	c.ctxCancel()

	close(c.exitChan)
	close(c.authChan)
	close(c.errChan)

	for _, chann := range c.resChans {
		close(chann)
	}
	for _, chann := range c.eventChans {
		close(chann)
	}

	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *connection) Auth(pass string) error {
	<-c.authChan
	res, err := c.SendCtx(c.ctx, newAuthCommand(pass))
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("connection closed")
	}

	return nil
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
	cType := fsgo.ChannelType(t)

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

		switch cType {
		case fsgo.TypCommandReply, fsgo.TypApiResponse:
			resChan, ok := c.resChans[cType]
			if ok {
				resChan <- newMessageResponse(msg)
			} else {
				c.errChan <- fmt.Errorf("no response channel for %s", cType)
			}
		case fsgo.TypAuthRequest:
			c.authChan <- msg
		case fsgo.TypEventPlain:
			eventChan, ok := c.eventChans[cType]
			if ok {
				eventChan <- newMessageEvent(msg)
			} else {
				c.errChan <- fmt.Errorf("no event channel for %s", cType)
			}
		case fsgo.TypDisconnect:
			c.exitChan <- true
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
		authChan: make(chan fsgo.Message, 1),
		resChans: map[fsgo.ChannelType]chan fsgo.Response{
			fsgo.TypCommandReply: make(chan fsgo.Response),
			fsgo.TypApiResponse:  make(chan fsgo.Response),
		},
		eventChans: map[fsgo.ChannelType]chan fsgo.Event{
			fsgo.TypEventPlain: make(chan fsgo.Event),
		},
	}
	c.ctx, c.ctxCancel = context.WithCancel(opts.Context)
	go func() {
		_ = c.recvLoop()
		// TODO: log error when receive loop error
	}()
	go func() {
		<-c.exitChan
		if c.onDisconnect != nil {
			c.onDisconnect(addr)
		}
		_ = c.Close()
	}()

	return c, nil
}

var _ fsgo.Connection = (*connection)(nil)
