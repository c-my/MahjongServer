package model

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestPlayerTile_SortHand(t *testing.T) {
	playerTile := PlayerTile{}
	playerTile.SetHandTiles(generateTiles())
	playerTile.SortHand()
	hand := playerTile.GetHandTiles()
	isSorted := sort.SliceIsSorted(hand, func(i, j int) bool {
		iValue := int(hand[i].Suit)*10 + hand[i].Number
		jValue := int(hand[j].Suit)*10 + hand[j].Number
		return iValue < jValue
	})
	if !isSorted {
		t.Errorf("hand tiles are not sorted")
	}
}

func generateTiles() []Tile {
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
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tiles), func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})
	return tiles
}
