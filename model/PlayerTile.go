package model

type PlayerTile struct {
	HandTiles  []Tile
	ShownTiles []ShownTile
}

func (t *PlayerTile) SetHandTiles(tiles []Tile) {
	t.HandTiles = tiles
	t.SortHand()
}

func (t *PlayerTile) GetHandTiles() []Tile {
	return t.HandTiles
}

func (t *PlayerTile) SortHand() {
	SortTiles(t.HandTiles)
}
