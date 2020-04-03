package controller

import (
	"encoding/json"
	"github.com/c-my/MahjongServer/container"
	"log"
	"net/http"
)

type createRoomMsg struct {
	UserID int    `json:"user_id"`
	Passwd string `json:"passwd"`
}

type joinRoomMsg struct {
	UserID   int    `json:"user_id"`
	RoomID   int    `json:"room_id"`
	Password string `json:"passwd"`
}

type result struct {
	Success bool `json:"success"`
	Reason  int  `json:"reason"`
}

func RoomCreateHandler(writer http.ResponseWriter, request *http.Request) {
	var msg createRoomMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		j, _ := json.Marshal(failMsg(0))
		log.Print("failed to create room: unknown request")
		writer.Write(j)
	}
	userID := msg.UserID
	PassWord := msg.Passwd
	roomID := container.GetHall().CreateRoom(userID, PassWord)
	if roomID == -1 {
		log.Print("failed to create room: max room size reached")
		j, _ := json.Marshal(failMsg(0))
		writer.Write(j)
	}
	log.Print("user[", userID, "] created room: ", roomID, " with password[", PassWord, "]")
	j, _ := json.Marshal(successMsg())
	writer.Write(j)
}

func RoomJoinHandler(writer http.ResponseWriter, request *http.Request) {
	var msg joinRoomMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		j, _ := json.Marshal(failMsg(0))
		log.Print("failed to join room: unknown request")
		writer.Write(j)
	}
	userID := msg.UserID
	roomID := msg.RoomID
	password := msg.Password

	success := container.GetHall().JoinRoom(userID, roomID, password)
	if success != container.Success {
		log.Print("failed to join room, reason: ", success)
		j, _ := json.Marshal(failMsg(success))
		writer.Write(j)
	} else {
		j, _ := json.Marshal(successMsg())
		writer.Write(j)
	}
}

func successMsg() result {
	return result{Success: true, Reason: 1}
}

func failMsg(reason int) result {
	return result{Success: false, Reason: reason}
}
