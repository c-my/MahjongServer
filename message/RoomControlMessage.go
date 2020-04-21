package message

import "github.com/c-my/MahjongServer/model"

type TableOrderMsg struct {
	MsgType    int `json:"msg_type"`
	TableOrder int `json:"table_order"`
}

type UserOrderMsg struct {
	MsgType    int `json:"msg_type"`
	TableOrder int `json:"user_order"`
}

type RoomInfoMsg struct {
	MsgType int `json:"msg_type"`
}

type GameResultMsg struct {
	MsgType   int        `json:"msg_type"`
	Winner    int        `json:"winner"`
	FinalTile model.Tile `json:"final_tile"`
}

type GetReadyMsg struct {
	MsgType   int     `json:"msg_type"`
	ReadyList [4]bool `json:"ready_list"`
}
