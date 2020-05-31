package rule

import (
	"github.com/c-my/MahjongServer/config"
	"github.com/c-my/MahjongServer/model"
	"testing"
)

func TestJinzhouRule_CanChow(t *testing.T) {
	j := NewJinzhouRule()
	var tests = []struct {
		hand    []model.Tile
		newTile model.Tile
		want    bool
	}{
		{[]model.Tile{{0, 1}, {0, 3}, {1, 4}}, model.Tile{0, 2}, true},
		{[]model.Tile{{0, 2}, {0, 3}, {1, 4}}, model.Tile{0, 2}, false},
		{[]model.Tile{{0, 2}, {0, 3}, {1, 4}}, model.Tile{0, 1}, true},
		{[]model.Tile{{0, 1}, {1, 3}, {2, 4}}, model.Tile{0, 2}, false},
		{[]model.Tile{{1, 7}, {1, 8}, {2, 4}}, model.Tile{1, 9}, true},
		{[]model.Tile{{1, 7}, {1, 8}, {2, 4}}, model.Tile{4, 1}, false},
	}
	for _, test := range tests {
		canChow, _ := j.CanChow(test.hand, test.newTile)
		if canChow != test.want {
			t.Errorf("wrong chow")
		}
	}
}

func TestJinzhouRule_CanPong(t *testing.T) {
	j := NewJinzhouRule()
	var tests = []struct {
		hand    []model.Tile
		newTile model.Tile
		want    bool
	}{
		{[]model.Tile{{4, 1}, {4, 1}, {4, 1}}, model.Tile{4, 1}, true},
		{[]model.Tile{{1, 4}, {4, 1}, {4, 1}}, model.Tile{4, 1}, true},
		{[]model.Tile{{0, 1}, {1, 1}, {3, 1}, {4, 1}}, model.Tile{4, 1}, false},
		{[]model.Tile{{0, 1}, {0, 3}, {1, 4}}, model.Tile{0, 2}, false},
	}
	for _, test := range tests {
		if j.CanPong(test.hand, test.newTile) != test.want {
			t.Errorf("wrong pong")
		}
	}
}

func TestJinzhouRule_CanConcealedKong(t *testing.T) {
	j := NewJinzhouRule()
	var tests = []struct {
		hand    []model.Tile
		newTile model.Tile
		want    bool
	}{
		{[]model.Tile{{4, 1}, {4, 1}, {4, 1}}, model.Tile{4, 1}, true},
		{[]model.Tile{{4, 1}, {4, 1}, {4, 1}}, model.Tile{1, 1}, false},
		{[]model.Tile{{0, 1}, {1, 1}, {3, 1}, {4, 1}}, model.Tile{4, 1}, false},
		{[]model.Tile{{0, 1}, {0, 3}, {1, 4}}, model.Tile{0, 2}, false},
	}
	for _, test := range tests {
		if j.CanConcealedKong(test.hand, test.newTile) != test.want {
			t.Errorf("wrong pong")
		}
	}
}

