package controller

import (
	"github.com/c-my/MahjongServer/container"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("WS: connect request received")
	vars := mux.Vars(r)
	user, ok := vars["userID"]
	hall := container.GetHall()
	conn, _ := upgrader.Upgrade(w, r, nil)

	if !ok {
		log.Println("Ignore: no userID provided")
		conn.Close()
		return
	}
	userID, err := strconv.Atoi(user)
	if err != nil {
		log.Println("Ignore: userID[", user, "] isn't a number")
		conn.Close()
		return
	}
	roomID := hall.GetRoomID(userID)
	if roomID < 0 {
		log.Println("Ignore: no room for user[", userID, "]")
		conn.Close()
		return
	}
	hall.AddConn(roomID, conn)
}
