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
