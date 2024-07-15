package internal

import (
	"bufio"
	"github.com/luongdev/fsgo"
	"github.com/luongdev/fsgo/errors"
	"io"
	"net/textproto"
	"strconv"
)

type message struct {
	headers map[string]string
	body    []byte

	r  *bufio.Reader
	tr *textproto.Reader
}

func (m *message) Header(header string) (string, bool) {

}

func (m *message) Variable(key string) (string, bool) {

}

func (m *message) parse() error {
	mHeaders, err := m.tr.ReadMIMEHeader()
	if err != nil && err.Error() != "EOF" {
		return errors.ErrCannotReadMIMEHeader
	}

	if mHeaders == nil || mHeaders.Get("Content-Type") == "" {
		return errors.ErrInvalidContentType
	}

	if lv := mHeaders.Get("Content-Length"); lv != "" {
		l, err := strconv.Atoi(lv)
		if err != nil {
			return errors.ErrInvalidContentLength
		}

		m.body = make([]byte, l)
		if _, err := io.ReadFull(m.r, m.body); err != nil {
			return err
		}
	}

	cType := mHeaders.Get("Content-Type")
}

func newMessage(r *bufio.Reader) fsgo.Message {
	return &message{}
}

var _ fsgo.Message = (*message)(nil)