func TestJinzhouRule_CanWin(t *testing.T) {
	j := NewJinzhouRule()
	var tests = []struct {
		hand    []model.Tile
		shown   []model.ShownTile
		newTile model.Tile
		want    bool
	}{
		{[]model.Tile{{0, 7}, {0, 8}, {0, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 7}, {2, 8}, {2, 9}, {2, 9}},
			[]model.ShownTile{{config.ConcealedKong, []model.Tile{{1, 3}, {1, 3}, {1, 3}, {1, 3}}}},
			model.Tile{2, 6},
			true},
		{[]model.Tile{{0, 7}, {0, 7}, {0, 9}, {0, 9},
			{1, 3},
			{2, 1}, {2, 1}, {2, 3}, {2, 3}, {2, 8}, {2, 8}, {2, 9}, {2, 9}},
			[]model.ShownTile{},
			model.Tile{1, 3},
			true}, //7对
		{[]model.Tile{{0, 2}, {0, 3}, {0, 4}, {0, 6}, {0, 7}, {0, 8},
			{2, 5}, {2, 7}, {2, 8}, {2, 8}},
			[]model.ShownTile{{config.ConcealedKong, []model.Tile{{1, 3}, {1, 3}, {1, 3}, {1, 3}}}},
			model.Tile{2, 6},
			false}, //缺幺九
		{[]model.Tile{{0, 2}, {0, 3}, {0, 4}, {0, 6}, {0, 7}, {0, 8},
			{2, 5}, {2, 7}, {2, 8}, {2, 8}},
			[]model.ShownTile{{config.ConcealedKong, []model.Tile{{1, 9}, {1, 9}, {1, 9}, {1, 9}}}},
			model.Tile{2, 6},
			true},
		{[]model.Tile{{0, 7}, {0, 8}, {0, 9},
			{1, 1}, {1, 2},
			{2, 1}, {2, 2}, {2, 3}, {2, 6}, {2, 7}, {2, 8}, {2, 9}},
			[]model.ShownTile{},
			model.Tile{1, 3},
			false},
		{[]model.Tile{{0, 9}, {0, 9}, {0, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 7}, {2, 8}, {2, 9}, {2, 9}},
			[]model.ShownTile{{config.ConcealedKong, []model.Tile{{1, 3}, {1, 3}, {1, 3}, {1, 3}}}},
			model.Tile{2, 6},
			true},
		{[]model.Tile{{0, 7}, {0, 8}, {0, 9},
			{2, 1}, {2, 2}, {2, 3}, {2, 7}, {2, 8}, {2, 9}, {2, 9}},
			[]model.ShownTile{{config.ConcealedKong, []model.Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}},
			model.Tile{2, 6},
			false}, //缺门
		{[]model.Tile{{0, 1}, {0, 1},
			{0, 4},	{0, 4}, {0, 4},
			{0, 6}, {0, 8},
			{1, 5}, {1, 6}, {1, 7},
			{2,9},{2,9},{2,9}},
			[]model.ShownTile{},
			model.Tile{0, 7},
			true},
	}
	for _, test := range tests {
		if j.CanWin(test.hand, test.shown, test.newTile) != test.want {
			t.Errorf("wrong win")
		}
	}
}

func BenchmarkJinzhouRule_CanWin_Normal(b *testing.B) {
	j := NewJinzhouRule()
	var test = struct {
		hand    []model.Tile
		shown   []model.ShownTile
		newTile model.Tile
		want    bool
	}{[]model.Tile{{0, 7}, {0, 8}, {0, 9},
		{2, 1}, {2, 2}, {2, 3}, {2, 7}, {2, 8}, {2, 9}, {2, 9}},
		[]model.ShownTile{{config.ConcealedKong, []model.Tile{{1, 3}, {1, 3}, {1, 3}, {1, 3}}}},
		model.Tile{2, 6}, true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.CanWin(test.hand, test.shown, test.newTile)
	}
}

func BenchmarkJinzhouRule_CanWin_SevenPair(b *testing.B) {
	j := NewJinzhouRule()
	var test = struct {
		hand    []model.Tile
		shown   []model.ShownTile
		newTile model.Tile
		want    bool
	}{[]model.Tile{{0, 7}, {0, 7}, {0, 9}, {0, 9},
		{1, 3},
		{2, 1}, {2, 1}, {2, 3}, {2, 3}, {2, 8}, {2, 8}, {2, 9}, {2, 9}},
		[]model.ShownTile{},
		model.Tile{1, 3},
		true} //7对
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.CanWin(test.hand, test.shown, test.newTile)
	}
}

func BenchmarkJinzhouRule_CanWin_LackDoor(b *testing.B) {
	j := NewJinzhouRule()
	var test = struct {
		hand    []model.Tile
		shown   []model.ShownTile
		newTile model.Tile
		want    bool
	}{[]model.Tile{{0, 7}, {0, 8}, {0, 9},
		{2, 1}, {2, 2}, {2, 3}, {2, 7}, {2, 8}, {2, 9}, {2, 9}},
		[]model.ShownTile{{config.ConcealedKong, []model.Tile{{0, 3}, {0, 3}, {0, 3}, {0, 3}}}},
		model.Tile{2, 6},
		false,
	} //缺门
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.CanWin(test.hand, test.shown, test.newTile)
	}
}
