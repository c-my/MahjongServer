package main

import "github.com/c-my/MahjongServer/rule"

type Hall struct {
	rooms   map[int]*Room
	players map[int]int

	maxRoomSize int
}

var hall *Hall = nil

func GetHall() *Hall {
	if hall != nil {
		return hall
	} else {
		hall = NewHall(100)
		return hall
	}
}

func NewHall(maxRoomSize int) *Hall {
	return &Hall{maxRoomSize: maxRoomSize}
}

func (h *Hall) CreateRoom(userID int, roomName string, passwd string) int {
	for i := 0; i < h.maxRoomSize; i++ {
		if !h.hasRoom(i) {
			h.rooms[i] = NewRoom(rule.NewJinzhouRule())
			h.players[userID] = i
			return i
		}
	}
	return -1
}

func (h *Hall) hasRoom(roomID int) bool {
	_, ok := h.rooms[roomID]
	return ok
}

func (h *Hall) destroyRoom(roomID int) {
	delete(h.rooms, roomID)
}
