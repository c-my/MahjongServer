package controller

import (
	"github.com/c-my/MahjongServer/container"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"testing"
)

func TestWsHandler(t *testing.T) {
	container.GetHall().CreateRoom(100, "100r", "100")
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:1114", Path: "/100"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	//time.Sleep(1)
	defer c.Close()

}
