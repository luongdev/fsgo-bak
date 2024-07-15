package internal

import (
	"bufio"
	"fmt"
	"github.com/luongdev/fsgo"
	"io"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

type message struct {
	headers map[string]string
	body    []byte

	r  *bufio.Reader
	tr *textproto.Reader
}

func (m *message) Header(header string) (string, bool) {
	return m.headers[header], true
}

func (m *message) Variable(variable string) (string, bool) {
	return m.Header(fmt.Sprintf("variable_%s", variable))
}

func (m *message) String() string {
	return fmt.Sprintf("%v\nbody=%s", m.headers, m.body)
}

func (m *message) parse() error {
	mh, err := m.tr.ReadMIMEHeader()
	if err != nil && err.Error() != "EOF" {
		return err
	}

	if mh == nil {
		return fsgo.NewRuntimeError("could not read MIME headers")
	}

	if mh.Get("Content-Type") == "" {
		return fsgo.NewSyntaxError("content-type missing")
	}

	if lstr := mh.Get("Content-Length"); lstr != "" {
		_len, err := strconv.Atoi(lstr)
		if err != nil {
			return fsgo.NewSyntaxError("content-type missing")
		}
		m.body = make([]byte, _len)
		if _, err := io.ReadFull(m.r, m.body); err != nil {
			return fsgo.NewRuntimeError("could not read body")
		}
	}

	for k, v := range mh {
		m.headers[k] = v[0]
		if strings.Contains(v[0], "%") {
			m.headers[k], err = url.QueryUnescape(v[0])
			if err != nil {
				//TODO: log error
				continue
			}
		}
	}

	return nil
}

func newMessage(r *bufio.Reader) (*message, error) {
	msg := &message{r: r, tr: textproto.NewReader(r), headers: make(map[string]string)}
	if err := msg.parse(); err != nil {
		return nil, err
	}

	return msg, nil
}

var _ fsgo.Message = (*message)(nil)
