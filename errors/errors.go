package errors

type CommandError struct {
	msg string
}

func (e *CommandError) Error() string {
	return e.msg
}

func NewCommandError(msg string) error {
	return &CommandError{msg: msg}
}

var (
	ErrInvalidCommand       = NewCommandError(InvalidCommand)
	ErrCannotReadMIMEHeader = NewCommandError(MIMEHeaderError)
	ErrInvalidContentType   = NewCommandError(InvalidContentType)
	ErrInvalidContentLength = NewCommandError(InvalidContentType)
)
