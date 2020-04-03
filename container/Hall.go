package container

import (
	"github.com/c-my/MahjongServer/rule"
	"github.com/gorilla/websocket"
)

const (
	Success       = 1
	NoSuchRoom    = -1
	NoSeat        = -2
	WrongPassword = -3
)

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
	return &Hall{maxRoomSize: maxRoomSize, rooms: map[int]*Room{}, players: map[int]int{}}
}

func (h *Hall) CreateRoom(userID int, passwd string) int {
	for i := 0; i < h.maxRoomSize; i++ {
		if !h.hasRoom(i) {
			room := NewRoom(rule.NewJinzhouRule())
			room.Start()
			room.password = passwd
			room.playerCount = 1
			h.rooms[i] = room
			h.players[userID] = i
			return i
		}
	}
	return -1
}

func (h *Hall) JoinRoom(userID int, roomID int, password string) int {
	if !h.hasRoom(roomID) {
		return NoSuchRoom
	}
	room := h.rooms[roomID]
	if room.playerCount >= 4 {
		return NoSeat
	}
	if room.password != password {
		return WrongPassword
	}
	h.players[userID] = roomID
	room.playerCount = room.playerCount + 1
	return Success
}

func (h *Hall) GetAllRoom() map[int]*Room {
	return h.rooms
}

func (h *Hall) GetRoomID(userID int) int {
	roomID, ok := h.players[userID]
	if !ok {
		return -1
	} else {
		return roomID
	}
}

func (h *Hall) AddConn(roomID int, conn *websocket.Conn) {
	h.rooms[roomID].AddConn(conn)
}

func (h *Hall) hasRoom(roomID int) bool {
	_, ok := h.rooms[roomID]
	return ok
}

func (h *Hall) destroyRoom(roomID int) {
	delete(h.rooms, roomID)
}
