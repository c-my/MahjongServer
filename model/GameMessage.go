package model

//type GameMsgRecv struct {
//	Action string `json:"action"`
//}

// Player actions
const (
	Start int = iota
	Deal      //发牌
	Chow      //吃
	Pong      //碰
	Kong      //杠
	Win
)

type GameMsgRecv struct {
	MsgType int `json:"msg_type"`

	PlayerID int  `json:"player_id"`
	Action   int  `json:"action"`
	Tile     Tile `json:"tile"`
}

type GameMsgSend struct {
	MsgType int `json:"msg_type"`

	CurrentTurn      int    `json:"current_turn"`
	CurrentTile      []Tile `json:"current_tile"`
	AvailableActions []int  `json:"available_actions"`
	LastTurn         int    `json:"last_turn"`
	LastTile         Tile   `json:"last_tile"`
	LastAction       int    `json:"last_action"`
	LastOpenTile     []Tile `json:"last_open_tile"`

	TilesCount []int `json:"tiles_count"`
}
