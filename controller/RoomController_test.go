package controller

import (
	"bytes"
	"encoding/json"
	"github.com/c-my/MahjongServer/container"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoomCreateHandler(t *testing.T) {
	body, _ := json.Marshal(createRoomMsg{UserID: 123, RoomName: "test room", Passwd: "ppp"})
	req, err := http.NewRequest("POST", "/room", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RoomCreateHandler)
	handler.ServeHTTP(rr, req)

	var msg result
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)

	if err != nil || msg.Success != true {
		t.Errorf("failed to create first room")
	} else {
		log.Print("roomID: ", msg.RoomID)
	}
}

func TestRoomJoinHandler(t *testing.T) {
	roomID := container.GetHall().CreateRoom(003, "test", "pp")
	body, _ := json.Marshal(joinRoomMsg{
		UserID:   004,
		RoomID:   roomID,
		Password: "pp",
	})
	req, err := http.NewRequest("PUT", "/room", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RoomJoinHandler)
	handler.ServeHTTP(rr, req)

	var msg result
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)
	if err != nil || msg.Success != true {
		t.Errorf("failed to join room as 2nd player")
	}

	body, _ = json.Marshal(joinRoomMsg{
		UserID: 005,
		RoomID: roomID,
	})
	req, err = http.NewRequest("PUT", "/room", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(RoomJoinHandler)
	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)
	if msg.Success == true {
		t.Errorf("join room with wrong password")
	}
	body, _ = json.Marshal(joinRoomMsg{
		UserID:   005,
		RoomID:   roomID,
		Password: "pp",
	})
	req, err = http.NewRequest("PUT", "/room", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(RoomJoinHandler)
	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)
	if err != nil || msg.Success != true {
		t.Errorf("failed to join room as 3rd player")
	}
	body, _ = json.Marshal(joinRoomMsg{
		UserID:   006,
		RoomID:   roomID,
		Password: "pp",
	})
	req, err = http.NewRequest("PUT", "/room", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(RoomJoinHandler)
	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)
	if err != nil || msg.Success != true {
		t.Errorf("failed to join room as 4th player")
	}
	body, _ = json.Marshal(joinRoomMsg{
		UserID:   006,
		RoomID:   roomID,
		Password: "pp",
	})
	req, err = http.NewRequest("PUT", "/room", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(RoomJoinHandler)
	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)
	if err != nil || msg.Success == true {
		t.Errorf("joined room as 5th player")
	}
}
