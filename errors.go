package fsgo

func NewRuntimeError(m string) error {
	return &RuntimeError{m: m}
}

type RuntimeError struct {
	m string
}

func (e *RuntimeError) Error() string {
	return e.m
}

func NewSyntaxError(m string) error {
	return &SyntaxError{m: m}
}

type SyntaxError struct {
	m string
}

func (e *SyntaxError) Error() string {
	return e.m
}
