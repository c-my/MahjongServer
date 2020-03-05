package message

const (
	GameMsgType int = iota
	ChatMsgType
)

type CommonMsg struct {
	MsgType int `json:"msg_type"`
	Content interface{}
}
