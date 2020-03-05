package model

import (
	"sort"
)

const (
	Character int = iota
	Bamboo
	Dot
	Wind
	Dragon
	Flower
	Season
	NAT
)

type Tile struct {
	Suit   int `json:"suit"`
	Number int `json:"number"`
}

func NewTile(suit int, number int) *Tile {
	return &Tile{Suit: suit, Number: number}
}

func (t *Tile) Equals(another Tile) bool {
	return t.Suit == another.Suit && t.Number == another.Number
}

func (t *Tile) Less(another Tile) bool {
	iValue := int(t.Suit)*10 + t.Number
	jValue := int(another.Suit)*10 + another.Number
	return iValue < jValue
}

func (t *Tile) IsSuit() bool {
	return t.Suit == Character || t.Suit == Bamboo || t.Suit == Dot
}

func SortTiles(tiles []Tile) {
	sort.Slice(tiles, func(i, j int) bool {
		iValue := int(tiles[i].Suit)*10 + tiles[i].Number
		jValue := int(tiles[j].Suit)*10 + tiles[j].Number
		return iValue < jValue
	})
}

//GetPairPos  return positions of all pair
//the tiles needs to be sorted
func GetPairPos(tiles []Tile) []int {
	pos := make([]int, 0)
	length := len(tiles) - 1
	for i := 0; i < length; i++ {
		if tiles[i].Equals(tiles[i+1]) {
			pos = append(pos, i)
			i++
		}
	}
	return pos
}

func RemovePair(tiles []Tile, pos int) []Tile {
	remainTiles := make([]Tile, 0, len(tiles)-2)
	remainTiles = append(remainTiles, tiles[:pos]...)
	remainTiles = append(remainTiles, tiles[pos+2:]...)
	return remainTiles
}

func FindAndRemoveTriplet(tiles *[]Tile) bool {
	var v = *tiles
	if IsTriplet(v[0], v[1], v[2]) {
		newTiles := make([]Tile, 0)
		*tiles = append(newTiles, v[3:]...)
		return true
	}
	return false
}

func FindAndRemoveSequence(tiles *[]Tile) bool {
	tmp := make([]Tile, 0)
	var v = *tiles
	for i := 1; i < len(v); i++ {
		switch {
		case v[i].Equals(v[i-1]):
			tmp = append(tmp, v[i])
		case v[i-1].GetRightTile() != nil && v[i].Equals(*(v[i-1].GetRightTile())):
			if v[i].Suit == v[0].Suit && v[i].Number-v[0].Number == 2 {
				tmp = append(tmp, v[i+1:]...)
				*tiles = tmp
				return true
			}
		default:
			return false
		}
	}
	return false
}

func IsAllSeqOrTriplet(tiles []Tile, tripCount int) bool {
	length := len(tiles)
	for i := 0; i < length/3; i++ {
		find := FindAndRemoveTriplet(&tiles)
		if find {
			tripCount--
		}
		if !find {
			find = FindAndRemoveSequence(&tiles)
		}
		if !find {
			return false
		}
	}
	return len(tiles) == 0 && tripCount <= 0
}

func HasCharacter(tiles []Tile) bool {
	for _, t := range tiles {
		if t.Suit == Character {
			return true
		}
		if t.Suit > Character {
			return false
		}
	}
	return false
}

func HasBamboo(tiles []Tile) bool {
	for _, t := range tiles {
		if t.Suit == Bamboo {
			return true
		}
		if t.Suit > Bamboo {
			return false
		}
	}
	return false
}

func HasDot(tiles []Tile) bool {
	for _, t := range tiles {
		if t.Suit == Dot {
			return true
		}
		if t.Suit > Dot {
			return false
		}
	}
	return false
}

func IsSequence(t1, t2, t3 Tile) bool {
	if !t1.IsSuit() || !t2.IsSuit() || !t3.IsSuit() {
		return false
	}
	if t2.GetLeftTile() == nil || t2.GetRightTile() == nil {
		return false
	}
	return t1.Equals(*t2.GetLeftTile()) && t3.Equals(*t2.GetRightTile())
}

func IsTriplet(t1, t2, t3 Tile) bool {
	return t1.Equals(t2) && t2.Equals(t3)
}

//GetTileCount return count of target in sorted tiles
func GetTileCount(tiles []Tile, target Tile) int {
	count := 0
	for i := 0; i < len(tiles); i++ {
		if tiles[i].Equals(target) {
			count++
		} else if target.Less(tiles[i]) {
			break
		}
	}
	return count
}

func (t *Tile) GetLeftTile() *Tile {
	if t.Suit == Character || t.Suit == Bamboo || t.Suit == Dot {
		if t.Number == 1 {
			return nil
		}
		return &Tile{
			Suit:   t.Suit,
			Number: t.Number - 1,
		}
	}
	return nil
}

func (t *Tile) GetRightTile() *Tile {
	if t.Suit == Character || t.Suit == Bamboo || t.Suit == Dot {
		if t.Number == 9 {
			return nil
		}
		return &Tile{
			Suit:   t.Suit,
			Number: t.Number + 1,
		}
	}
	return nil
}
