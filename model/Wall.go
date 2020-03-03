package model

import (
	"math/rand"
	"time"
)

type Wall struct {
	tiles     []Tile
	front     int
	back      int
	generator func() []Tile
}

func NewWall(generator func() []Tile) *Wall {
	t := generator()
	wall := &Wall{tiles: t, front: 0, back: len(t) - 1, generator: generator}
	wall.Shuffle()
	return wall
}

func (w *Wall) FrontDraw() Tile {
	if w.back < w.front {
		panic("front-draw from empty wall")
	}
	tile := w.tiles[w.front]
	w.front++
	return tile
}

func (w *Wall) backDraw() Tile {
	if w.back < w.front {
		panic("back-draw from empty wall")
	}
	tile := w.tiles[w.back]
	w.back--
	return tile
}

func (w *Wall) ReShuffle() {
	w.tiles = w.generator()
	w.front = 0
	w.back = len(w.tiles) - 1
	w.Shuffle()
}

func (w *Wall) Length() int {
	return w.back - w.front + 1
}

func (w *Wall) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(w.tiles), func(i, j int) {
		w.tiles[i], w.tiles[j] = w.tiles[j], w.tiles[i]
	})
}
