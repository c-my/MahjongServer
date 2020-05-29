package rule

import (
	"github.com/c-my/MahjongServer/config"
	"github.com/c-my/MahjongServer/model"
)

type ShenyangRule struct {
}

func NewShenyangRule() *ShenyangRule {
	return &ShenyangRule{}
}

func (r *ShenyangRule) GenerateTiles() []model.Tile {
	tiles := make([]model.Tile, 0)
	// generate all tiles
	for s := 0; s < 4; s++ {
		for i := 1; i <= 9; i++ {
			tiles = append(tiles, model.Tile{model.Character, i},
				model.Tile{model.Bamboo, i},
				model.Tile{model.Dot, i})
		}
		tiles = append(tiles, model.Tile{model.Dragon, 1})
		tiles = append(tiles, model.Tile{model.Dragon, 2})
		tiles = append(tiles, model.Tile{model.Dragon, 3})
		tiles = append(tiles, model.Tile{model.Wind, 1})
		tiles = append(tiles, model.Tile{model.Wind, 2})
		tiles = append(tiles, model.Tile{model.Wind, 3})
		tiles = append(tiles, model.Tile{model.Wind, 4})
	}
	return tiles
}

func (r *ShenyangRule) CanChow(tiles []model.Tile, newTile model.Tile) (bool, []int) {
	var canChow = false
	var chowTypes []int
	if !newTile.IsSuit() {
		return false, append(chowTypes, config.NAC)
	}
	if r.canLeftChow(tiles, newTile) {
		canChow = true
		chowTypes = append(chowTypes, config.LeftChow)
	}
	if r.canMidChow(tiles, newTile) {
		canChow = true
		chowTypes = append(chowTypes, config.MidChow)
	}
	if r.canRightChow(tiles, newTile) {
		canChow = true
		chowTypes = append(chowTypes, config.RightChow)
	}
	return canChow, chowTypes
}

func (r *ShenyangRule) CanPong(tiles []model.Tile, newTile model.Tile) bool {
	return model.GetTileCount(tiles, newTile) >= 2
}

func (r *ShenyangRule) CanExposedKong(tiles []model.Tile, newTile model.Tile) bool {
	return model.GetTileCount(tiles, newTile) == 3
}

func (r *ShenyangRule) CanConcealedKong(tiles []model.Tile, newTile model.Tile) bool {
	return r.CanExposedKong(tiles, newTile)
}

func (r *ShenyangRule) CanAddedKong(shown []model.ShownTile, newTile model.Tile) bool {
	for _, v := range shown {
		if v.ShownType == config.Pong && v.Tiles[0].Equals(newTile) {
			return true
		}
	}
	return false
}

func (r *ShenyangRule) CanWin(handTiles []model.Tile, showTiles []model.ShownTile, newTile model.Tile) bool {
	//必须开门
	if len(showTiles) == 0 {
		return false
	}

	tarTiles := make([]model.Tile, len(handTiles))
	copy(tarTiles, handTiles)

	tarTiles = append(tarTiles, newTile)
	model.SortTiles(tarTiles)

	var pos = model.GetPairPos(tarTiles)
	if len(pos) == 0 { //没有对
		return false
	}

	//有幺九
	if !r.hasOneOrNine(handTiles, showTiles) {
		return false
	}
	//不缺门
	if !r.CheckDoor(handTiles, showTiles) {
		return false
	}
	var requestedTripletCount = 1
	shownTripCount := model.GetPongCount(showTiles) + model.GetExposedKongCount(showTiles) +
		model.GetAddedKongCount(showTiles) + model.GetConcealedKongCount(showTiles)

	requestedTripletCount -= shownTripCount
	lastPairTile := model.Tile{
		Suit:   model.NAT,
		Number: 0,
	}
	for _, p := range pos {
		if tarTiles[p].Equals(lastPairTile) {
			continue
		} else {
			lastPairTile = tarTiles[p]
		}
		cards := model.RemovePair(tarTiles, p)
		if model.IsAllSeqOrTriplet(cards, requestedTripletCount) {
			//中发白为对子，可以免去杠
			if lastPairTile.Suit == model.Dragon {
				return true
			}
			for _, v := range showTiles {
				if v.ShownType == config.AddedKong ||
					v.ShownType == config.ExposedKong ||
					v.ShownType == config.ConcealedKong {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (r *ShenyangRule) hasOneOrNine(handTiles []model.Tile, showTiles []model.ShownTile) bool {
	for _, v := range handTiles {
		if v.Suit == model.Dragon || v.Suit == model.Wind || v.Number == 1 || v.Number == 9 {
			return true
		}
	}
	for _, v := range showTiles {
		for _, t := range v.Tiles {
			if t.Suit == model.Dragon || t.Suit == model.Wind || t.Number == 1 || t.Number == 9 {
				return true
			}
		}
	}
	return false
}
func (r *ShenyangRule) CheckDoor(hand []model.Tile, shown []model.ShownTile) bool {
	hasCharacter := false
	hasBamboo := false
	hasDot := false
	for _, t := range hand {
		switch t.Suit {
		case model.Character:
			hasCharacter = true
		case model.Bamboo:
			hasBamboo = true
		case model.Dot:
			hasDot = true
		}
	}
	for _, v := range shown {
		for _, t := range v.Tiles {
			switch t.Suit {
			case model.Character:
				hasCharacter = true
			case model.Bamboo:
				hasBamboo = true
			case model.Dot:
				hasDot = true
			}
		}
	}
	return hasCharacter && hasBamboo && hasDot
}

func (r *ShenyangRule) canLeftChow(tiles []model.Tile, newTile model.Tile) bool {
	right := newTile.GetRightTile()
	rright := right.GetRightTile()
	if right == nil || rright == nil {
		return false
	}
	return model.GetTileCount(tiles, *right) != 0 && model.GetTileCount(tiles, *rright) != 0
}

func (r *ShenyangRule) canMidChow(tiles []model.Tile, newTile model.Tile) bool {
	left := newTile.GetLeftTile()
	right := newTile.GetRightTile()
	if left == nil || right == nil {
		return false
	}
	return model.GetTileCount(tiles, *left) != 0 && model.GetTileCount(tiles, *right) != 0
}

func (r *ShenyangRule) canRightChow(tiles []model.Tile, newTile model.Tile) bool {
	left := newTile.GetLeftTile()
	lleft := left.GetLeftTile()
	if left == nil || lleft == nil {
		return false
	}
	return model.GetTileCount(tiles, *left) != 0 && model.GetTileCount(tiles, *lleft) != 0
}
