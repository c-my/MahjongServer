package container

import (
	"github.com/c-my/MahjongServer/game"
	"github.com/c-my/MahjongServer/message"
	"github.com/c-my/MahjongServer/rule"
	"github.com/gorilla/websocket"
)

type Room struct {
	connManager    *ConnManager
	mahjongManager *game.MahjongManager

	gameRecvCh chan message.GameMsgRecv
	gameSendCh chan message.GameMsgSend

	password    string
	playerCount int
}

func NewRoom(rule rule.MahjongRule) *Room {
	gameRecvChannel := make(chan message.GameMsgRecv)
	gameSendChannel := make(chan message.GameMsgSend)
	return &Room{
		gameRecvCh:     gameRecvChannel,
		gameSendCh:     gameSendChannel,
		connManager:    NewConnManager(4, gameRecvChannel, gameSendChannel),
		mahjongManager: game.NewMahjongManager(gameRecvChannel, gameSendChannel, rule),
	}
}

func (r *Room) AddConn(conn *websocket.Conn) {
	r.connManager.SetConn(conn)
}

func (r *Room) Start() {
	r.mahjongManager.Start()
}
