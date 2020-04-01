package message

type TableOrderMsg struct {
	MsgType    int `json:"msg_type"`
	TableOrder int `json:"table_order"`
}

type RoomInfoMsg struct {
	MsgType int `json:"msg_type"`
}

type GameResultMsg struct {
	MsgType int `json:"msg_type"`
	Winner  int `json:"winner"`
}
