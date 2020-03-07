package message

import "github.com/c-my/MahjongServer/model"

//type GameMsgRecv struct {
//	Action string `json:"action"`
//}

// Player actions
const (
	Start         int = iota
	Deal              //发牌
	Chow              //吃
	Pong              //碰
	ExposedKong       //杠
	ConcealedKong     //暗杠
	AddedKong         //补杠
	Win               //胡
	Cancel            //取消操作
	Discard           //打牌
)

// 吃牌的类型
const (
	LeftChow  int = iota //作为左边牌吃
	MidChow              //作为中间牌吃
	RightChow            //作为右边牌吃
	NAC                  //不是吃
)

type GameMsgRecv struct {
	MsgType int `json:"msg_type"`

	TableOrder int        `json:"table_order"`
	Action     int        `json:"action"`
	Tile       model.Tile `json:"tile"`
	ChowType   int        `json:"chow_type"`
}

type GameMsgSend struct {
	MsgType int `json:"msg_type"`

	TableOrder       int          `json:"table_order"`       // 发送给谁，以后可能会用到
	CurrentTurn      int          `json:"current_turn"`      // 当前轮到的玩家
	CurrentTile      model.Tile `json:"current_tile"`      // 玩家收到的牌
	AvailableActions []int        `json:"available_actions"` //玩家可以进行的动作，打牌、吃、碰等
	LastTurn         int          `json:"last_turn"`         // 上一个玩家
	LastAction       int          `json:"last_action"`       // 上一个玩家的动作

	PlayerTile []model.PlayerTile `json:"player_tile"` //全局麻将牌牌信息
	WallCount  int                `json:"wall_count"`  // 剩余牌墙数量
}
