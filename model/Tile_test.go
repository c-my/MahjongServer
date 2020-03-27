package model

import "testing"

func TestNewTile(t *testing.T) {
	tile := NewTile(Character, 3)
	if tile.Suit != Character || tile.Number != 3 {
		t.Errorf("wrong attribute value")
	}
}

func TestTile_GetLeftTile(t *testing.T) {
	var tests = []struct {
		input Tile
		want  *Tile
	}{
		{Tile{0, 1}, nil},
		{Tile{1, 1}, nil},
		{Tile{2, 1}, nil},
		{Tile{0, 2}, &Tile{0, 1}},
		{Tile{1, 3}, &Tile{1, 2}},
		{Tile{4, 1}, nil},
	}
	for _, test := range tests {
		if test.input.GetLeftTile() == nil && test.want == nil {
			continue
		} else if !test.input.GetLeftTile().Equals(*test.want) {
			t.Errorf("get wrong left tile")
		}
	}
}

func TestTile_GetRightTile(t *testing.T) {
	var tests = []struct {
		input Tile
		want  *Tile
	}{
		{Tile{0, 9}, nil},
		{Tile{1, 9}, nil},
		{Tile{2, 9}, nil},
		{Tile{0, 2}, &Tile{0, 3}},
		{Tile{1, 3}, &Tile{1, 4}},
		{Tile{4, 1}, nil},
	}
	for _, test := range tests {
		if test.input.GetRightTile() == nil && test.want == nil {
			continue
		} else if !test.input.GetRightTile().Equals(*test.want) {
			t.Errorf("get wrong left tile")
		}
	}
}

func TestIsAllSeqOrTriplet(t *testing.T) {
	var tests = []struct {
		tiles []Tile
		want  bool
	}{
		{[]Tile{{0, 7}, {0, 8}, {0, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 6}, {2, 7}, {2, 8}},
			true,
		},
	}
	for _, test := range tests {
		if IsAllSeqOrTriplet(test.tiles, 0) != test.want {
			t.Errorf("not all seq or triplet")
		}
	}
}

func TestGetTileCount(t *testing.T) {
	tiles := []Tile{{0, 7}, {0, 8}, {0, 9},
		{1, 1}, {1, 1}, {1, 1},
		{1, 2}, {1, 2},
		{2, 1}, {2, 2},
		{2, 3}, {2, 6},
		{2, 7}, {2, 8},
		{2, 9}}
	var tests = []struct {
		input Tile
		want  int
	}{
		{Tile{0, 7}, 1},
		{Tile{1, 2}, 2},
		{Tile{1, 1}, 3},
		{Tile{3, 1}, 0},
		{Tile{2, 1}, 1},
		{Tile{2, 2}, 1},
		{Tile{0, 0}, 0},
	}

	for _, test := range tests {
		if GetTileCount(tiles, test.input) != test.want {
			t.Errorf("get wrong tile count")
		}
	}
}

func TestGetTilePos(t *testing.T) {
	tiles := []Tile{{0, 7}, {0, 8}, {0, 9},
		{1, 1}, {1, 1}, {1, 1},
		{1, 2}, {1, 2},
		{2, 1}, {2, 2},
		{2, 3}, {2, 6},
		{2, 7}, {2, 8},
		{2, 9}}
	var tests = []struct {
		input Tile
		want  int
	}{
		{Tile{0, 7}, 0},
		{Tile{0, 8}, 1},
		{Tile{0, 9}, 2},
		{Tile{1, 1}, 3},
		{Tile{1, 2}, 6},
		{Tile{2, 1}, 8},
		{Tile{2, 2}, 9},
		{Tile{2, 3}, 10},
		{Tile{2, 6}, 11},
		{Tile{2, 7}, 12},
		{Tile{2, 8}, 13},
		{Tile{2, 9}, 14},
		{Tile{0, 0}, -1},
		{Tile{1, 3}, -1},
	}
	for _, test := range tests {
		if GetTilePos(tiles, test.input) != test.want {
			t.Errorf("get wrong tile postion")
		}
	}
}
