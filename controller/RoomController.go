package controller

import (
	"encoding/json"
	"github.com/c-my/MahjongServer/container"
	"log"
	"net/http"
)

type createRoomMsg struct {
	UserID   int    `json:"user_id"`
	RoomName string `json:"room_name"`
	Passwd   string `json:"passwd"`
}

type joinRoomMsg struct {
	UserID int `json:"user_id"`
	RoomID int `json:"room_id"`
}

type result struct {
	Success bool `json:"Success"`
	RoomID  int  `json:"room_id"`
}

func RoomCreateHandler(writer http.ResponseWriter, request *http.Request) {
	var msg createRoomMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		j, _ := json.Marshal(failMsg())
		log.Print("failed to create room: unknown request")
		writer.Write(j)
	}
	userID := msg.UserID
	roomName := msg.RoomName
	PassWord := msg.RoomName
	roomID := container.GetHall().CreateRoom(userID, roomName, PassWord)
	if roomID == -1 {
		log.Print("failed to create room: max room size reached")
		j, _ := json.Marshal(failMsg())
		writer.Write(j)
	}
	log.Print("created room:", roomID)
	j, _ := json.Marshal(successMsg(roomID))
	writer.Write(j)
}

func RoomJoinHandler(writer http.ResponseWriter, request *http.Request) {
	var msg joinRoomMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		j, _ := json.Marshal(failMsg())
		log.Print("failed to join room: unknown request")
		writer.Write(j)
	}
	userID := msg.UserID
	roomID := msg.RoomID

	success := container.GetHall().JoinRoom(userID, roomID)
	if !success {
		j, _ := json.Marshal(failMsg())
		writer.Write(j)
	} else {
		j, _ := json.Marshal(successMsg(roomID))
		writer.Write(j)
	}
}

func successMsg(roomID int) result {
	return result{Success: true, RoomID: roomID}
}

func failMsg() result {
	return result{Success: false}
}
