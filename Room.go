package main

import (
	"github.com/c-my/MahjongServer/game"
	"github.com/c-my/MahjongServer/model"
	"github.com/gorilla/websocket"
)

type Room struct {
	connManager    *ConnManager
	mahjongManager game.Manager

	gameRecvCh chan model.GameMsgRecv
	gameSendCh chan model.GameMsgSend
}

func NewRoom() *Room {
	gameRecvChannel := make(chan model.GameMsgRecv)
	gameSendChannel := make(chan model.GameMsgSend)
	return &Room{
		gameRecvCh:     gameRecvChannel,
		gameSendCh:     gameSendChannel,
		connManager:    NewConnManager(4, gameRecvChannel, gameSendChannel),
		mahjongManager: game.NewJinzhouMahjong(gameRecvChannel, gameSendChannel),
	}
}

func (r *Room) AddConn(conn *websocket.Conn) {
	r.connManager.SetConn(conn)
}

func (r *Room) Start() {
	r.mahjongManager.Start()
}
