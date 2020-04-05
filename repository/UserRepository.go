package repository

import (
	"github.com/c-my/MahjongServer/datamodel"
	"github.com/c-my/MahjongServer/datasource"
	"github.com/jinzhu/gorm"
)

var (
	UserRepo = NewUserRepository()
)

type UserRepository struct {
	source *gorm.DB
}

// SelectByID selects user by uid
func (r *UserRepository) SelectByUsername(username string) (user datamodel.User, notfound bool) {
	notfound = r.source.Where("username = ?", username).First(&user).RecordNotFound()
	return
}

func (r *UserRepository) Append(user datamodel.User) bool {
	var u datamodel.User
	if err := r.source.Where("username = ?", user.UserName).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			r.source.Create(&user)
			return true
		}
	}
	return false
}

// NewUserRepository is
func NewUserRepository() *UserRepository {
	db := datasource.DB
	if !db.HasTable(&datamodel.User{}) {
		db.CreateTable(&datamodel.User{})
	}
	return &UserRepository{source: db}
}
