package container

import (
	"github.com/c-my/MahjongServer/game"
	"github.com/c-my/MahjongServer/message"
	"github.com/c-my/MahjongServer/rule"
	"github.com/gorilla/websocket"
)

type Room struct {
	roomID         int
	connManager    *game.ConnManager
	mahjongManager *game.MahjongManager

	gameRecvCh   chan message.GameMsgRecv
	gameSendCh   chan message.GameMsgSend
	tableOrderCh chan int
	gameResultCh chan message.GameResultMsg
	getReadyCh   chan message.GetReadyMsg
	chatCh       chan message.ChatMsg

	destroyCh chan int

	password    string
	playerCount int
}

func NewRoom(roomID int, rule rule.MahjongRule, destroyCh chan int) *Room {
	gameRecvChannel := make(chan message.GameMsgRecv)
	gameSendChannel := make(chan message.GameMsgSend)
	tableOrderChannel := make(chan int)
	gameResultChannel := make(chan message.GameResultMsg)
	getReadyChannel := make(chan message.GetReadyMsg)
	joinChannel := make(chan message.JoinMsg)
	chatChannel := make(chan message.ChatMsg)
	exitChannel := make(chan bool)
	loopOverChannel := make(chan bool)
	return &Room{
		roomID:         roomID,
		gameRecvCh:     gameRecvChannel,
		gameSendCh:     gameSendChannel,
		tableOrderCh:   tableOrderChannel,
		gameResultCh:   gameResultChannel,
		getReadyCh:     getReadyChannel,
		chatCh:         chatChannel,
		destroyCh:      destroyCh,
		connManager:    game.NewConnManager(roomID, 4, gameRecvChannel, gameSendChannel, tableOrderChannel, gameResultChannel, getReadyChannel, joinChannel, chatChannel, exitChannel, loopOverChannel, destroyCh),
		mahjongManager: game.NewMahjongManager(gameRecvChannel, gameSendChannel, tableOrderChannel, gameResultChannel, getReadyChannel, joinChannel, chatChannel, loopOverChannel, rule),
	}
}

func (r *Room) AddConn(conn *websocket.Conn) {
	r.connManager.AddConn(conn)
}

func (r *Room) AddUserInfo(info message.UserInfo) {
	r.mahjongManager.AddUserInfo(info)
}

func (r *Room) Start() {
	r.mahjongManager.Start()
}
