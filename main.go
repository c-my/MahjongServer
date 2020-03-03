package main

import (
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

var room *Room

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("connection received")
	conn, _ := upgrader.Upgrade(w, r, nil)
	room.AddConn(conn)
}

func main() {
	room = NewRoom()
	room.Start()

	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:1114", nil)
}
