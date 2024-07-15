package fsgo

type MessageType string

const (
	TypEventPlain   MessageType = "text/event-plain"
	TypEventJson    MessageType = "text/event-json"
	TypApiResponse  MessageType = "api/response"
	TypCommandReply MessageType = "command/reply"
	TypDisconnect   MessageType = "text/disconnect-notice"
	TypeAuthRequest MessageType = "auth/request"
)

const (
	ReadBuffSize = 1024 << 6
	EOL          = "\r\n"
)
