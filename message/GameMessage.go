package message

import "github.com/c-my/MahjongServer/model"

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
	Discard //打牌
)

type GameMsgRecv struct {
	MsgType int `json:"msg_type"`

	PlayerID int        `json:"player_id"`
	Action   int        `json:"action"`
	Tile     model.Tile `json:"tile"`
}

type GameMsgSend struct {
	MsgType int `json:"msg_type"`

	CurrentTurn      int          `json:"current_turn"`
	CurrentTile      []model.Tile `json:"current_tile"`
	AvailableActions []int        `json:"available_actions"`
	LastTurn         int          `json:"last_turn"`
	LastTile         model.Tile   `json:"last_tile"`
	LastAction       int          `json:"last_action"`
	LastOpenTile     []model.Tile `json:"last_open_tile"`

	TilesCount []int `json:"tiles_count"`
}
