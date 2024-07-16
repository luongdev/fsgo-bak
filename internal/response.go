package internal

import (
	"bufio"
	"github.com/luongdev/fsgo"
	"strings"
)

const replyText = "Reply-Text"

type response struct {
	*message

	raw string
	err error
}

func (r *response) Error() error {
	return r.err
}

func newMessageResponse(msg *message) *response {
	if msg == nil {
		return nil
	}

	res := &response{message: msg}
	repText, ok := msg.Header(replyText)
	if ok {
		res.raw = repText
	} else {
		res.raw = string(msg.body)
	}

	res.parse()

	return res
}

func newRawResponse(r *bufio.Reader) *response {
	msg, _ := newMessage(r)
	return newMessageResponse(msg)
}

func (r *response) parse() {
	if strings.HasPrefix(r.raw, "-ERR") {
		r.err = fsgo.NewCommandError(r.raw[5:])
		return
	}

	if strings.HasPrefix(r.raw, "Job-UUID: ") {
		r.raw = r.raw[10:]
		return
	}

	if strings.HasPrefix(r.raw, "+OK") {
		r.raw = r.raw[4:]
		return
	}
}

var _ fsgo.Response = (*response)(nil)
