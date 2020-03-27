package message

type CommonMsg struct {
	MsgType int `json:"msg_type"`
	Content interface{}
}
