package model

import (
	"sort"
	"testing"
	"time"
)

func TestWall_Shuffle(t *testing.T) {
	wall := NewWall(generateSortedTiles)
	wall.Shuffle()
	isSorted := sort.SliceIsSorted(wall.tiles, func(i, j int) bool {
		iValue := int(wall.tiles[i].Suit)*10 + wall.tiles[i].Number
		jValue := int(wall.tiles[j].Suit)*10 + wall.tiles[j].Number
		return iValue < jValue
	})
	if isSorted {
		t.Errorf("wall is not shuffled")
	}
}

func TestWall_ReShuffle(t *testing.T) {
	wall:=NewWall(generateSortedTiles)
	oldTiles := make([]Tile, wall.Length())
	copy(oldTiles, wall.tiles)
	// sleep to get a different seed
	time.Sleep(1)
	wall.ReShuffle()
	if isTileSliceEqual(oldTiles, wall.tiles){
		t.Errorf("reshuffled walls is same as former")
	}
}

func TestWall_FrontDraw(t *testing.T) {
	wall := NewWall(generateTiles)
	firstTile := wall.tiles[0]
	oldFront := wall.front
	newTile := wall.FrontDraw()
	newFront := wall.front
	if firstTile != newTile || newFront != oldFront+1 {
		t.Errorf("something wrong when draw from front")
	}
}

func TestWall_BackDraw(t *testing.T) {
	wall := NewWall(generateTiles)
	lastTile := wall.tiles[len(wall.tiles)-1]
	oldBack := wall.back
	newTile := wall.BackDraw()
	newBack := wall.back
	if lastTile != newTile || newBack != oldBack-1 {
		t.Errorf("something wrong when draw from back")
	}
}

func TestWall_Length(t *testing.T) {
	wall := NewWall(generateSortedTiles)
	initLen := wall.Length()
	wall.FrontDraw()
	wall.BackDraw()
	if initLen != wall.Length()+2 {
		t.Errorf("wrong wall length caculated")
	}
}

func generateSortedTiles() []Tile {
	tiles := make([]Tile, 0)
	// generate all tiles
	for s := 0; s < 4; s++ {
		for i := 1; i <= 9; i++ {
			tiles = append(tiles, Tile{Character, i},
				Tile{Bamboo, i},
				Tile{Dot, i})
		}
		tiles = append(tiles, Tile{Dragon, 1})
	}
	sort.Slice(tiles, func(i, j int) bool {
		iValue := int(tiles[i].Suit)*10 + tiles[i].Number
		jValue := int(tiles[j].Suit)*10 + tiles[j].Number
		return iValue < jValue
	})
	return tiles
}

func isTileSliceEqual(s1, s2 []Tile)bool{
	if len(s1)!= len(s2){
		return false
	}
	for i:=0;i< len(s1);i++{
		if !s1[i].Equals(s2[i]){
			return false
		}
	}
	return true
}
