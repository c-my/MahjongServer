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

	gameRecvCh   chan message.GameMsgRecv
	gameSendCh   chan message.GameMsgSend
	tableOrderCh chan int
	gameResultCh chan message.GameResultMsg

	password    string
	playerCount int
}

func NewRoom(rule rule.MahjongRule) *Room {
	gameRecvChannel := make(chan message.GameMsgRecv)
	gameSendChannel := make(chan message.GameMsgSend)
	tableOrderChannel := make(chan int)
	gameResultChannel := make(chan message.GameResultMsg)
	return &Room{
		gameRecvCh:     gameRecvChannel,
		gameSendCh:     gameSendChannel,
		tableOrderCh:   tableOrderChannel,
		gameResultCh:   gameResultChannel,
		connManager:    NewConnManager(4, gameRecvChannel, gameSendChannel, tableOrderChannel, gameResultChannel),
		mahjongManager: game.NewMahjongManager(gameRecvChannel, gameSendChannel, tableOrderChannel, gameResultChannel, rule),
	}
}

func (r *Room) AddConn(conn *websocket.Conn) {
	r.connManager.SetConn(conn)
}

func (r *Room) Start() {
	r.mahjongManager.Start()
}
