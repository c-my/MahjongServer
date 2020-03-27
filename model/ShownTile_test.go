package model

import (
	"github.com/c-my/MahjongServer/config"
	"testing"
)

func TestGetChowCount(t *testing.T) {
	var tests = []struct {
		input []ShownTile
		want  int
	}{
		{[]ShownTile{{config.ConcealedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}}, 0},
		{[]ShownTile{}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}}}, 1},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}},
			{config.Chow, []Tile{{1, 5}, {1, 6}, {1, 7}}}}, 2},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}},
			{config.Pong, []Tile{{3, 1}, {3, 1}, {3, 1}}}}, 0},
	}
	for _, test := range tests {
		if GetChowCount(test.input) != test.want {
			t.Errorf("get wrong chow count")
		}
	}
}

func TestGetPongCount(t *testing.T) {
	var tests = []struct {
		input []ShownTile
		want  int
	}{
		{[]ShownTile{{config.ConcealedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}}, 0},
		{[]ShownTile{}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}}}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}},
			{config.Chow, []Tile{{1, 5}, {1, 6}, {1, 7}}}}, 0},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}}}, 1},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}},
			{config.Pong, []Tile{{3, 1}, {3, 1}, {3, 1}}}}, 2},
	}
	for _, test := range tests {
		if GetPongCount(test.input) != test.want {
			t.Errorf("get wrong pong count")
		}
	}
}

func TestGetExposedKongCount(t *testing.T) {
	var tests = []struct {
		input []ShownTile
		want  int
	}{
		{[]ShownTile{{config.ConcealedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}}, 0},
		{[]ShownTile{}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}}}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}},
			{config.Chow, []Tile{{1, 5}, {1, 6}, {1, 7}}}}, 0},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}}}, 0},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}},
			{config.Pong, []Tile{{3, 1}, {3, 1}, {3, 1}}}}, 0},
		{[]ShownTile{{config.ExposedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}}, 1},
		{[]ShownTile{{config.ExposedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}},
			{config.ExposedKong, []Tile{{3, 1}, {3, 1}, {3, 1}, {3, 1}}}}, 2},
	}
	for _, test := range tests {
		if GetExposedKongCount(test.input) != test.want {
			t.Errorf("get wrong exposed kong count")
		}
	}
}

func TestGetConcealedKongCount(t *testing.T) {
	var tests = []struct {
		input []ShownTile
		want  int
	}{
		{[]ShownTile{{config.ConcealedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}}, 1},
		{[]ShownTile{}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}}}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}},
			{config.Chow, []Tile{{1, 5}, {1, 6}, {1, 7}}}}, 0},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}}}, 0},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}},
			{config.Pong, []Tile{{3, 1}, {3, 1}, {3, 1}}}}, 0},
		{[]ShownTile{{config.ExposedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}}, 0},
		{[]ShownTile{{config.ExposedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}},
			{config.ConcealedKong, []Tile{{3, 1}, {3, 1}, {3, 1}, {3, 1}}}}, 1},
	}
	for _, test := range tests {
		if GetConcealedKongCount(test.input) != test.want {
			t.Errorf("get wrong concealed kong count")
		}
	}
}

func TestGetAddedKongCount(t *testing.T) {
	var tests = []struct {
		input []ShownTile
		want  int
	}{
		{[]ShownTile{{config.AddedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}}, 1},
		{[]ShownTile{}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}}}, 0},
		{[]ShownTile{{config.Chow, []Tile{{0, 3}, {0, 4}, {0, 5}}},
			{config.Chow, []Tile{{1, 5}, {1, 6}, {1, 7}}}}, 0},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}}}, 0},
		{[]ShownTile{{config.Pong, []Tile{{0, 3}, {0, 3}, {0, 3}}},
			{config.Pong, []Tile{{3, 1}, {3, 1}, {3, 1}}}}, 0},
		{[]ShownTile{{config.ExposedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}}, 0},
		{[]ShownTile{{config.AddedKong, []Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}},
			{config.ConcealedKong, []Tile{{3, 1}, {3, 1}, {3, 1}, {3, 1}}}}, 1},
	}
	for _, test := range tests {
		if GetAddedKongCount(test.input) != test.want {
			t.Errorf("get wrong added-kong count")
		}
	}
}
