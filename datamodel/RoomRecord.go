package datamodel

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/postgres"

type Record struct {
	gorm.Model
	RoomNumber	int
	Password	string
	UserID0	int
	UserID1 int
	UserID2	int
	UserID3 int
}