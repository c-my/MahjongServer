package game

import (
	"github.com/c-my/MahjongServer/model"
	"log"
)

type JinzhouMahjong struct {
	wall       *model.Wall
	playerTile []model.PlayerTile

	gameRecvCh chan model.GameMsgRecv
	gameSendCh chan model.GameMsgSend
}

//NewJinzhouMahjong return new Mahjong instance, with initiated game channels
// and tile-wall
func NewJinzhouMahjong(gameRecvCh chan model.GameMsgRecv, gameSendCh chan model.GameMsgSend) *JinzhouMahjong {
	mahjong := JinzhouMahjong{}
	mahjong.wall = model.NewWall(mahjong.GenerateTiles)
	mahjong.playerTile = make([]model.PlayerTile, 4)
	mahjong.gameRecvCh = gameRecvCh
	mahjong.gameSendCh = gameSendCh
	return &mahjong
}

//Start will start the main loop of a game in a new goroutine
func (m *JinzhouMahjong) Start() {
	go gameLoop(m.gameRecvCh, m.gameSendCh)
}

func (j JinzhouMahjong) GenerateTiles() []model.Tile {
	tiles := make([]model.Tile, 0)
	// generate all tiles
	for s := 0; s < 4; s++ {
		for i := 1; i <= 9; i++ {
			tiles = append(tiles, model.Tile{model.Character, i},
				model.Tile{model.Bamboo, i},
				model.Tile{model.Dot, i})
		}
		tiles = append(tiles, model.Tile{model.Dragon, 1})
	}
	return tiles
}

func (j JinzhouMahjong) GetFirstPlayer() {
	panic("implement me")
}

func (j JinzhouMahjong) CanChow(tiles []model.Tile, newTile model.Tile) bool {
	if newTile.Suit != model.Character && newTile.Suit != model.Bamboo && newTile.Suit != model.Dot {
		return false
	}
	panic("implement me")
}

func (j JinzhouMahjong) CanPong(tiles []model.Tile, newTile model.Tile) bool {
	panic("implement me")
}

func (j JinzhouMahjong) CanKong(tiles []model.Tile, newTile model.Tile) bool {
	panic("implement me")
}

func (j JinzhouMahjong) CanWin(tiles []model.Tile, newTile model.Tile) bool {
	panic("implement me")
}

func gameLoop(gameRecvCh chan model.GameMsgRecv, gameSendCh chan model.GameMsgSend) {
	for {
		select {
		case msg := <-gameRecvCh:
			log.Println("receive from game loop: ", msg)
			// TODO:main logic here
			var msgSend model.GameMsgSend
			gameSendCh <- msgSend
		}
	}
}
