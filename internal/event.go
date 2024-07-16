package internal

import (
	"fmt"
	"github.com/luongdev/fsgo"
)

const (
	uniqueID     = "Unique-ID"
	callUniqueID = "Caller-Unique-ID"
)

type event struct {
	*message
}

func (m *event) UID() string {
	u, ok := m.Header(uniqueID)
	if !ok {
		return ""
	}

	return u
}

func (m *event) CallID() string {
	u, ok := m.Header(callUniqueID)
	if !ok {
		return ""
	}

	return u
}

func (m *event) String() string {
	return fmt.Sprintf("%v body=%s", m.headers, m.body)
}

func newMessageEvent(msg *message) *event {
	if msg == nil {
		return nil
	}

	e := &event{message: msg}

	return e
}

var _ fsgo.Event = (*event)(nil)
