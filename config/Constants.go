package config

// Player actions
const (
	Start         int = iota
	Deal              //发牌（开局）
	Chow              //吃
	Pong              //碰
	ExposedKong       //杠
	ConcealedKong     //暗杠
	AddedKong         //补杠
	Win               //胡
	Cancel            //取消操作
	Discard           //打牌
	Draw              //发牌（一张）
)

// 吃牌的类型
const (
	LeftChow  int = iota //作为左边牌吃
	MidChow              //作为中间牌吃
	RightChow            //作为右边牌吃
	NAC                  //不是吃
)

const (
	GameMsgType int = iota
	ChatMsgType
)
