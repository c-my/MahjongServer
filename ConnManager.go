package main

import (
	"github.com/c-my/MahjongServer/model"
	"github.com/gorilla/websocket"
	"log"
)

type ConnManager struct {
	playersCount int
	conns        []websocket.Conn

	gameCh chan model.GameMessage
}

func NewConnManager(playersCount int, gameCh chan model.GameMessage) *ConnManager {
	return &ConnManager{playersCount: playersCount,
		conns:  make([]websocket.Conn, 0),
		gameCh: gameCh,
	}
}

func (m *ConnManager) SetConn(conn *websocket.Conn) {
	m.conns = append(m.conns, *conn)
	go connListener(conn, m.gameCh)
}

func connListener(conn *websocket.Conn, ch chan model.GameMessage) {
	for {
		//var msg model.GameMessage
		_,msg,_:=conn.ReadMessage()

		//err := conn.ReadJSON(&msg)
		//if err != nil {
		//	panic("error reading json")
		//}
		// TODO: judge if is game message
		log.Println("receive message: ", string(msg))
		ch <- model.GameMessage(msg)
		msg = []byte(<-ch)
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}
