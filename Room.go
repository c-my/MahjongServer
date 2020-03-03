package main

import (
	"github.com/c-my/MahjongServer/game"
	"github.com/c-my/MahjongServer/model"
	"github.com/gorilla/websocket"
)

type Room struct {
	connManager    *ConnManager
	mahjongManager game.Manager

	gameCh chan model.GameMessage
}

func NewRoom() *Room {
	gameChannel := make(chan model.GameMessage)
	return &Room{
		gameCh:         gameChannel,
		connManager:    NewConnManager(4, gameChannel),
		mahjongManager: game.NewJinzhouMahjong(gameChannel),
	}
}

func (r *Room) AddConn(conn *websocket.Conn) {
	r.connManager.SetConn(conn)
}

func (r *Room) Start() {
	r.mahjongManager.Start()
}

