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

	gameRecvCh   chan message.GameMsgRecv
	gameSendCh   chan message.GameMsgSend
	tableOrderCh chan int
	gameResultCh chan message.GameResultMsg
	getReadyCh   chan message.GetReadyMsg
}

func NewConnManager(playersCount int,
	gameRecvCh chan message.GameMsgRecv,
	gameSendCh chan message.GameMsgSend,
	tableOrderCh chan int,
	gameResultCh chan message.GameResultMsg,
	getReadyCh chan message.GetReadyMsg) *ConnManager {
	connManager := ConnManager{playersCount: playersCount,
		conns:        make([]websocket.Conn, 0),
		gameRecvCh:   gameRecvCh,
		gameSendCh:   gameSendCh,
		tableOrderCh: tableOrderCh,
		gameResultCh: gameResultCh,
		getReadyCh:   getReadyCh,
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
		case first := <-m.tableOrderCh:
			log.Println("send table order")
			for i := 0; i < 4; i++ {
				msg := message.TableOrderMsg{
					MsgType:    config.TableOrderMsgType,
					TableOrder: i,
				}
				m.conns[(first+i)%4].WriteJSON(msg)
			}
		case msg := <-m.gameResultCh:
			log.Println("broadcast game result")
			for _, conn := range m.conns {
				conn.WriteJSON(msg)
			}
		case msg := <-m.getReadyCh:
			log.Println("broadcast ready message")
			for _, conn := range m.conns {
				conn.WriteJSON(msg)
			}
		}
	}
}

func (m *ConnManager) SendMsgTo(order int, msg []byte) {
	if order >= len(m.conns) {
		return
	}
	m.conns[order].WriteMessage(websocket.TextMessage, msg)
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
