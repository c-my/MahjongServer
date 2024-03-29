package model

type PlayerTile struct {
	HandTiles  []Tile      `json:"hand_tiles"`
	ShownTiles []ShownTile `json:"shown_tiles"`
	DropTiles  []Tile      `json:"drop_tiles"`
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
