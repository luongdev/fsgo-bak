package fsgo

type ChannelType string

const (
	TypEventPlain   ChannelType = "text/event-plain"
	TypEventJson    ChannelType = "text/event-json"
	TypApiResponse  ChannelType = "api/response"
	TypCommandReply ChannelType = "command/reply"
	TypDisconnect   ChannelType = "text/disconnect-notice"
	TypAuthRequest  ChannelType = "auth/request"
)

const (
	ReadBuffSize = 1024 << 6
	EOL          = "\r\n"
)
