package controller

import (
	"encoding/json"
	"github.com/c-my/MahjongServer/datamodel"
	"github.com/c-my/MahjongServer/repository"
	"log"
	"net/http"
)

type createUserMsg struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Gender   int    `json:"gender"`
}

type loginMsg struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type userResult struct {
	Success bool `json:"success"`
	UserID  uint `json:"user_id"`
}

func UserCreateHandler(writer http.ResponseWriter, request *http.Request) {
	var msg createUserMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		return
	}
	newUser := datamodel.User{
		UserName: msg.UserName,
		Password: msg.Password,
		NickName: msg.Nickname,
		Gender:   msg.Gender,
	}
	success, newID := repository.UserRepo.Append(newUser)
	if !success { //user already exist
		//return fail result
		log.Print("register failed: user already exist")
		res, _ := json.Marshal(userResult{false, 0})
		writer.Write(res)
		return
	} else {
		log.Print("register success")
		res, _ := json.Marshal(userResult{true, newID})
		writer.Write(res)
	}
}

func UserLoginHandler(writer http.ResponseWriter, request *http.Request) {
	var msg createUserMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		return
	}
	loginUser := datamodel.User{
		UserName: msg.UserName,
		Password: msg.Password,
	}
	u, notFound := repository.UserRepo.SelectByUsername(loginUser.UserName)
	if notFound {
		log.Print("login failed: user not exist")
		res, _ := json.Marshal(userResult{false, 0})
		writer.Write(res)
		return
	}
	if encodePassword(loginUser.Password) != u.Password {
		log.Print("login failed: wrong password")
		res, _ := json.Marshal(userResult{false, 0})
		writer.Write(res)
		return
	}
	log.Print("login success")
	res, _ := json.Marshal(userResult{true, u.ID})
	writer.Write(res)
}

func encodePassword(password string) string {
	//TODO:use a real encoder
	return password
}
