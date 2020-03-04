package main

import (
	"encoding/json"
	"github.com/c-my/MahjongServer/model"
	"github.com/gorilla/websocket"
	"log"
)

type ConnManager struct {
	playersCount int
	conns        []websocket.Conn

	gameRecvCh chan model.GameMsgRecv
	gameSendCh chan model.GameMsgSend
}

func NewConnManager(playersCount int, gameRecvCh chan model.GameMsgRecv, gameSendCh chan model.GameMsgSend) *ConnManager {
	return &ConnManager{playersCount: playersCount,
		conns:      make([]websocket.Conn, 0),
		gameRecvCh: gameRecvCh,
		gameSendCh: gameSendCh,
	}
}

func (m *ConnManager) SetConn(conn *websocket.Conn) {
	m.conns = append(m.conns, *conn)
	go connListener(conn, m.gameRecvCh)
}

func connListener(conn *websocket.Conn, gameCh chan model.GameMsgRecv) {
	for {
		commonMsg := model.CommonMsg{}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message")
		}
		err = json.Unmarshal(msg, &commonMsg)
		if err != nil {
			panic("not a valid message")
		}

		// estimate type of message
		switch commonMsg.MsgType {
		case model.GameMsgType:
			log.Println("got a game msg")
			var gameMsg model.GameMsgRecv
			err = json.Unmarshal(msg, &gameMsg)
			if err != nil {
				panic("fail to decode game msg")
			}
			gameCh <- gameMsg
			gameMsgSend := <-gameCh
			conn.WriteJSON(gameMsgSend)
		}
	}
}
