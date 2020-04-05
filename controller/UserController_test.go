package controller

import (
	"bytes"
	"encoding/json"
	"github.com/c-my/MahjongServer/datamodel"
	"github.com/c-my/MahjongServer/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserCreateHandler(t *testing.T) {
	name := time.Now().String()
	body, _ := json.Marshal(createUserMsg{
		UserName: name,
		Password: "qwer123",
	})
	req, err := http.NewRequest("POST", "/register/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserCreateHandler)
	handler.ServeHTTP(rr, req)

	var msg userResult
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)
	if err != nil || msg.Success != true {
		t.Errorf("failed to create user")
	}

	req, err = http.NewRequest("POST", "/register/", bytes.NewReader(body))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)

	if err != nil || msg.Success == true {
		t.Errorf("create a repeated user")
	}
}

func TestUserLoginHandler(t *testing.T) {
	//login fail test
	body, _ := json.Marshal(loginMsg{
		UserName: time.Now().String(),
		Password: "ppp",
	})
	req, err := http.NewRequest("POST", "/login/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserLoginHandler)
	handler.ServeHTTP(rr, req)
	var msg userResult
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)
	if err != nil || msg.Success == true {
		t.Errorf("login as nonexisten user")
	}

	//login success test
	name := time.Now().String()
	repository.UserRepo.Append(datamodel.User{
		UserName: name,
		Password: "pass",
	})
	body, _ = json.Marshal(loginMsg{
		UserName: name,
		Password: "pass",
	})
	req, err = http.NewRequest("POST", "/login/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(UserLoginHandler)
	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.String()), &msg)
	if err != nil || msg.Success != true {
		t.Errorf("login failed")
	}
}
