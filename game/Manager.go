package game

import "github.com/c-my/MahjongServer/model"

type Manager interface {
	GenerateTiles() []model.Tile
	CanChow(tiles []model.Tile, newTile model.Tile) bool
	CanPong(tiles []model.Tile, newTile model.Tile) bool
	CanKong(tiles []model.Tile, newTile model.Tile) bool
	CanWin(tiles []model.Tile, newTile model.Tile) bool
}
