package model

import (
	"sort"
	"testing"
)

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
	var p *Tile
	p = nil
	if p.GetLeftTile() != nil {
		t.Errorf("left of nil isn't nil")
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
	var p *Tile
	p = nil
	if p.GetRightTile() != nil {
		t.Errorf("right of nil isn't nil")
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
		{[]Tile{{0, 7}, {0, 7}, {0, 7},
			{2, 1}, {2, 2}, {2, 3}, {2, 6}, {2, 7}, {2, 8}},
			true,
		},
		{[]Tile{{0, 7}, {0, 7}, {0, 8},
			{2, 1}, {2, 2}, {2, 3}, {2, 6}, {2, 7}, {2, 8}},
			false,
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

func TestIsSequence(t *testing.T) {
	var tests = []struct {
		input []Tile
		want  bool
	}{
		{[]Tile{{3, 1}, {3, 1}, {3, 1}}, false},
		{[]Tile{{1, 1}, {3, 1}, {2, 1}}, false},
		{[]Tile{{1, 1}, {1, 1}, {1, 1}}, false},
		{[]Tile{{1, 1}, {0, 1}, {1, 1}}, false},
		{[]Tile{{2, 5}, {2, 6}, {2, 7}}, true},
		{[]Tile{{2, 7}, {2, 8}, {2, 9}}, true},
		{[]Tile{{1, 1}, {1, 2}, {1, 3}}, true},
		{[]Tile{{0, 1}, {0, 2}, {0, 3}}, true},
	}
	for _, test := range tests {
		if IsSequence(test.input[0], test.input[1], test.input[2]) != test.want {
			t.Errorf("wrong judge sequence")
		}
	}
}

func TestHasCharacter(t *testing.T) {
	var tests = []struct {
		input []Tile
		want  bool
	}{
		{[]Tile{{0, 7}, {0, 8}, {0, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 6}, {2, 7}, {2, 8}},
			true,
		},
		{[]Tile{{1, 7}, {1, 8}, {1, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 6},
			{2, 7}, {2, 8}, {3, 1}},
			false,
		},
	}
	for _, test := range tests {
		if HasCharacter(test.input) != test.want {
			t.Errorf("wrong jugde existence of character tile")
		}
	}
}

func TestHasBamboo(t *testing.T) {
	var tests = []struct {
		input []Tile
		want  bool
	}{
		{[]Tile{{0, 7}, {0, 8}, {0, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 6}, {2, 7}, {2, 8}},
			false,
		},
		{[]Tile{{1, 7}, {1, 8}, {1, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 6},
			{2, 7}, {2, 8}, {3, 1}},
			true,
		},
	}
	for _, test := range tests {
		if HasBamboo(test.input) != test.want {
			t.Errorf("wrong jugde existence of bamboo tile")
		}
	}
}

func TestHasDot(t *testing.T) {
	var tests = []struct {
		input []Tile
		want  bool
	}{
		{[]Tile{{0, 7}, {0, 8}, {0, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 6}, {2, 7}, {2, 8}},
			true,
		},
		{[]Tile{{1, 7}, {1, 8}, {1, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 6},
			{2, 7}, {2, 8}, {3, 1}},
			true,
		},
	}
	for _, test := range tests {
		if HasDot(test.input) != test.want {
			t.Errorf("wrong jugde existence of dot tile")
		}
	}
}

func TestGetPairPos(t *testing.T) {
	var tests = []struct {
		input []Tile
		want  []int
	}{
		{[]Tile{}, []int{}},
		{[]Tile{{0, 7}, {3, 1}, {3, 1}, {3, 1}}, []int{1}},
	}
	for _, test := range tests {
		res := GetPairPos(test.input)
		if len(res) != len(test.want) {
			t.Errorf("wrong pair pos")
		}
		sort.Ints(res)
		for i := 0; i < len(res); i++ {
			if res[i] != test.want[i] {
				t.Errorf("wrong pair pos")
			}
		}
	}
}

func TestRemoveTile(t *testing.T) {
	var tests = []struct {
		input []Tile
		tile  Tile
		count int
		want  []Tile
	}{
		{[]Tile{{0, 1}, {0, 2}, {0, 3}},
			Tile{0, 1},
			1,
			[]Tile{{0, 2}, {0, 3}}},
		{[]Tile{{0, 1}, {0, 2}, {0, 3}, {1, 2}, {1, 2}},
			Tile{1, 2},
			2,
			[]Tile{{0, 1}, {0, 2}, {0, 3}}},
		{[]Tile{{1, 2}, {1, 2}},
			Tile{1, 2},
			2,
			[]Tile{}},
	}
	for _, test := range tests {
		res := RemoveTile(test.input, test.tile, test.count)
		l := len(res)
		if l != len(test.want) {
			t.Errorf("failed to remove tile")
		}
		for i := 0; i < l; i++ {
			if !res[i].Equals(test.want[i]) {
				t.Errorf("failed to remove tile")
			}
		}
	}
}

func TestRemovePair(t *testing.T) {
	var tests = []struct {
		tiles []Tile
		pos   int
		want  []Tile
	}{
		{[]Tile{{0, 1}, {0, 2}, {0, 3}, {1, 2}, {1, 2}},
			3,
			[]Tile{{0, 1}, {0, 2}, {0, 3}}},
	}
	for _, test := range tests {
		res := RemovePair(test.tiles, test.pos)
		l := len(res)
		if l != len(test.want) {
			t.Errorf("failed to remove pair")
		}
		for i := 0; i < l; i++ {
			if !res[i].Equals(test.want[i]) {
				t.Errorf("failed to remove pair")
			}
		}
	}
}
