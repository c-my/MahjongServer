package main

import (
	"encoding/json"
	"github.com/c-my/MahjongServer/container"
	"github.com/c-my/MahjongServer/controller"
	"github.com/c-my/MahjongServer/rule"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var room *container.Room

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("WS: connect request received")
	conn, _ := upgrader.Upgrade(w, r, nil)
	room.AddConn(conn)
}

func main() {
	configure := readConfig()
	room = container.NewRoom(rule.NewJinzhouRule())
	room.Start()

	router := mux.NewRouter()

	router.HandleFunc("/register/", controller.UserCreateHandler).Methods("POST")
	router.HandleFunc("/login/", controller.UserLoginHandler).Methods("POST")
	router.HandleFunc("/room/", controller.RoomCreateHandler).Methods("POST")
	router.HandleFunc("/room/", controller.RoomJoinHandler).Methods("PUT")
	router.HandleFunc("/ws/{userID}", controller.WsHandler)
	//router.HandleFunc("/", wsHandler)

	http.ListenAndServe(configure.URL+":"+configure.Port, router)
}

func readConfig() Configure {
	content, err := ioutil.ReadFile("configure.json")
	if err != nil {
		log.Fatal("couldn't read configure.json: " + err.Error())
	}
	var conf Configure
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal("couldn't parse configure.json: " + err.Error())
	}
	return conf
}

type Configure struct {
	DatabaseURL    string `json:"DB_URL"`
	DatabasePort   string `json:"DB_PORT"`
	DatabaseName   string `json:"DB_NAME"`
	DatabaseUser   string `json:"DB_USER"`
	DatabasePasswd string `json:"DB_PASSWD"`

	URL  string `json:"HTTP_URL"`
	Port string `json:"HTTP_PORT"`
}
