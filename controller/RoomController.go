package controller

import (
	"net/http"
)

func RoomPostHandler(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}
	form := request.Form
	userID:=form["UserID"]
	roomName :=form["RoomName"]
	PassWord := form["Password"]
	GetHall().CreateRoom(userID, roomName, PassWord)
}
