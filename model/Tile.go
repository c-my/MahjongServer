package model

type Suit int

const (
	Character Suit = iota
	Bamboo
	Dot
	Wind
	Dragon
	Flower
	Season
)

type Tile struct {
	Suit   Suit `json:"suit"`
	Number int  `json:"number"`
}

func NewTile(suit Suit, number int) *Tile {
	return &Tile{Suit: suit, Number: number}
}

func (t *Tile) Equals(another Tile) bool {
	return t.Suit == another.Suit && t.Number == another.Number
}
