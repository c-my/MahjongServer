package model

import (
	"github.com/c-my/MahjongServer/config"
)

type ShownTile struct {
	ShownType int    `json:"shown_type"`
	Tiles     []Tile `json:"tiles"`
}

func GetChowCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == config.Chow {
			count++
		}
	}
	return count
}
func GetPongCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == config.Pong {
			count++
		}
	}
	return count
}

func GetExposedKongCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == config.ExposedKong {
			count++
		}
	}
	return count
}

func GetConcealedKongCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == config.ConcealedKong {
			count++
		}
	}
	return count
}

func GetAddedKongCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == config.AddedKong {
			count++
		}
	}
	return count
}
