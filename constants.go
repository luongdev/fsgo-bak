package fsgo

const (
	TypAuth         = "auth/request"
	TypDisconnect   = "text/disconnect-notice"
	TypEventJson    = "text/event-json"
	TypEventPlain   = "text/event-plain"
	TypApiResponse  = "api/response"
	TypCommandReply = "command/reply"
)

const (
	EOF            = "\r\n"
	ReadBufferSize = 1024 << 6
)
