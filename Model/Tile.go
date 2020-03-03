package Model

type Suit int

const (
	Character Suit = iota
	Bamboo
	Dot
	Bonus
	Wind
	Dragon
	Flower
	Season
)

type Tile struct {
	Suit   Suit `json:"suit"`
	Number int  `json:"number"`
}
