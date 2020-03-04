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
	go m.gameLoop(m.gameRecvCh, m.gameSendCh)
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

func (j *JinzhouMahjong) gameLoop(gameRecvCh chan model.GameMsgRecv, gameSendCh chan model.GameMsgSend) {
	for {
		select {
		case msg := <-gameRecvCh:
			log.Println("receive from game loop: ", msg)
			switch msg.Action {
			case model.Start:
				j.dealTile()
				availableActions := make([]int, 0)
				availableActions = append(availableActions, model.Deal)
				tilesCount := make([]int, 0)
				for i := 0; i < 4; i++ {
					tilesCount = append(tilesCount, len(j.playerTile[i].HandTiles))
				}
				var msgSend model.GameMsgSend
				msgSend.MsgType = model.GameMsgType
				msgSend.AvailableActions = availableActions
				msgSend.CurrentTile = j.playerTile[0].HandTiles
				msgSend.TilesCount = tilesCount
				gameSendCh <- msgSend
			}
			// TODO:main logic here

		}
	}
}

func (j *JinzhouMahjong) dealTile() {
	for i := 0; i < 3; i++ {
		j.playerTile[0].HandTiles = append(j.playerTile[0].HandTiles, j.wall.FrontDraw())
		j.playerTile[0].HandTiles = append(j.playerTile[0].HandTiles, j.wall.FrontDraw())
		j.playerTile[0].HandTiles = append(j.playerTile[0].HandTiles, j.wall.FrontDraw())
		j.playerTile[0].HandTiles = append(j.playerTile[0].HandTiles, j.wall.FrontDraw())
		j.playerTile[1].HandTiles = append(j.playerTile[1].HandTiles, j.wall.FrontDraw())
		j.playerTile[1].HandTiles = append(j.playerTile[1].HandTiles, j.wall.FrontDraw())
		j.playerTile[1].HandTiles = append(j.playerTile[1].HandTiles, j.wall.FrontDraw())
		j.playerTile[1].HandTiles = append(j.playerTile[1].HandTiles, j.wall.FrontDraw())
		j.playerTile[2].HandTiles = append(j.playerTile[2].HandTiles, j.wall.FrontDraw())
		j.playerTile[2].HandTiles = append(j.playerTile[2].HandTiles, j.wall.FrontDraw())
		j.playerTile[2].HandTiles = append(j.playerTile[2].HandTiles, j.wall.FrontDraw())
		j.playerTile[2].HandTiles = append(j.playerTile[2].HandTiles, j.wall.FrontDraw())
		j.playerTile[3].HandTiles = append(j.playerTile[3].HandTiles, j.wall.FrontDraw())
		j.playerTile[3].HandTiles = append(j.playerTile[3].HandTiles, j.wall.FrontDraw())
		j.playerTile[3].HandTiles = append(j.playerTile[3].HandTiles, j.wall.FrontDraw())
		j.playerTile[3].HandTiles = append(j.playerTile[3].HandTiles, j.wall.FrontDraw())
	}
	j.playerTile[0].HandTiles = append(j.playerTile[0].HandTiles, j.wall.FrontDraw())
	j.playerTile[1].HandTiles = append(j.playerTile[1].HandTiles, j.wall.FrontDraw())
	j.playerTile[2].HandTiles = append(j.playerTile[2].HandTiles, j.wall.FrontDraw())
	j.playerTile[3].HandTiles = append(j.playerTile[3].HandTiles, j.wall.FrontDraw())

	for i := 0; i <= 3; i++ {
		j.playerTile[i].SortHand()
	}
}
