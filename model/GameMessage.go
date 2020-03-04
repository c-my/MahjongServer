package model

//type GameMsgRecv struct {
//	Action string `json:"action"`
//}

const (
	Start int = iota
	Chow
	Pong
	Kong
	Win
)

type GameMsgRecv struct {
	MsgType int `json:"msg_type"`

	PlayerID int  `json:"player_id"`
	Action   int  `json:"action"`
	Tile     Tile `json:"tile"`
}

type GameMsgSend struct {
	CurrentTurn      int    `json:"current_turn"`
	CurrentTile      []Tile `json:"current_tile"`
	AvailableActions []int  `json:"available_actions"`
	LastTurn         int    `json:"last_turn"`
	LastTile         Tile   `json:"last_tile"`
	LastAction       int    `json:"last_action"`
	LastOpenTile     []Tile `json:"last_open_tile"`
}
