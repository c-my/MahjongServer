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

func (r *UserRepository) SelectByID(userID int) (user datamodel.User, notfound bool) {
	notfound = r.source.Where("ID = ?", userID).First(&user).RecordNotFound()
	return
}

func (r *UserRepository) AddWinRecord(userID uint) {
	var user datamodel.User
	notfound := r.source.Where("ID = ?", userID).First(&user).RecordNotFound()
	if notfound {
		return
	}
	user.WinTime = user.WinTime + 1
	r.source.Save(&user)
}

func (r *UserRepository) AddGameRecord(userID uint) {
	var user datamodel.User
	notfound := r.source.Where("ID = ?", userID).First(&user).RecordNotFound()
	if notfound {
		return
	}
	user.GameTime = user.GameTime + 1
	r.source.Save(&user)
}

func (r *UserRepository) Append(user datamodel.User) (bool, uint) {
	var u datamodel.User
	if err := r.source.Where("username = ?", user.UserName).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			r.source.Create(&user)
			return true, user.ID
		}
	}
	return false, 0
}

// NewUserRepository is
func NewUserRepository() *UserRepository {
	db := datasource.DB
	if !db.HasTable(&datamodel.User{}) {
		db.CreateTable(&datamodel.User{})
	}
	return &UserRepository{source: db}
}
