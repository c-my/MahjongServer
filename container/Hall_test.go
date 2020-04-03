package container

import "testing"

func TestGetHall(t *testing.T) {
	hall := GetHall()
	if hall == nil {
		t.Errorf("get nil hall")
	}
	hall2 := GetHall()
	if hall != hall2 {
		t.Errorf("get different hall")
	}
}

func TestHall_CreateRoom(t *testing.T) {
	hall := GetHall()
	roomID := hall.CreateRoom(0001, "ppp")
	if roomID == -1 {
		t.Errorf("failed to create first room")
	}
}

func TestHall_JoinRoom(t *testing.T) {
	hall := GetHall()
	roomID := hall.CreateRoom(0002, "ppp")
	if hall.JoinRoom(0003, roomID, "ppp") != Success {
		t.Errorf("2nd player joined failed")
	}
	if hall.JoinRoom(0004, roomID, "ppp") != Success {
		t.Errorf("3rd player joined failed")
	}
	if hall.JoinRoom(0005, roomID, "pppaaa") == Success {
		t.Errorf("joined room via wrong password")
	}
	if hall.JoinRoom(0005, roomID, "ppp") != Success {
		t.Errorf("4th player joined failed")
	}
	if hall.JoinRoom(0006, roomID, "ppp") == Success {
		t.Errorf("5th player joined room")
	}
}

func TestHall_GetRoomID(t *testing.T) {
	hall := GetHall()
	roomID := hall.CreateRoom(0007, "ppp")
	if hall.GetRoomID(7) != roomID {
		t.Errorf("get wrong roomID")
	}
	if hall.GetRoomID(8) >= 0 {
		t.Errorf("get ghost roomID")
	}
	hall.JoinRoom(8, roomID, "ppp")
	if hall.GetRoomID(8) != roomID {
		t.Errorf("get wrong roomID")
	}
}
