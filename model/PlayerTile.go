package model

import "sort"

type PlayerTile struct {
	handTiles []Tile
	showTiles []Tile
}

func (t *PlayerTile) SetHandTiles(tiles []Tile) {
	t.handTiles = tiles
	t.SortHand()
}

func (t *PlayerTile) GetHandTiles() []Tile {
	return t.handTiles
}

func (t *PlayerTile) SortHand() {
	sort.Slice(t.handTiles, func(i, j int) bool {
		iValue := int(t.handTiles[i].Suit)*10 + t.handTiles[i].Number
		jValue := int(t.handTiles[j].Suit)*10 + t.handTiles[j].Number
		return iValue < jValue
	})
}
