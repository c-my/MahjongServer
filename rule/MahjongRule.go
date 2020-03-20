package rule

import "github.com/c-my/MahjongServer/model"

type MahjongRule interface {
	GenerateTiles() []model.Tile
	CanChow(hand []model.Tile, newTile model.Tile) (bool, []int)
	CanPong(hand []model.Tile, newTile model.Tile) bool
	CanExposedKong(hand []model.Tile, newTile model.Tile) bool
	CanConcealedKong(hand []model.Tile, newTile model.Tile) bool
	CanAddedKong(shown []model.ShownTile, newTile model.Tile) bool
	CanWin(handTiles []model.Tile, showTiles []model.ShownTile, newTile model.Tile) bool
}
