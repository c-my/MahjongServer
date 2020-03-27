package container

import (
	"encoding/json"
	"github.com/c-my/MahjongServer/config"
	"github.com/c-my/MahjongServer/message"
	"github.com/gorilla/websocket"
	"log"
)

type ConnManager struct {
	playersCount int
	conns        []websocket.Conn

	gameRecvCh chan message.GameMsgRecv
	gameSendCh chan message.GameMsgSend
}

func NewConnManager(playersCount int, gameRecvCh chan message.GameMsgRecv, gameSendCh chan message.GameMsgSend) *ConnManager {
	connManager := ConnManager{playersCount: playersCount,
		conns:      make([]websocket.Conn, 0),
		gameRecvCh: gameRecvCh,
		gameSendCh: gameSendCh,
	}
	go connManager.Broadcaster()
	return &connManager
}

func (m *ConnManager) SetConn(conn *websocket.Conn) {
	m.conns = append(m.conns, *conn)
	go connListener(conn, m.gameRecvCh, m.gameSendCh)
}

func (m *ConnManager) Broadcaster() {
	for {
		select {
		case msg := <-m.gameSendCh:
			log.Println("broadcast game message")
			for _, conn := range m.conns {
				conn.WriteJSON(msg)
			}
		}
	}
}

func connListener(conn *websocket.Conn, gameRecvCh chan message.GameMsgRecv, gameSendCh chan message.GameMsgSend) {
	for {
		commonMsg := message.CommonMsg{}
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
		case config.GameMsgType:
			log.Println("got a game msg")
			var gameMsg message.GameMsgRecv
			err = json.Unmarshal(msg, &gameMsg)
			if err != nil {
				panic("fail to decode game msg")
			}
			gameRecvCh <- gameMsg
		}
	}
}
