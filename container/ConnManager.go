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
	conns        []*websocket.Conn

	gameRecvCh   chan message.GameMsgRecv
	gameSendCh   chan message.GameMsgSend
	tableOrderCh chan int
	gameResultCh chan message.GameResultMsg
	getReadyCh   chan message.GetReadyMsg
	joinCh       chan message.JoinMsg
	chatCh       chan message.ChatMsg
	exitCh       chan bool
	loopOverCh   chan bool
	destroyCh    chan int

	roomID int
}

func NewConnManager(roomID int,
	playersCount int,
	gameRecvCh chan message.GameMsgRecv,
	gameSendCh chan message.GameMsgSend,
	tableOrderCh chan int,
	gameResultCh chan message.GameResultMsg,
	getReadyCh chan message.GetReadyMsg,
	joinCh chan message.JoinMsg,
	chatCh chan message.ChatMsg,
	exitCh chan bool,
	loopOverCh chan bool,
	destroyCh chan int) *ConnManager {
	connManager := ConnManager{playersCount: playersCount,
		conns:        make([]*websocket.Conn, 0),
		gameRecvCh:   gameRecvCh,
		gameSendCh:   gameSendCh,
		tableOrderCh: tableOrderCh,
		gameResultCh: gameResultCh,
		getReadyCh:   getReadyCh,
		joinCh:       joinCh,
		chatCh:       chatCh,
		exitCh:       exitCh,
		loopOverCh:   loopOverCh,
		destroyCh:    destroyCh,
		roomID:       roomID,
	}
	go connManager.Broadcaster()
	return &connManager
}

func (m *ConnManager) AddConn(conn *websocket.Conn) {
	index := len(m.conns)
	msg := message.TableOrderMsg{
		MsgType:    config.TableOrderMsgType,
		TableOrder: index,
	}
	conn.WriteJSON(msg)
	m.conns = append(m.conns, conn)
	go m.connListener(conn, m.gameRecvCh, m.chatCh)
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
		case msg := <-m.joinCh:
			log.Println("broadcast join message")
			for _, conn := range m.conns {
				conn.WriteJSON(msg)
			}
		case _ = <-m.exitCh:
			log.Println("broadcaster exit")
			m.destroyCh <- m.roomID
			return
		}
	}
}

func (m *ConnManager) SendMsgTo(order int, msg []byte) {
	if order >= len(m.conns) {
		return
	}
	m.conns[order].WriteMessage(websocket.TextMessage, msg)
}

func (m *ConnManager) checkConn(c *websocket.Conn) {
	var user int
	for index, conn := range m.conns {
		if c == conn {
			log.Println("user: ", index, " disconnected")
			user = index
		}
	}
	u := user
	user = u
	for _, conn := range m.conns {
		conn.Close()
	}
	m.exitCh <- true
	m.loopOverCh <- true
}

func (m *ConnManager) connListener(conn *websocket.Conn, gameRecvCh chan message.GameMsgRecv, chatCh chan message.ChatMsg) {
	for {
		log.Println("waiting msg")
		commonMsg := message.CommonMsg{}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			//TODO: handle when user force exit
			m.checkConn(conn)
			log.Println("failed to read message")
			return
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
		case config.ChatMsgType:
			log.Println("got a chat msg")
			var chatMsg message.ChatMsg
			err = json.Unmarshal(msg, &chatMsg)
			if err != nil {
				panic("fail to decode chat msg")
			}
			for _, conn := range m.conns {
				conn.WriteJSON(chatMsg)
			}
		}
	}
}
