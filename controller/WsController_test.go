package controller

import (
	"github.com/c-my/MahjongServer/container"
	"github.com/c-my/MahjongServer/rule"
	"github.com/gorilla/websocket"
	"net/url"
	"testing"
)

func TestWsHandler(t *testing.T) {
	container.GetHall().CreateRoom(100, "100", rule.NewJinzhouRule())
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:1114", Path: "/100"}
	_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err == nil {
		t.Errorf("connected via wrong userid")
	}
	u = url.URL{Scheme: "ws", Host: "127.0.0.1:1114", Path: "/some_string"}
	_, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err == nil {
		t.Errorf("connected via non-number param")
	}

}
