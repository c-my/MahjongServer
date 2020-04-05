package datamodel

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/postgres"

type User struct {
	gorm.Model
	UserName string `gorm:"UNIQUE;Column:username"`
	Password string
}
