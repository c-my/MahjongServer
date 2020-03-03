package game

import (
	"github.com/c-my/MahjongServer/model"
	"log"
)

type JinzhouMahjong struct {
	wall       *model.Wall
	playerTile []model.PlayerTile

	gameCh chan model.GameMessage
}

func NewJinzhouMahjong(gameCh chan model.GameMessage) *JinzhouMahjong {
	mahjong := JinzhouMahjong{}
	mahjong.wall = model.NewWall(mahjong.GenerateTiles)
	mahjong.playerTile = make([]model.PlayerTile, 4)
	mahjong.gameCh = gameCh
	return &mahjong
}

func (m *JinzhouMahjong) Start() {
	go gameLoop(m.gameCh)
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

func gameLoop(gameCh chan model.GameMessage){
	for {
		select {
		case msg := <-gameCh:
			log.Println("receive from game loop: ", msg)
			// TODO:main logic here
			gameCh<-"game received"
		}
	}
}
