package main

import (
	"github.com/c-my/MahjongServer/container"
	"github.com/c-my/MahjongServer/controller"
	"github.com/c-my/MahjongServer/rule"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var room *container.Room

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("WS: connect request received")
	conn, _ := upgrader.Upgrade(w, r, nil)
	room.AddConn(conn)
}

func main() {
	room = container.NewRoom(rule.NewJinzhouRule())
	room.Start()

	router := mux.NewRouter()

	router.HandleFunc("/room/", controller.RoomCreateHandler).Methods("POST")
	router.HandleFunc("/room/", controller.RoomJoinHandler).Methods("PUT")
	router.HandleFunc("/ws/{userID}", controller.WsHandler)
	//router.HandleFunc("/", wsHandler)
	http.ListenAndServe("127.0.0.1:1114", router)
}
