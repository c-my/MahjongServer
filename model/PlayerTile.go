package model

import "sort"

type PlayerTile struct {
	HandTiles []Tile
	ShowTiles []Tile
}

func (t *PlayerTile) SetHandTiles(tiles []Tile) {
	t.HandTiles = tiles
	t.SortHand()
}

func (t *PlayerTile) GetHandTiles() []Tile {
	return t.HandTiles
}

func (t *PlayerTile) SortHand() {
	sort.Slice(t.HandTiles, func(i, j int) bool {
		iValue := int(t.HandTiles[i].Suit)*10 + t.HandTiles[i].Number
		jValue := int(t.HandTiles[j].Suit)*10 + t.HandTiles[j].Number
		return iValue < jValue
	})
}
